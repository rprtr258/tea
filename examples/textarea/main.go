package textarea

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"context"
	"fmt"

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

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{textarea.Blink}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return []tea.Cmd{tea.Quit}
		default:
			if !m.textarea.Focused() {
				cmds = append(cmds, m.textarea.Focus()...)
			}
		}
	}

	return append(cmds, m.textarea.Update(msg)...)
}

func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf(
		"Tell me a story.\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
