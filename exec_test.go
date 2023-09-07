package tea

import (
	"bytes"
	"context"
	"os/exec"
	"reflect"
	"testing"

	"github.com/rprtr258/assert"
)

type msgExecFinished struct{ err error }

type testExecModel struct {
	cmd string
	err error
}

func (m *testExecModel) Init(f func(...Cmd)) {
	c := exec.Command(m.cmd)
	f(ExecProcess(c, func(err error) Msg {
		return msgExecFinished{err}
	}))
}

func (m *testExecModel) Update(msg Msg, f func(...Cmd)) {
	if msg, ok := msg.(msgExecFinished); ok {
		m.err = msg.err
		f(Quit)
	}
}

func (m *testExecModel) View(Viewbox) {}

func TestTeaExec(t *testing.T) {
	for name, test := range map[string]struct {
		cmd       string
		expectErr error
	}{
		"true": {
			cmd:       "true",
			expectErr: nil,
		},
		"false": {
			cmd:       "false",
			expectErr: &exec.ExitError{},
		},
		"invalid command": {
			cmd:       "invalid",
			expectErr: &exec.Error{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			var in bytes.Buffer

			m := &testExecModel{cmd: test.cmd}
			_, err := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf).Run()
			assert.NoError(t, err)
			assert.Equal(t, reflect.TypeOf(test.expectErr), reflect.TypeOf(m.err))
		})
	}
}
