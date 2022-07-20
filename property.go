package observer

import "sync"

// Property is an object that is continuously updated by one or more
// publishers. It is completely goroutine safe: you can use Property
// concurrently from multiple goroutines.
type Property[T any] interface {
	// Value returns the current value for this property.
	Value() T

	// Update sets a new value for this property.
	Update(value T)

	// Observe returns a newly created Stream for this property.
	Observe() Stream[T]
}

// NewProperty creates a new Property with the initial value value.
// It returns the created Property.
func NewProperty[T any](value T) Property[T] {
	return &property[T]{state: newState(value)}
}

type property[T any] struct {
	sync.RWMutex
	state *state[T]
}

func (p *property[T]) Value() T {
	p.RLock()
	defer p.RUnlock()
	return p.state.value
}

func (p *property[T]) Update(value T) {
	p.Lock()
	defer p.Unlock()
	p.state = p.state.update(value)
}

func (p *property[T]) Observe() Stream[T] {
	p.RLock()
	defer p.RUnlock()
	return &stream[T]{state: p.state}
}
