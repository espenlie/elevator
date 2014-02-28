package fsm

import (
	"fmt"
	"testing"
)

func TestMultipleSources(t *testing.T) {
	fsm := NewFSM(
		"one",
		Events{
			{Name: "first", Src: []string{"one"}, Dst: "two"},
			{Name: "second", Src: []string{"two"}, Dst: "three"},
			{Name: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
		},
		Callbacks{},
	)

	fsm.Event("first")
	if fsm.Current() != "two" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.FailNow()
	}
	fsm.Event("first")
	fsm.Event("second")
	if fsm.Current() != "three" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.FailNow()
	}
}
