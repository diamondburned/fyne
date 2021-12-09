package wayland

import (
	"image"
	"log"
	"math"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/driver/wayland/internal/painter"
	"fyne.io/fyne/v2/internal/driver/wayland/internal/swizzle"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	"golang.org/x/sys/unix"
)

type canvas struct {
	common.Canvas
	img     image.RGBA // need swizzling
	content fyne.CanvasObject

	win *waylandWindow
	fd  *os.File
	buf buffer

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)

	queued bool // draw
}

var _ fyne.Canvas = (*canvas)(nil)

func newCanvas(w *waylandWindow) canvas {
	fd, err := unix.MemfdCreate("fyne_wl_shm", unix.MFD_HUGE_8MB)
	must(err, "cannot create memfd for wl_shm")

	return canvas{
		win: w,
		fd:  os.NewFile(uintptr(fd), "fyne_wl_shm"),
	}
}

func (c *canvas) destroy() {
	c.Lock()
	c.buf.destroy()
	c.fd.Close()
	c.Unlock()
}

type buffer struct {
	*client.Buffer
	pix []uint8
}

func remapFd(s *session, f *os.File, img *image.RGBA) buffer {
	if len(img.Pix) > math.MaxInt32 {
		log.Panicf("canvas size %d overflows int32", len(img.Pix))
	}

	// TODO: round this up so we can overalloc.
	err := f.Truncate(int64(len(img.Pix)))
	must(err, "cannot resize Wayland image buffer file")

	data, err := unix.Mmap(
		int(f.Fd()), 0, len(img.Pix),
		unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	must(err, "cannot mmap Wayland image buffer file")

	pool, err := s.shm.CreatePool(f.Fd(), int32(len(img.Pix)))
	must(err, "cannot create Wayland shm pool")

	defer func() {
		err = pool.Destroy()
		must(err, "cannot destroy Wayland shm pool")
	}()

	buff, err := pool.CreateBuffer(0,
		int32(img.Rect.Max.X), int32(img.Rect.Max.Y), int32(img.Stride),
		uint32(client.ShmFormatArgb8888))
	must(err, "cannot create Wayland shm buffer")

	b := buffer{
		Buffer: buff,
		pix:    data,
	}

	buff.AddReleaseHandler(func(client.BufferReleaseEvent) {
		b.destroy()
	})

	return b
}

func (b *buffer) destroy() {
	var err error

	err = b.Buffer.Destroy()
	must(err, "cannot destroy Wayland buffer")

	err = unix.Munmap(b.pix)
	must(err, "cannot unmap Wayland canvas buffer")
}

func (b *buffer) flush(img *image.RGBA) {
	copy(b.pix, img.Pix)
	swizzle.BGRA(b.pix)
}

func (c *canvas) reset(x, y int) {
	// check that the new sizes don't overflow int32
	if x > math.MaxInt32 || y > math.MaxInt32 {
		log.Panicf("(x, y) = (%d, %d) overflows int32", x, y)
	}
	// check if point is either zero or is what we already have
	if pt := image.Pt(x, y); pt == image.ZP || pt == c.img.Rect.Max {
		return
	}

	c.img = *image.NewRGBA(image.Rect(0, 0, x, y))
	c.win.state.use(func(s *session) {
		c.buf = remapFd(s, c.fd, &c.img)
	})

	err := c.win.surface.Attach(c.buf.Buffer, 0, 0)
	must(err, "cannot attach new Wayland canvas buffer")
}

func (c *canvas) damage(pos fyne.Position, size fyne.Size) {
	x := int32(internal.ScaleInt(c, pos.X))
	y := int32(internal.ScaleInt(c, pos.Y))
	w := int32(internal.ScaleInt(c, size.Width))
	h := int32(internal.ScaleInt(c, size.Height))

	err := c.win.surface.DamageBuffer(x, y, w, h)
	must(err, "cannot damage Wayland region")
}

func (c *canvas) damageAll() {
	err := c.win.surface.DamageBuffer(0, 0, math.MaxInt32, math.MaxInt32)
	must(err, "cannot damage Wayland surface")
}

// redraw redraws the application into the internal buffer. It does not mark the
// region as damaged, so the user has to do that.
//
func (c *canvas) redraw() {
	// TODO: do damage tracking and only walk the widgets that we actually need
	// to redraw.
	painter.Paint(&c.img, c)
	c.buf.flush(&c.img)
}

func (c *canvas) commit() {
	err := c.win.surface.Commit()
	must(err, "cannot commit Wayland surface")
}

// queueDraw queues the draw to the server. It flushes the application state
// into the local buffer.
func (c *canvas) queueDraw() {
	c.redraw()

	if c.queued {
		return
	}
	c.queued = true

	f, err := c.win.surface.Frame()
	must(err, "cannot create surface frame callback")

	f.AddDoneHandler(func(client.CallbackDoneEvent) {
		c.Lock()
		c.queued = false
		c.commit()
		c.Unlock()
	})
}

func (c *canvas) Refresh(obj fyne.CanvasObject) {
	c.damage(obj.Position(), obj.Size())
	c.queueDraw()
}

func (c *canvas) Scale() float32 {
	// TODO: see client.NewOutput.
	return 1.0
}

func (c *canvas) PixelCoordinateForPosition(pos fyne.Position) (x, y int) {
	return internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
}

func (c *canvas) Resize(size fyne.Size) {
	panic("TODO")
}

func (c *canvas) Size() fyne.Size {
	return fyne.Size{
		Width:  internal.UnscaleInt(c, c.img.Rect.Max.X),
		Height: internal.UnscaleInt(c, c.img.Rect.Max.Y),
	}
}

func (c *canvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *canvas) Content() fyne.CanvasObject {
	c.RLock()
	defer c.RUnlock()

	return c.content
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.Lock()
	c.content = content
	c.Unlock()

	c.Redraw()
}

// Redraw draws everything into an internal buffer and signals to the compositor
// to draw on the next frame.
func (c *canvas) Redraw() {
	c.RLock()
	defer c.RUnlock()

	c.damageAll()
	c.queueDraw()
}

// Capture copies the internal RGBA buffer.
func (c *canvas) Capture() image.Image {
	img := image.NewRGBA(c.img.Rect)
	copy(c.img.Pix, img.Pix)
	return img
}

func (c *canvas) SetOnKeyDown(typed func(*fyne.KeyEvent)) {
	c.Lock()
	defer c.Unlock()
	c.onKeyDown = typed
}

func (c *canvas) SetOnKeyUp(typed func(*fyne.KeyEvent)) {
	c.Lock()
	defer c.Unlock()
	c.onKeyUp = typed
}

func (c *canvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.Lock()
	defer c.Unlock()
	c.onTypedKey = typed
}

func (c *canvas) SetOnTypedRune(typed func(rune)) {
	c.Lock()
	defer c.Unlock()
	c.onTypedRune = typed
}

func (c *canvas) OnKeyDown() func(*fyne.KeyEvent) {
	c.RLock()
	defer c.RUnlock()
	return c.onKeyDown
}

func (c *canvas) OnKeyUp() func(*fyne.KeyEvent) {
	c.RLock()
	defer c.RUnlock()
	return c.onKeyUp
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.RLock()
	defer c.RUnlock()
	return c.onTypedKey
}

func (c *canvas) OnTypedRune() func(rune) {
	c.RLock()
	defer c.RUnlock()
	return c.onTypedRune
}

func (c *canvas) handleKeyboardKey(ev client.KeyboardKeyEvent) {
	var keysym uint32
	c.win.state.use(func(s *session) {
		keysym = s.keymap.OneSym(ev.Key)
	})

	k := fyne.KeyEvent{}
}
