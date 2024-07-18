package state

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandmadeFSM(t *testing.T) {
	type State int
	type Event rune
	startState := State(0)
	stateA := State(1)
	stateB := State(2)
	stateC := State(3)

	transitions := Transitions[State, Event]{
		startState: map[Event]State{'a': stateA},
		stateA:     map[Event]State{'b': stateB},
		stateB:     map[Event]State{'c': stateC},
	}

	for name, test := range map[string]struct {
		input          string
		expectedStatus Status
	}{
		"empty string":            {"", StatusNormal},
		"non matching string":     {"x", StatusFail},
		"matching string":         {"abc", StatusSuccess},
		"partial matching string": {"ab", StatusNormal},
	} {
		t.Run(name, func(t *testing.T) {
			_, status := FSM(startState, transitions, []Event(test.input))
			require.Equal(t, test.expectedStatus, status)
		})
	}
}
