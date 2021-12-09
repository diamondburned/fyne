package wayland

import (
	"errors"
	"io"
	"log"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/wayland/internal/xkbcommon"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	"github.com/rajveermalviya/go-wayland/wayland/cursor"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

type session struct {
	display     *client.Display
	registry    *client.Registry
	shm         *client.Shm
	compositor  *client.Compositor
	xdgWmBase   *xdg_shell.WmBase
	seat        *client.Seat
	keyboard    *client.Keyboard
	keymap      *xkbcommon.KeyMap
	cursorTheme *cursor.Theme
}

func (s *session) init() {
	var err error

	s.display, err = client.Connect("")
	must(err, "cannot connect to Wayland display")
	s.display.AddErrorHandler(func(err client.DisplayErrorEvent) {
		must(errors.New(err.Message), "wayland display error")
	})

	s.registry, err = s.display.GetRegistry()
	must(err, "cannot get Wayland display registry")
	s.registry.AddGlobalHandler(s.handleRegistryGlobalEvent)

	// TODO: replace w/ routine that checks fields and exit on timeout.
	s.displayRoundTrip()
	s.displayRoundTrip()

	s.cursorTheme, err = cursor.LoadTheme("", 24, s.shm)
	must(err, "cannot load Wayland cursor theme")
}

func (s *session) attachKeyboard() {
	var err error

	s.keyboard, err = s.seat.GetKeyboard()
	must(err, "cannot get Wayland keyboard")

	s.keyboard.AddKeymapHandler(s.handleKeyboardKeymap)
	s.keyboard.AddModifiersHandler(s.handleKeyboardModifiers)
}

func (s *session) detachKeyboard() error {
	if s.keyboard == nil {
		return nil
	}

	err := s.keyboard.Release()
	s.keyboard = nil

	return err
}

func (s *session) displayRoundTrip() {
	c, err := s.display.Sync()
	must(err, "cannot sync Wayland display")
	defer c.Destroy()

	var done bool
	c.AddDoneHandler(func(client.CallbackDoneEvent) {
		done = true
	})

	for !done {
		s.dispatch()
	}
}

func (s *session) destroy() {
	errors := []error{
		s.cursorTheme.Destroy(),
		s.detachKeyboard(),
		s.seat.Release(),
		s.xdgWmBase.Destroy(),
		s.shm.Destroy(),
		s.compositor.Destroy(),
		s.registry.Destroy(),
		s.display.Destroy(),
		s.display.Context().Close(),
	}

	for _, err := range errors {
		if err != nil {
			fyne.LogError("Wayland session close error:", err)
		}
	}
}

func (s *session) mustBind(ev client.RegistryGlobalEvent, id client.Proxy) {
	err := s.registry.Bind(ev.Name, ev.Interface, ev.Version, s.compositor)
	must(err, "cannot bind into Wayland registry interface "+ev.Interface)
}

func (s *session) dispatch() bool {
	err := s.display.Context().Dispatch()
	if !errors.Is(err, net.ErrClosed) {
		must(err, "cannot dispatch Wayland")
	}
	return err == nil
}

func (s *session) handleRegistryGlobalEvent(ev client.RegistryGlobalEvent) {
	switch ev.Interface {
	case "wl_compositor":
		s.compositor = client.NewCompositor(s.display.Context())
		s.mustBind(ev, s.compositor)
	case "wl_shm":
		s.shm = client.NewShm(s.display.Context())
		s.mustBind(ev, s.shm)
	case "xdg_wm_base":
		s.xdgWmBase = xdg_shell.NewWmBase(s.display.Context())
		s.mustBind(ev, s.xdgWmBase)
	case "wl_seat":
		s.seat = client.NewSeat(s.display.Context())
		s.mustBind(ev, s.seat)
		s.seat.AddCapabilitiesHandler(s.handleSeatCapabilities)
	}
}

func (s *session) handleWmBasePing(ev xdg_shell.WmBasePingEvent) {
	err := s.xdgWmBase.Pong(ev.Serial)
	must(err, "cannot ping Wayland server")
}

func (s *session) handleSeatCapabilities(ev client.SeatCapabilitiesEvent) {
	has := func(cap client.SeatCapability) bool {
		return ev.Capabilities&uint32(cap) != 0
	}

	if has(client.SeatCapabilityKeyboard) {
		s.attachKeyboard()
	} else {
		s.detachKeyboard()
	}

	// TODO pointer
	if has(client.SeatCapabilityPointer) {
	} else {
	}
}

func (s *session) handleKeyboardKeymap(ev client.KeyboardKeymapEvent) {
	f := os.NewFile(ev.Fd, "key.map")
	defer f.Close()

	data := make([]byte, ev.Size)

	_, err := io.ReadFull(f, data)
	must(err, "cannot read XKB keymap from fd")

	// https://wayland-book.com/seat/keyboard.html#key-maps
	switch ev.Format {
	case 0: // no_keymap
		return
	case 1:
		s.keymap, err = xkbcommon.Parse(data)
		must(err, "cannot parse keymap from Wayland server")
	default:
		log.Panicln("server sent unknown keymap format", ev.Format)
	}
}

func (s *session) handleKeyboardModifiers(ev client.KeyboardModifiersEvent) {
	if s.keymap == nil {
		return
	}
	s.keymap.UpdateMask(ev.ModsDepressed, ev.ModsLatched, ev.ModsLocked, 0, 0, ev.Group)
}

func must(err error, desc string) {
	if err != nil {
		log.Panicln(desc+":", err)
	}
}
