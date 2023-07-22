package tea

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type msgString string

func TestEvery(t *testing.T) {
	expected := msgString("every ms")
	msg := Every(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}

func TestTick(t *testing.T) {
	expected := msgString("tick")
	msg := Tick(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}

type msgError struct{ err error }

func TestSequentially(t *testing.T) {
	expectedMsgError := msgError{errors.New("some err")}
	expectedMsgString := msgString("some msg")

	nilReturnCmd := func() Msg {
		return nil
	}

	for name, test := range map[string]struct {
		cmds     []Cmd
		expected Msg
	}{
		"all nil": {
			cmds:     []Cmd{nilReturnCmd, nilReturnCmd},
			expected: nil,
		},
		"null cmds": {
			cmds:     []Cmd{nil, nil},
			expected: nil,
		},
		"one error": {
			cmds: []Cmd{
				nilReturnCmd,
				func() Msg {
					return expectedMsgError
				},
				nilReturnCmd,
			},
			expected: expectedMsgError,
		},
		"some msg": {
			cmds: []Cmd{
				nilReturnCmd,
				func() Msg {
					return expectedMsgString
				},
				nilReturnCmd,
			},
			expected: expectedMsgString,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, Sequentially(test.cmds...)())
		})
	}
}
