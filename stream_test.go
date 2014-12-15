package observer

import "testing"

func TestStreamInitialValue(t *testing.T) {
	state1 := newState(10)
	stream1 := &stream{state: state1}
	if val := stream1.Value(); val != 10 {
		t.Fatalf("Expecting 10 but got %#v\n", val)
	}
}

func TestStreamUpdate(t *testing.T) {
	state1 := newState(10)
	stream1 := &stream{state: state1}
	select {
	case <-stream1.Changes():
		t.Fatalf("Expecting no changes\n")
	default:
	}
	state2 := state1.update(15)
	stream2 := &stream{state: state2}
	select {
	case <-stream2.Changes():
		t.Fatalf("Expecting no changes\n")
	default:
	}
	if val := stream2.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
	select {
	case <-stream1.Changes():
	default:
		t.Fatalf("Expecting changes\n")
	}
	stream1.Next()
	if val := stream1.Value(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
}

func TestStreamNextReturnsValue(t *testing.T) {
	state1 := newState(10)
	stream1 := &stream{state: state1}
	state1.update(15)
	if val := stream1.Next(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
}

func TestStreamDetectsChanges(t *testing.T) {
	state1 := newState(10)
	stream1 := &stream{state: state1}
	if has := stream1.HasNext(); has {
		t.Fatalf("Expecting no changes\n")
	}
	state1.update(15)
	if has := stream1.HasNext(); !has {
		t.Fatalf("Expecting changes\n")
	}
	if val := stream1.Next(); val != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", val)
	}
}

func TestStreamWaitsNext(t *testing.T) {
	state1 := newState(10)
	stream1 := &stream{state: state1}
	for i := 15; i <= 100; i++ {
		state1 = state1.update(i)
		if val := stream1.WaitNext(); val != i {
			t.Fatalf("Expecting %#v but got %#v\n", i, val)
		}
	}
	if has := stream1.HasNext(); has {
		t.Fatalf("Expecting no changes\n")
	}
}
