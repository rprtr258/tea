package chat

// A simple program demonstrating text area component.

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/textarea"
	"github.com/rprtr258/tea/components/viewport"
	"github.com/rprtr258/tea/styles"
)

type model struct {
	viewport    viewport.Model
	lines       [][2]string // [sender, message]
	textarea    textarea.Model
	senderStyle styles.Style
}

func newModel() *model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = styles.Style{} // Remove cursor line styling
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 5)

	return &model{
		textarea:    ta,
		viewport:    vp,
		senderStyle: styles.Style{}.Foreground(styles.FgColor("5")),
		lines: [][2]string{
			{"system", "Welcome to the chat room!"},
			{"system", "Type a message and press Enter to send."},
		},
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textarea.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			f(tea.Quit)
			return
		case tea.KeyEnter:
			m.lines = append(m.lines, [2]string{"You", m.textarea.Value()})
			m.textarea.Reset()
			m.viewport.GotoBottom(len(m.lines))
		}
	default:
		// do not update viewport on keys, they must go to textarea
		m.viewport.Update(msg)
	}

	m.textarea.Update(msg, f)
}

func (m *model) View(vb tea.Viewbox) {
	m.viewport.View(vb, func(vb tea.Viewbox, i int) {
		if i >= len(m.lines) {
			return
		}

		x := vb.Styled(m.senderStyle).WriteLine(m.lines[i][0] + ": ")
		vb.PaddingLeft(x).WriteLine(m.lines[i][1])
	})
	m.textarea.View(vb.PaddingTop(m.viewport.Height))
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
