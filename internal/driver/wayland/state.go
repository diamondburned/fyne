package wayland

import (
	"sync"
	"sync/atomic"
)

// TODO: refactor this into an event loop instead of mutex. This is horrible.

type state struct {
	mu      sync.Mutex
	windows map[*window]struct{}
	session session
	running uint32
}

func newWaylandState() *state {
	return &state{
		windows: make(map[*window]struct{}),
	}
}

func (s *state) use(f func(s *session)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.session == (session{}) {
		s.session.init()
	}

	f(&s.session)
}

func (s *state) removeWindow(w *window) {
	s.mu.Lock()
	delete(s.windows, w)
	s.mu.Unlock()
}

func (s *state) run() {
	if !atomic.CompareAndSwapUint32(&s.running, 0, 1) {
		panic("run called more than once")
	}
	defer atomic.StoreUint32(&s.running, 0)

	s.use(func(*session) {})

	done := false
	for s.session.dispatch() && !done {
		s.mu.Lock()
		done = len(s.windows) == 0
		s.mu.Unlock()
	}
}

func (s *state) quit() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for w := range s.windows {
		w.Close()
		<-w.dead // TODO: potential deadlock here
	}

	if s.session != (session{}) {
		s.session.destroy()
	}
}
