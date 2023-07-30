package tea

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type initCmdModel struct {
	initCmds []Cmd
}

func (m *initCmdModel) Init() []Cmd {
	return m.initCmds
}

func (m *initCmdModel) Update(msg Msg) []Cmd {
	switch msg.(type) { //nolint:gocritic
	case MsgKey:
		return []Cmd{Quit}
	}

	return nil
}

func (m *initCmdModel) View(r Renderer) {
	r.Write("success\n")
}

//nolint:lll
func TestMsgClear(t *testing.T) {
	for name, test := range map[string]struct {
		cmds     []Cmd
		expected string
	}{
		"clear_screen": {
			cmds:     []Cmd{ClearScreen},
			expected: "\x1b[?25l\x1b[2J\x1b[1;1H\x1b[1;1Hsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"altscreen": {
			cmds:     []Cmd{EnterAltScreen, ExitAltScreen},
			expected: "\x1b[?25l\x1b[?1049h\x1b[2J\x1b[1;1H\x1b[1;1H\x1b[?25l\x1b[?1049l\x1b[?25lsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"altscreen_autoexit": {
			cmds:     []Cmd{EnterAltScreen},
			expected: "\x1b[?25l\x1b[?1049h\x1b[2J\x1b[1;1H\x1b[1;1H\x1b[?25lsuccess\r\n\x1b[2;0H\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l\x1b[?1049l\x1b[?25h",
		},
		"mouse_cellmotion": {
			cmds:     []Cmd{EnableMouseCellMotion},
			expected: "\x1b[?25l\x1b[?1002hsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"mouse_allmotion": {
			cmds:     []Cmd{EnableMouseAllMotion},
			expected: "\x1b[?25l\x1b[?1003hsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"mouse_disable": {
			cmds:     []Cmd{EnableMouseAllMotion, DisableMouse},
			expected: "\x1b[?25l\x1b[?1003h\x1b[?1002l\x1b[?1003lsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"cursor_hide": {
			cmds:     []Cmd{HideCursor},
			expected: "\x1b[?25l\x1b[?25lsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
		"cursor_hideshow": {
			cmds:     []Cmd{HideCursor, ShowCursor},
			expected: "\x1b[?25l\x1b[?25l\x1b[?25hsuccess\r\n\x1b[0D\x1b[2K\x1b[?25h\x1b[?1002l\x1b[?1003l",
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			p := NewProgram(context.Background(), &initCmdModel{append(test.cmds, Quit)}).
				WithInput(&in).
				WithOutput(&out)

			_, err := p.Run()
			assert.NoError(t, err)

			assert.Equal(t, test.expected, out.String())
		})
	}
}
