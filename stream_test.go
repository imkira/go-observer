package observer

import "testing"

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
