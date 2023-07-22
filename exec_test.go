package tea

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type msgExecFinished struct{ err error }

type testExecModel struct {
	cmd string
	err error
}

func (m *testExecModel) Init() Cmd {
	c := exec.Command(m.cmd)
	return ExecProcess(c, func(err error) Msg {
		return msgExecFinished{err}
	})
}

func (m *testExecModel) Update(msg Msg) Cmd {
	switch msg := msg.(type) { //nolint:gocritic
	case msgExecFinished:
		m.err = msg.err
		return Quit
	}

	return nil
}

func (m *testExecModel) View(Renderer) {}

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
			assert.IsType(t, test.expectErr, m.err)
		})
	}
}
