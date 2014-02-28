package main

import (
//   "drivers"
    "fmt"
    "statemachine"
)

func TestMultipleSources(t *testing.T) {
	fsm := NewFSM(
		"INIT",
		Events{
			{Name: "IDLE", Src: []string{"INIT", "DOOR_OPEN", "ESTOP"}, Dst: []string{"UP", "DOOR_OPEN", "DOWN", "ESTOP", "ERROR"}},
			{Name: "UP", Src: []string{"IDLE", "DOOR_OPEN"}, Dst: []string{"DOOR_OPEN", "ESTOP", "ERROR"}},
			{Name: "DOWN", Src: []string{"IDLE", "DOOR_OPEN"}, Dst: []string{"DOOR_OPEN", "ESTOP", "ERROR"}},
			{Name: "DOOR_OPEN", Src: []string{"IDLE", "UP", "DOWN"}, Dst: []string{"UP", "DOWN", "ESTOP", "ERROR"}},
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

func main() {
    for {
        fmt.Println("hei")
    }
}
