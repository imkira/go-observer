package observer

import (
	"fmt"
	"testing"
	"time"
)

func TestStreamInitialValue(t *testing.T) {
	state := newState(10)
	stream := &stream{state: state}
	if val := stream.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
}

func TestStreamUpdate(t *testing.T) {
	state1 := newState(10)
	state2 := state1.update(15)
	stream := &stream{state: state1}
	if val := stream.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
	state2.update(15)
	if val := stream.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
}

func TestStreamNextValue(t *testing.T) {
	state1 := newState(10)
	stream := &stream{state: state1}
	state2 := state1.update(15)
	if val := stream.Next(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	state2.update(20)
	if val := stream.Next(); val != 20 {
		t.Fatalf("Expecting 20 but got %#v\n", val)
	}
}

func TestStreamDetectsChanges(t *testing.T) {
	state := newState(10)
	stream := &stream{state: state}
	select {
	case <-stream.Changes():
		t.Fatalf("Expecting no changes\n")
	default:
	}
	go func() {
		time.Sleep(1 * time.Second)
		state.update(15)
	}()
	select {
	case <-stream.Changes():
	case <-time.After(2 * time.Second):
		t.Fatalf("Expecting changes\n")
	}
	select {
	case <-stream.Changes():
	default:
		t.Fatalf("Expecting changes\n")
	}
	if val := stream.Next(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	select {
	case <-stream.Changes():
		t.Fatalf("Expecting no changes\n")
	default:
	}
}

func TestStreamHasChanges(t *testing.T) {
	state := newState(10)
	stream := &stream{state: state}
	if stream.HasNext() {
		t.Fatalf("Expecting no changes\n")
	}
	state.update(15)
	if !stream.HasNext() {
		t.Fatalf("Expecting changes\n")
	}
	if val := stream.Next(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
}

func TestStreamWaitsNext(t *testing.T) {
	state := newState(10)
	stream := &stream{state: state}
	for i := 15; i <= 100; i++ {
		state = state.update(i)
		if val := stream.WaitNext(); val != i {
			t.Fatalf("Expecting %#v but got %#v\n", i, val)
		}
	}
	if stream.HasNext() {
		t.Fatalf("Expecting no changes\n")
	}
}

func TestStreamClone(t *testing.T) {
	state := newState(10)
	stream1 := &stream{state: state}
	stream2 := stream1.Clone()
	if stream2.HasNext() {
		t.Fatalf("Expecting no changes\n")
	}
	if val := stream2.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
	state.update(15)
	if !stream1.HasNext() {
		t.Fatalf("Expecting changes\n")
	}
	if !stream2.HasNext() {
		t.Fatalf("Expecting changes\n")
	}
	stream1.Next()
	if val := stream1.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	if val := stream2.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
	stream2.Next()
	if val := stream2.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	if stream1.HasNext() {
		t.Fatalf("Expecting no changes\n")
	}
	if stream2.HasNext() {
		t.Fatalf("Expecting no changes\n")
	}
}

func TestStreamPeek(t *testing.T) {
	state := newState(10)
	stream := &stream{state: state}
	state = state.update(15)
	if val := stream.Peek(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	state = state.update(20)
	if val := stream.Peek(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	stream.Next()
	if val := stream.Peek(); val != 20 {
		t.Fatalf("Expecting 20 but got %#v\n", val)
	}
}

func TestStreamConcurrencyWithClones(t *testing.T) {
	initial := 1000
	final := 2000
	prop := NewProperty(initial)
	stream := prop.Observe()
	var cherrs []chan error
	for i := 0; i < 1000; i++ {
		cherr := make(chan error, 1)
		cherrs = append(cherrs, cherr)
		go testStreamRead(stream.Clone(), initial, final, cherr)
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

func testStreamRead(s Stream, initial, final int, err chan error) {
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
