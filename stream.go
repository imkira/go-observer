package observer

import "context"

// Stream represents the list of values a property is updated to.  For every
// property update, that value is appended to the list in the order they
// happen. The value is discarded once you advance the stream.  Please note
// that Stream is not goroutine safe: you cannot use the same stream on
// multiple goroutines concurrently. If you want to use multiple streams for
// the same property, either use Property.Observe (goroutine-safe) or use
// Stream.Clone (before passing it to another goroutine).
type Stream[T any] interface {
	// Value returns the current value for this stream.
	Value() T

	// Changes returns the channel that is closed when a new value is available.
	Changes() chan struct{}

	// Next advances this stream to the next state.
	// You should never call this unless Changes channel is closed.
	Next() T

	// HasNext checks whether there is a new value available.
	HasNext() bool

	// WaitNext waits for Changes to be closed, advances the stream and returns
	// the current value.
	WaitNext() T

	// WaitNextFiltered does the same as WaitNext but only returns when filterFunc evaluates to true.
	WaitNextFiltered(filterFunc func(T) bool) T

	// WaitNextCtx does the same as WaitNext but returns earlier with an error if the given context is cancelled first.
	WaitNextCtx(ctx context.Context) (T, error)

	// WaitNextCtxFiltered does the same as WaitNextCtx and additionally only returns the value when filterFunc evaluates to true.
	WaitNextCtxFiltered(ctx context.Context, filterFunc func(T) bool) (T, error)

	// Clone creates a new independent stream from this one but sharing the same
	// Property. Updates to the property will be reflected in both streams but
	// they may have different values depending on when they advance the stream
	// with Next.
	Clone() Stream[T]

	// Peek return the value in the next state
	// You should never call this unless Changes channel is closed.
	Peek() T
}

type stream[T any] struct {
	state *state[T]
}

func (s *stream[T]) Clone() Stream[T] {
	return &stream[T]{state: s.state}
}

func (s *stream[T]) Value() T {
	return s.state.value
}

func (s *stream[T]) Changes() chan struct{} {
	return s.state.done
}

func (s *stream[T]) Next() T {
	s.state = s.state.next
	return s.state.value
}

func (s *stream[T]) HasNext() bool {
	select {
	case <-s.state.done:
		return true
	default:
		return false
	}
}

func (s *stream[T]) WaitNext() T {
	<-s.state.done
	s.state = s.state.next
	return s.state.value
}

func (s *stream[T]) WaitNextFiltered(filterFunc func(T) bool) T {
	for {
		val := s.WaitNext()
		if evaluateFilterFunc[T](val, filterFunc) {
			return val
		}
	}
}

func (s *stream[T]) WaitNextCtx(ctx context.Context) (T, error) {
	select {
	case <-s.Changes():
		// ensure that context is not canceled, only then advance the stream
		if ctx.Err() == nil {
			return s.Next(), nil
		}

	case <-ctx.Done():
	}

	var zeroVal T
	return zeroVal, ctx.Err()
}

func (s *stream[T]) WaitNextCtxFiltered(ctx context.Context, filterFunc func(T) bool) (T, error) {
	for {
		val, err := s.WaitNextCtx(ctx)
		if err != nil {
			return val, err
		}
		if evaluateFilterFunc[T](val, filterFunc) {
			return val, nil
		}
	}
}

func (s *stream[T]) Peek() T {
	return s.state.next.value
}

func evaluateFilterFunc[T any](val T, filterFunc func(T) bool) bool {
	return filterFunc == nil || filterFunc(val)
}
