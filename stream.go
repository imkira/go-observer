package observer

// Stream represents the list of values a property is updated to.  For every
// property update, that value is appended to the list in the order they
// happen. The value is discarded once you advance the stream.  Please note
// that Stream is not goroutine safe: you cannot use the same stream on
// multiple goroutines concurrently. If you want to use multiple streams for
// the same property, either use Property.Observe (goroutine-safe) or use
// Stream.Clone (before passing it to another goroutine).
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

	// Clone creates a new independent stream from this one but sharing the same
	// Property. Updates to the property will be reflected in both streams but
	// they may have different values depending on when they advance the stream
	// with Next.
	Clone() Stream
}

type stream struct {
	state *state
}

func (s *stream) Clone() Stream {
	return &stream{state: s.state}
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
	<-s.state.done
	s.state = s.state.next
	return s.state.value
}
