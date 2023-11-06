package split_editors //nolint:revive,stylecheck

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/textarea"
	"github.com/rprtr258/tea/styles"
)

const (
	initialInputs = 2
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 5
)

var (
	cursorStyle = styles.Style{}.Foreground(styles.FgColor("212"))

	cursorLineStyle = styles.Style{}.
			Background(styles.BgColor("57")).
			Foreground(styles.FgColor("230"))

	placeholderStyle = styles.Style{}.
				Foreground(styles.FgColor("238"))

	endOfBufferStyle = styles.Style{}.
				Foreground(styles.FgColor("235"))

	focusedPlaceholderStyle = styles.Style{}.
				Foreground(styles.FgColor("99"))

	focusedBorderStyle = styles.Style{}
	// Border(styles.RoundedBorder).
	// BorderForeground(styles.FgColor("238"))

	blurredBorderStyle = styles.Style{}
	// Border(styles.HiddenBorder)
)

type keymap = struct {
	next, prev, add, remove, quit key.Binding
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.Binding{Keys: []string{"down"}}
	t.KeyMap.LinePrevious = key.Binding{Keys: []string{"up"}}
	t.Blur()
	return t
}

type model struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int
}

func newModel() *model {
	m := &model{
		inputs: make([]textarea.Model, initialInputs),
		help:   help.New(),
		keymap: keymap{
			next: key.Binding{
				Keys: []string{"tab"},
				Help: key.Help{"tab", "next"},
			},
			prev: key.Binding{
				Keys: []string{"shift+tab"},
				Help: key.Help{"shift+tab", "prev"},
			},
			add: key.Binding{
				Keys: []string{"ctrl+n"},
				Help: key.Help{"ctrl+n", "add an editor"},
			},
			remove: key.Binding{
				Keys: []string{"ctrl+w"},
				Help: key.Help{"ctrl+w", "remove an editor"},
			},
			quit: key.Binding{
				Keys: []string{"esc", "ctrl+c"},
				Help: key.Help{"esc", "quit"},
			},
		},
	}
	for i := 0; i < initialInputs; i++ {
		m.inputs[i] = newTextarea()
	}
	m.inputs[m.focus].Focus()
	m.updateKeybindings()
	return m
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textarea.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			f(tea.Quit)
			return
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			f(m.inputs[m.focus].Focus()...)
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			f(m.inputs[m.focus].Focus()...)
		case key.Matches(msg, m.keymap.add):
			m.inputs = append(m.inputs, newTextarea())
		case key.Matches(msg, m.keymap.remove):
			m.inputs = m.inputs[:len(m.inputs)-1]
			if m.focus > len(m.inputs)-1 {
				m.focus = len(m.inputs) - 1
			}
		}
	case tea.MsgWindowSize:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.updateKeybindings()
	m.sizeInputs()

	// Update all textareas
	for i := range m.inputs {
		m.inputs[i].Update(msg, f)
	}
}

func (m *model) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

func (m *model) updateKeybindings() {
	m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m *model) View(vb tea.Viewbox) {
	for i := range m.inputs {
		m.inputs[i].View(vb.Row(i))
	}

	m.help.ShortHelpView(vb.Padding(tea.PaddingOptions{Top: 2 + len(m.inputs)}), []key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.add,
		m.keymap.remove,
		m.keymap.quit,
	})
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).WithAltScreen().Run()
	return err
}
