package observer

type state[T any] struct {
	value T
	next  *state[T]
	done  chan struct{}
}

func newState[T any](value T) *state[T] {
	return &state[T]{
		value: value,
		done:  make(chan struct{}),
	}
}

func (s *state[T]) update(value T) *state[T] {
	s.next = newState(value)
	close(s.done)
	return s.next
}
