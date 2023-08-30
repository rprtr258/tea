package textarea

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/textarea"
)

type msgErr error

type model struct {
	textarea textarea.Model
}

func initialModel() *model {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.Focus()

	return &model{
		textarea: ti,
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textarea.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			f(tea.Quit)
		default:
			if !m.textarea.Focused() {
				f(m.textarea.Focus()...)
			}
		}
	}

	f(m.textarea.Update(msg)...)
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine(0, 0, "Tell me a story.")
	m.textarea.View(vb.Padding(tea.PaddingOptions{Top: 1, Bottom: 2}))
	vb.WriteLine(2+m.textarea.Height(), 0, "(ctrl+c to quit)")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
