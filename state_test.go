package observer

import "testing"

func testStateNew(t *testing.T, state *state, val interface{}) {
	if state.value != val {
		t.Fatalf("Expecting %#v but got %#v\n", val, state.value)
	}
	if state.next != nil {
		t.Fatalf("Expecting no next but got %#v\n", state.next)
	}
	select {
	case <-state.done:
		t.Fatalf("Expecting not done\n")
	default:
	}
}

func TestStateNew(t *testing.T) {
	state := newState(10)
	testStateNew(t, state, 10)
}

func TestStateUpdate(t *testing.T) {
	state1 := newState(10)
	state2 := state1.update(15)
	if state1 == state2 {
		t.Fatalf("Expecting different states\n")
	}
	if state2.value != 15 {
		t.Fatalf("Expecting 15 but got %#v\n", state1.value)
	}
	if state1.next == nil {
		t.Fatalf("Expecting next but got %#v\n", state1.next)
	}
	select {
	case <-state1.done:
	default:
		t.Fatalf("Expecting done\n")
	}
	testStateNew(t, state2, 15)
}
