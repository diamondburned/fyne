package wayland

type event interface{ ev() }

type (
	windowClosedEvent struct{ w *window }
	quitEvent         struct{ done chan struct{} }
)

func (windowClosedEvent) ev() {}
func (quitEvent) ev()         {}

type waylandLoop struct {
	state *waylandState

	ev chan event
}

func newWaylandLoop(s *waylandState) *waylandLoop {
	l := waylandLoop{
		ev:      make(chan event, 1),
		windows: make(map[*window]struct{}),
	}

	return &l
}

func (l *waylandLoop) run() {
	go func() {
		l.state.dispatch()
	}()
	for ev := range l.ev {
		switch ev := ev.(type) {
		case windowClosedEvent:
			delete(l.windows, ev.w)
			continue
		case quitEvent:
			return
		}
	}

	l.state.quit()
}
