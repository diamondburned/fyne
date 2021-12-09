package wayland

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

type waylandWindow struct {
	state *state

	surface     *client.Surface
	xdgSurface  *xdg_shell.Surface
	xdgTopLevel *xdg_shell.Toplevel
}

func (w *waylandWindow) init() {
	w.state.use(func(s *session) {
		var err error

		w.surface, err = s.compositor.CreateSurface()
		must(err, "cannot create Wayland surface")
		w.xdgSurface, err = s.xdgWmBase.GetXdgSurface(w.surface)
		must(err, "cannot get Wayland XDG surface")
		w.xdgTopLevel, err = w.xdgSurface.GetToplevel()
		must(err, "cannot get Wayland XDG top-level")
	})
}

func (w *waylandWindow) destroy() {
	var err error

	err = w.xdgTopLevel.Destroy()
	must(err, "cannot destroy Wayland XDG top-level")
	err = w.xdgSurface.Destroy()
	must(err, "cannot destroy Wayland XDG surface")
	err = w.surface.Destroy()
	must(err, "cannot destroy Wayland surface")
}

type window struct {
	mu     *sync.RWMutex
	window *waylandWindow
	canvas canvas

	onClose        func()
	closeIntercept func()

	dead  chan struct{}
	appID string
	title string

	fixedSize bool
}

var _ fyne.Window = (*window)(nil)

func newWindow(s *state, appID, title string) *window {
	wlwin := waylandWindow{
		state: s,
	}

	w := window{
		window: &wlwin,
		canvas: newCanvas(&wlwin),

		appID: appID,
		title: title,
		dead:  make(chan struct{}),
	}

	// hack
	w.mu = &w.canvas.Canvas.RWMutex

	w.window.init()
	w.window.xdgSurface.AddConfigureHandler(w.handleSurfaceConfigure)
	w.window.xdgTopLevel.AddConfigureHandler(w.handleToplevelConfigure)
	w.window.xdgTopLevel.AddCloseHandler(w.handleToplevelClose)

	err := w.window.xdgTopLevel.SetAppId(appID)
	must(err, "cannot set Wayland top-level application ID")

	w.SetTitle(title)
	w.canvas.commit()

	return &w
}

func (w *window) handleSurfaceConfigure(ev xdg_shell.SurfaceConfigureEvent) {
	var err error

	err = w.window.xdgSurface.AckConfigure(ev.Serial)
	must(err, "cannot send Wayland surface ack")

	// SurfaceConfigure should always follow after a toplevel configure event,
	// so we can assume that the buffer has already been reallocated.
	w.canvas.commit()
}

func (w *window) handleToplevelConfigure(ev xdg_shell.ToplevelConfigureEvent) {
	w.canvas.reset(int(ev.Width), int(ev.Height))
}

func (w *window) handleToplevelClose(ev xdg_shell.ToplevelCloseEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closeIntercept != nil {
		w.closeIntercept()
		return
	}

	w.close()
}

func (w *window) Title() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.title
}

func (w *window) SetTitle(title string) {
	w.mu.Lock()
	w.title = title
	w.mu.Unlock()

	err := w.window.xdgTopLevel.SetTitle(title)
	must(err, "cannot set Wayland top-level title")
}

func (w *window) FullScreen() bool {
	panic("implement me")
}

func (w *window) SetFullScreen(bool) {
	panic("implement me")
}

func (w *window) Resize(size fyne.Size) {
	panic("implement me")
	w.resize()
}

func (w *window) resize() {
	if w.fixedSize {
		pt := w.canvas.img.Rect.Max
		w.window.xdgTopLevel.SetMinSize(int32(pt.X), int32(pt.Y))
		w.window.xdgTopLevel.SetMaxSize(int32(pt.X), int32(pt.Y))
	}
}

func (w *window) RequestFocus() {
	panic("implement me")
}

func (w *window) FixedSize() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.fixedSize
}

func (w *window) SetFixedSize(fixedSize bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.fixedSize = fixedSize

	if fixedSize {
		w.resize()
	} else {
		w.window.xdgTopLevel.SetMinSize(0, 0)
		w.window.xdgTopLevel.SetMaxSize(0, 0)
	}
}

func (w *window) CenterOnScreen() {
	panic("implement me")
}

func (w *window) Padded() bool {
	panic("implement me")
}

func (w *window) SetPadded(bool) {
	panic("implement me")
}

func (w *window) Icon() fyne.Resource {
	panic("implement me")
}

func (w *window) SetIcon(fyne.Resource) {
	panic("implement me")
}

func (w *window) SetMaster() {
	panic("implement me")
}

func (w *window) MainMenu() *fyne.MainMenu {
	panic("implement me")
}

func (w *window) SetMainMenu(*fyne.MainMenu) {
	panic("implement me")
}

func (w *window) SetOnClosed(f func()) {
	w.mu.Lock()
	w.onClose = f
	w.mu.Unlock()
}

func (w *window) SetCloseIntercept(f func()) {
	w.mu.Lock()
	w.closeIntercept = f
	w.mu.Unlock()
}

func (w *window) Show() {
	// can't figure out what the unminimize thing is.
	w.canvas.Redraw()
}

func (w *window) Hide() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.window.xdgTopLevel.SetMinimized()
}

func (w *window) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.close()
}

func (w *window) close() {
	w.canvas.destroy()
	w.window.destroy()

	if w.onClose != nil {
		w.onClose()
	}

	close(w.dead)
}

func (w *window) ShowAndRun() {
	w.Show()
	w.window.state.run()
}

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.Content()
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.canvas.SetContent(content)
}

func (w *window) Canvas() fyne.Canvas {
	return &w.canvas
}

func (w *window) Clipboard() fyne.Clipboard {
	panic("implement me")
}
