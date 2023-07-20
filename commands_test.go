package tea

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stringMsg string

func (stringMsg) isBubbleteaMsg() {}

func TestEvery(t *testing.T) {
	expected := stringMsg("every ms")
	msg := Every(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}

func TestTick(t *testing.T) {
	expected := stringMsg("tick")
	msg := Tick(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}

type errorMsg struct {
	MsgImplementation
	error
}

func (errorMsg) isBubbleteaMsg() {}

func TestSequentially(t *testing.T) {
	expectedErrMsg := errorMsg{error: errors.New("some err")}
	expectedStrMsg := stringMsg("some msg")

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
					return expectedErrMsg
				},
				nilReturnCmd,
			},
			expected: expectedErrMsg,
		},
		"some msg": {
			cmds: []Cmd{
				nilReturnCmd,
				func() Msg {
					return expectedStrMsg
				},
				nilReturnCmd,
			},
			expected: expectedStrMsg,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, Sequentially(test.cmds...)())
		})
	}
}

func TestBatch(t *testing.T) {
	t.Run("nil cmd", func(t *testing.T) {
		assert.Nil(t, Batch(nil))
	})
	t.Run("empty cmd", func(t *testing.T) {
		assert.Nil(t, Batch())
	})
	t.Run("single cmd", func(t *testing.T) {
		assert.Len(t, Batch(Quit)(), 1)
	})
	t.Run("mixed nil cmds", func(t *testing.T) {
		assert.Len(t, Batch(nil, Quit, nil, Quit, nil, nil)(), 2)
	})
}
