package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandmadeFSM(t *testing.T) {
	type event rune
	type cmd struct{}
	type state = State[event, cmd]
	var stateA, stateB, stateC, stateFail state
	stateStart := func(ev event) (state, []cmd) {
		if ev == 'a' {
			return stateA, nil
		}
		return stateFail, nil
	})
	stateA = StateFunc[event,cmd](func(ev event) (state, []cmd) {
		if ev == 'b' {
			return stateA, nil
		}
		return stateFail, nil
	})
	stateB = StateFunc[event,cmd](func(ev event) (state, []cmd) {
		if ev == 'c' {
			return stateC, nil
		}
		return stateFail, nil
	})

	for name, test := range map[string]struct {
		input         string
		expectedState state
	}{
		"empty string": {"", stateStart},
		"non matching string":     {"x", StatusFail},
		"matching string":         {"abc", stateC},
		"partial matching string": {"ab", stateB},
	} {
		t.Run(name, func(t *testing.T) {
			state := stateStart
			for _, r := range test.input {
				state, _ = state(event(r))
			}
			require.Equal(t, test.expectedState, state)
		})
	}
}

func TestLoadFSM(t *testing.T) {
	type startedLoading struct{ At time.Time }
	type finishedLoading struct{ At time.Time }
	type eventh = h2[startedLoading, finishedLoading]
	type event = Sum2[startedLoading, finishedLoading]

	var newEvent = New2[startedLoading, finishedLoading]()

	type cmdStartLoadingAnimation struct{}
	type cmdDisplayPopup struct{ text string }
	type cmdh = h2[cmdStartLoadingAnimation, cmdDisplayPopup]
	type cmd = Sum2[cmdStartLoadingAnimation, cmdDisplayPopup]

	var newCmd = New2[cmdStartLoadingAnimation, cmdDisplayPopup]()

	type state = State[event, cmd]
	var stateLoading func(startedAt time.Time) state
	var stateLoaded func(startedAt time.Time, loadedAt time.Time) state

	stateInitial := StateFunc[event, cmd](func(ev event) (s state, cmds []cmd) {
		ev(eventh{
			A: func(ev startedLoading) {
				s = stateLoading(ev.At)
				cmds = []cmd{newCmd.A(cmdStartLoadingAnimation{})}
			},
		})
		return
	})
	stateLoading = func(startedAt time.Time) state {
		return StateFunc[event, cmd](func(ev event) (s state, cmds []cmd) {
			ev(eventh{
				B: func(ev finishedLoading) {
					s = stateLoaded(startedAt, ev.At)
					cmds = []cmd{newCmd.B(cmdDisplayPopup{fmt.Sprintf(
						`Loading finished in %d milliseconds!`,
						ev.At.Sub(startedAt)/time.Millisecond,
					)})}
				},
			})
			return
		})
	}
	stateLoaded = func(startedAt time.Time, loadedAt time.Time) state {
		return StateFunc[event, cmd](func(ev event) (s state, cmds []cmd) {
			return
		})
	}

	{
		initialState := stateInitial
		state := state(initialState)
		fmt.Println(state)
		state, cmds0 := state.Handle(newEvent.A(startedLoading{time.Now()}))
		fmt.Println(state)
		state, cmds1 := state.Handle(newEvent.B(finishedLoading{time.Now()}))
		fmt.Println(state)
		for _, cmd := range append(cmds0, cmds1...) { // handle commands
			cmd(cmdh{
				A: func(cmd cmdStartLoadingAnimation) {
					fmt.Println(`loading animation started`)
				},
				B: func(cmd cmdDisplayPopup) {
					fmt.Println(`displaying popup with text "${cmd.text}"`)
				},
			})
		}
	}
}
