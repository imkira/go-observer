package observer

import (
	"fmt"
	"testing"
)

func TestPropertyInitialValue(t *testing.T) {
	prop := NewProperty(10)
	if val := prop.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
	stream := prop.Observe()
	for i := 0; i <= 100; i++ {
		if val := stream.Value(); val != 10 {
			t.Fatalf("Expecting 10 but got %#v\n", val)
		}
	}
}

func TestPropertyInitialObserve(t *testing.T) {
	prop := NewProperty(10)
	var prevStream Stream
	for i := 0; i <= 100; i++ {
		stream := prop.Observe()
		if stream == prevStream {
			t.Fatalf("Expecting different stream\n")
		}
		if val := stream.Value(); val != 10 {
			t.Fatalf("Expecting 10 but got %#v\n", val)
		}
		prevStream = stream
	}
}

func TestPropertyObserveAfterUpdate(t *testing.T) {
	prop := NewProperty(10)
	if val := prop.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
	prop.Update(15)
	if val := prop.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	stream := prop.Observe()
	if val := stream.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
}

func TestPropertyMultipleConcurrentObservers(t *testing.T) {
	observer := func(s Stream, initial, final int, err chan error) {
		val := s.Value().(int)
		if val != initial {
			err <- fmt.Errorf("Expecting %#v but got %#v\n", initial, val)
			return
		}
		for i := initial + 1; i <= final; i++ {
			prevVal := val
			val = s.WaitNext().(int)
			expected := prevVal + 1
			if val != expected {
				err <- fmt.Errorf("Expecting %#v but got %#v\n", expected, val)
				return
			}
		}
		close(err)
	}
	initial := 1000
	final := 2000
	prop := NewProperty(initial)
	var cherrs []chan error
	for i := 0; i < 1000; i++ {
		cherr := make(chan error, 1)
		cherrs = append(cherrs, cherr)
		go observer(prop.Observe(), initial, final, cherr)
	}
	done := make(chan bool)
	go func(prop Property, initial, final int, done chan bool) {
		defer close(done)
		for i := initial + 1; i <= final; i++ {
			prop.Update(i)
		}
	}(prop, initial, final, done)
	for _, cherr := range cherrs {
		if err := <-cherr; err != nil {
			t.Fatal(err)
		}
	}
	<-done
}
