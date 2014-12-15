package observer

type state struct {
	value interface{}
	next  *state
	done  chan struct{}
}

func newState(value interface{}) *state {
	return &state{
		value: value,
		done:  make(chan struct{}),
	}
}

func (s *state) update(value interface{}) *state {
	s.next = newState(value)
	close(s.done)
	return s.next
}
