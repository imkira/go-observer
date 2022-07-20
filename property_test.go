package observer

import (
	"sync"
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
	var prevStream Stream[int]
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

func TestPropertyMultipleConcurrentReaders(t *testing.T) {
	initial := 1000
	final := 2000
	prop := NewProperty(initial)
	var cherrs []chan error
	for i := 0; i < 1000; i++ {
		cherr := make(chan error, 1)
		cherrs = append(cherrs, cherr)
		go testStreamRead(prop.Observe(), initial, final, cherr)
	}
	done := make(chan bool)
	go func(prop Property[int], initial, final int, done chan bool) {
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

func TestPropertyMultipleConcurrentReadersWriters(t *testing.T) {
	wg := &sync.WaitGroup{}
	writer := func(prop Property[int], times int) {
		defer wg.Done()
		for i := 0; i <= times; i++ {
			val := prop.Value()
			prop.Update(val + 1)
			prop.Observe()
		}
	}
	prop := NewProperty(0)
	times := 1000
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go writer(prop, times)
	}
	wg.Wait()
}
