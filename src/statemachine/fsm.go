package statemachine

import (
	"fmt"
	"strings"
)

// FSM is the state machine that holds the current state.
//
// It has to be created with NewFSM to function properly.
type FSM struct {
	// current is the state that the FSM is currently in.
	current string

	// transitions maps events and source states to destination states.
	transitions map[eKey]string

	// callbacks maps events and targers to callback functions.
	callbacks map[cKey]Callback

	// transition is the internal transition functions used either directly
	// or when Transition is called in an asynchronous state transition.
	transition func()
}

// Event is the info that get passed as a reference in the callbacks.
type Event struct {
	// FSM is a reference to the current FSM.
	FSM *FSM

	// Event is the event name.
	Event string

	// Src is the state before the transition.
	Src string

	// Dst is the state after the transition.
	Dst string

	// Err is an optional error that can be returned from a callback.
	Err error

	// Args is a optinal list of arguments passed to the callback.
	Args []interface{}

	// canceled is an internal flag set if the transition is canceled.
	canceled bool

	// async is an internal flag set if the transition should be asynchronous
	async bool
}

// Cancel can be called in before_<EVENT> or leave_<STATE> to cancel the
// current transition before it happens.
func (e *Event) Cancel() {
	e.canceled = true
}

// Async can be called in leave_<STATE> to do an asynchronous state transition.
//
// The current state transition will be on hold in the old state until a final
// call to Transition is made. This will comlete the transition and possibly
// call the other callbacks.
func (e *Event) Async() {
	e.async = true
}

// Events is a shorthand for defining the transition map in NewFSM.
type Events []EventDesc

// EventDesc represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type EventDesc struct {
	// Name is the event name used when calling for a transition.
	Name string

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []string

	// Dst is the destination state that the FSM will be in if the transition
	// succeds.
	Dst string
}

// Callbacks is a shorthand for defining the callbacks in NewFSM.a
type Callbacks map[string]Callback

// Callback is a function type that callbacks should use. Event is the current
// event info as the callback happens.
type Callback func(*Event)

// Type of callback used in the internal mapping to functions.
type callbackType int

const (
	noCallback callbackType = iota
	beforeEvent
	leaveState
	enterState
	afterEvent
)

// cKey is a struct key used for keeping the callbacks mapped to a target.
type cKey struct {
	// target is either the name of a state or an event depending on which
	// callback type the key refers to. It can also be "" for a non-targeted
	// callback like before_event.
	target string

	// callbackType is the situation when the callback will be run.
	callbackType callbackType
}

// eKey is a struct key used for storing the transition map.
type eKey struct {
	// event is the name of the event that the keys refers to.
	event string

	// src is the source from where the event can transition.
	src string
}

// NewFSM constructs a FSM from events and callbacks.
//
// The events and transitions are specified as a slice of Event structs
// specified as Events. Each Event is mapped to one or more internal
// transitions from Event.Src to Event.Dst.
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after eftering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
//
// There are also two short form versions for the most commonly used callbacks.
// They are simply the name of the event or state:
//
// 1. <NEW_STATE> - called after entering <NEW_STATE>
//
// 2. <EVENT> - called after event named <EVENT>
//
// If both a shorthand version and a full version is specified it is undefined
// which version of the callback will end up in the internal map. This is due
// to the psuedo random nature of Go maps. No checking for multiple keys is
// currently performed.
func NewFSM(initial string, events Events, callbacks Callbacks) *FSM {
	var f FSM
	f.current = initial
	f.transitions = make(map[eKey]string)
	f.callbacks = make(map[cKey]Callback)

	// Build transition map and store sets of all events and states.
	allEvents := make(map[string]bool)
	allStates := make(map[string]bool)
	for _, e := range events {
		for _, src := range e.Src {
			f.transitions[eKey{e.Name, src}] = e.Dst
			allStates[src] = true
			allStates[e.Dst] = true
		}
		allEvents[e.Name] = true
	}

	// Map all callbacks to events/states.
	for name, c := range callbacks {
		var target string
		var cType callbackType

		switch {
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "event" {
				target = ""
				cType = beforeEvent
			} else if _, ok := allEvents[target]; ok {
				cType = beforeEvent
			}
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "state" {
				target = ""
				cType = leaveState
			} else if _, ok := allStates[target]; ok {
				cType = leaveState
			}
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "state" {
				target = ""
				cType = enterState
			} else if _, ok := allStates[target]; ok {
				cType = enterState
			}
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "event" {
				target = ""
				cType = afterEvent
			} else if _, ok := allEvents[target]; ok {
				cType = afterEvent
			}
		default:
			target = name
			if _, ok := allStates[target]; ok {
				cType = enterState
			} else if _, ok := allEvents[target]; ok {
				cType = afterEvent
			}
		}

		if cType != noCallback {
			f.callbacks[cKey{target, cType}] = c
		}
	}

	return &f
}

// Current returns the current state of the FSM.
func (f *FSM) Current() string {
	return f.current
}

// Is returns true if state is the current state.
func (f *FSM) Is(state string) bool {
	return state == f.current
}

// Can returns true if event can occur in the current state.
func (f *FSM) Can(event string) bool {
	_, ok := f.transitions[eKey{event, f.current}]
	return ok && (f.transition == nil)
}

// Can returns true if event can not occure in the current state.
// It is a convenience method to help code read nicely.
func (f *FSM) Cannot(event string) bool {
	return !f.Can(event)
}

// Event initiates a state transition with the named event.
//
// The call takes a variable number of arguments that will be passed to the
// callback, if defined.
//
// It will return nil if the state change is ok or one of these errors:
//
// - event X inappropriate because previous transition did not complete
//
// - event X inappropriate in current state Y
//
// - event X does not exist
//
// - internal error on state transition
//
// The last error should never occur in this situation and is a sign of an
// internal bug.
func (f *FSM) Event(event string, args ...interface{}) error {
	if f.transition != nil {
		return fmt.Errorf("event %s inappropriate because previous transition did not complete", event)
	}

	dst, ok := f.transitions[eKey{event, f.current}]
	if !ok {
		found := false
		for ekey, _ := range f.transitions {
			if ekey.event == event {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("event %s inappropriate in current state %s", event, f.current)
		} else {
			return fmt.Errorf("event %s does not exist", event)
		}
	}

	if f.current == dst {
		return nil
	}

	e := &Event{f, event, f.current, dst, nil, args, false, false}

	// Call the before_ callbacks, first the named then the general version.
	if fn, ok := f.callbacks[cKey{event, beforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return e.Err
		}
	}
	if fn, ok := f.callbacks[cKey{"", beforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return e.Err
		}
	}

	f.transition = func() {
		// Do the state transition.
		f.current = dst

		// Call the enter_ callbacks, first the named then the general version.
		if fn, ok := f.callbacks[cKey{f.current, enterState}]; ok {
			fn(e)
		}
		if fn, ok := f.callbacks[cKey{"", enterState}]; ok {
			fn(e)
		}

		// Call the after_ callbacks, first the named then the general version.
		if fn, ok := f.callbacks[cKey{event, afterEvent}]; ok {
			fn(e)
		}
		if fn, ok := f.callbacks[cKey{"", afterEvent}]; ok {
			fn(e)
		}
	}

	// Call the leave_ callbacks, first the named then the general version.
	if fn, ok := f.callbacks[cKey{f.current, leaveState}]; ok {
		fn(e)
		if e.canceled {
			f.transition = nil
			return e.Err
		} else if e.async {
			return e.Err
		}
	}
	if fn, ok := f.callbacks[cKey{"", leaveState}]; ok {
		fn(e)
		if e.canceled {
			f.transition = nil
			return e.Err
		} else if e.async {
			return e.Err
		}
	}

	// Perform the rest of the transition, if not asynchronous.
	err := f.Transition()
	if err != nil {
		return fmt.Errorf("internal error on state transition")
	}

	return e.Err
}

// Transition completes an asynchrounous state change.
//
// The callback for leave_<STATE> must prviously have called Async on its
// event to have initiated an asynchronous state transition.
func (f *FSM) Transition() error {
	if f.transition == nil {
		return fmt.Errorf("transition inappropriate because no state change in progress")
	}
	f.transition()
	f.transition = nil
	return nil
}
