package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/rprtr258/assert"
)

func TestHandmadeFSM(t *testing.T) {
	type event rune
	type state = State[event, struct{}]
	var stateA, stateB, stateC, stateFail state
	stateStart := State[event, struct{}](func(ev event) (state, struct{}) {
		if ev == 'a' {
			return stateA, struct{}{}
		}
		return stateFail, struct{}{}
	})
	stateA = State[event, struct{}](func(ev event) (state, struct{}) {
		if ev == 'b' {
			return stateA, struct{}{}
		}
		return stateFail, struct{}{}
	})
	stateB = State[event, struct{}](func(ev event) (state, struct{}) {
		if ev == 'c' {
			return stateC, struct{}{}
		}
		return stateFail, struct{}{}
	})

	type testcase struct {
		input         string
		expectedState state
	}
	assert.Table(t, map[string]testcase{
		"empty string":            {"", stateStart},
		"non matching string":     {"x", stateFail},
		"matching string":         {"abc", stateC},
		"partial matching string": {"ab", stateB},
	}, func(t *testing.T, test testcase) {
		state := stateStart
		for _, r := range test.input {
			state, _ = state.Handle(event(r))
		}
		assert.Equal(t, test.expectedState, state)
	})
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
	type cmd = []Sum2[cmdStartLoadingAnimation, cmdDisplayPopup]

	var newCmd = New2[cmdStartLoadingAnimation, cmdDisplayPopup]()

	type state = State[event, cmd]
	var stateLoading func(startedAt time.Time) state
	var stateLoaded func(startedAt time.Time, loadedAt time.Time) state

	stateInitial := State[event, cmd](func(ev event) (s state, cmds cmd) {
		ev(eventh{
			A: func(ev startedLoading) {
				s = stateLoading(ev.At)
				cmds = cmd{newCmd.A(cmdStartLoadingAnimation{})}
			},
		})
		return
	})
	stateLoading = func(startedAt time.Time) state {
		return State[event, cmd](func(ev event) (s state, cmds cmd) {
			ev(eventh{
				B: func(ev finishedLoading) {
					s = stateLoaded(startedAt, ev.At)
					cmds = cmd{newCmd.B(cmdDisplayPopup{fmt.Sprintf(
						`Loading finished in %d milliseconds!`,
						ev.At.Sub(startedAt)/time.Millisecond,
					)})}
				},
			})
			return
		})
	}
	stateLoaded = func(time.Time, time.Time) state {
		return State[event, cmd](func(ev event) (s state, cmds cmd) {
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
