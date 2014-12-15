package observer

// Stream represents the list of values a property is updated to.
// For every property update, that value is appended to the list in the order
// they happen. The value is discarded once you advance the stream.
// Please note that Stream is not goroutine safe: if you have multiple
// goroutines in which you want to observe a property, you must create at least
// one Stream for each goroutine.
type Stream interface {
	// Value returns the current value for this stream.
	Value() interface{}

	// Changes returns the channel that is closed when a new value is available.
	Changes() chan struct{}

	// Next advances this stream to the next state.
	// You should never call this unless Changes channel is closed.
	Next() interface{}

	// HasNext checks whether there is a new value available.
	HasNext() bool

	// WaitNext waits for Changes to be closed, advances the stream and returns
	// the current value.
	WaitNext() interface{}

	// Clone creates a new independent stream from this one.
	Clone() Stream
}

type stream struct {
	state *state
}

func (s *stream) Clone() Stream {
	return &stream{state: newState(s.state.value)}
}

func (s *stream) Value() interface{} {
	return s.state.value
}

func (s *stream) Changes() chan struct{} {
	return s.state.done
}

func (s *stream) Next() interface{} {
	s.state = s.state.next
	return s.state.value
}

func (s *stream) HasNext() bool {
	select {
	case <-s.state.done:
		return true
	default:
		return false
	}
}

func (s *stream) WaitNext() interface{} {
	select {
	case <-s.state.done:
		s.state = s.state.next
		return s.state.value
	}
}
