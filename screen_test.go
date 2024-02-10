package tea

import (
	"bytes"
	"context"
	"testing"

	"github.com/rprtr258/assert"
)

type initCmdModel struct{}

func (m *initCmdModel) Init(func(...Cmd)) {}

func (m *initCmdModel) Update(msg Msg, f func(...Cmd)) {
	if _, ok := msg.(MsgKey); ok {
		f(Quit)
	}
}

func (m *initCmdModel) View(vb Viewbox) {
	vb.WriteLine("success")
}

//nolint:lll // uuh
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
			p := NewProgram(context.Background(), &initCmdModel{}).
				WithInput(&in).
				WithOutput(&out)
			ch := make(chan struct{})
			go func() {
				for _, cmd := range test.cmds {
					p.Send(cmd())
				}
				p.Send(Quit())
				ch <- struct{}{}
			}()
			_, err := p.Run()
			assert.NoError(t, err)
			<-ch

			assert.Equal(t, test.expected, out.String())
		})
	}
}
