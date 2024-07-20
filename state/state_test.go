package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func compare4[A, B, C, D comparable](a, b Sum4[A, B, C, D]) (res bool) {
	a(h4[A, B, C, D]{
		A: func(a A) { b(h4[A, B, C, D]{A: func(b A) { res = true }}) },
		B: func(a B) { b(h4[A, B, C, D]{B: func(b B) { res = true }}) },
		C: func(a C) { b(h4[A, B, C, D]{C: func(b C) { res = true }}) },
		D: func(a D) { b(h4[A, B, C, D]{D: func(b D) { res = true }}) },
	})
	return
}

func TestHandmadeFSM(t *testing.T) {
	type stateStart struct{}
	type stateA struct{}
	type stateB struct{}
	type stateC struct{}
	// TODO: stateFail
	type stateh = h4[stateStart, stateA, stateB, stateC]
	type state = Sum4[stateStart, stateA, stateB, stateC]

	var newState = New4[stateStart, stateA, stateB, stateC]()

	type event rune

	type cmd = struct{}

	update := func(state state, ev event) (s state, _ []cmd) {
		state(stateh{
			A: func(state stateStart) {
				if ev == 'a' {
					s = newState.B(stateA{})
				}
			},
			B: func(state stateA) {
				if ev == 'b' {
					s = newState.C(stateB{})
				}
			},
			C: func(stateB) {
				if ev == 'c' {
					s = newState.D(stateC{})
				}
			},
		})
		return
	}

	initialState := newState.A(stateStart{})
	hook := NewHook(initialState, update)

	for name, test := range map[string]struct {
		input         string
		expectedState state
	}{
		"empty string": {"", newState.A(stateStart{})},
		// "non matching string":     {"x", StatusFail},
		"matching string":         {"abc", newState.D(stateC{})},
		"partial matching string": {"ab", newState.C(stateB{})},
	} {
		t.Run(name, func(t *testing.T) {
			hook.Send([]event(test.input)...)
			require.True(t, compare4(test.expectedState, hook.State()))
		})
	}
}

func TestLoadFSM(t *testing.T) {
	type stateInitial struct{}
	type stateLoading struct{ startedAt time.Time }
	type stateLoaded struct {
		stateLoading
		loadedAt time.Time
	}
	type stateh = h3[stateInitial, stateLoading, stateLoaded]
	type state = Sum3[stateInitial, stateLoading, stateLoaded]

	var newState = New3[stateInitial, stateLoading, stateLoaded]()

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

	update := func(state state, ev event) (s state, cmds []cmd) {
		state(stateh{
			A: func(state stateInitial) {
				ev(eventh{
					A: func(ev startedLoading) {
						s = newState.B(stateLoading{startedAt: ev.At})
						cmds = []cmd{newCmd.A(cmdStartLoadingAnimation{})}
					},
				})
			},
			B: func(state stateLoading) {
				ev(eventh{
					B: func(ev finishedLoading) {
						s = newState.C(stateLoaded{state, ev.At})
						cmds = []cmd{newCmd.B(cmdDisplayPopup{fmt.Sprintf(
							`Loading finished in %d milliseconds!`,
							ev.At.Sub(state.startedAt)/time.Millisecond,
						)})}
					},
				})
			},
			C: func(stateLoaded) {},
		})
		return
	}

	initialState := newState.A(stateInitial{})
	hook := NewHook(initialState, update)
	fmt.Println(hook.State())
	hook.Send(newEvent.A(startedLoading{time.Now()}))
	fmt.Println(hook.State())
	hook.Send(newEvent.B(finishedLoading{time.Now()}))
	fmt.Println(hook.State())
	for _, cmd := range hook.cmds { // handle commands
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
