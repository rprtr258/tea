package spinner

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

type model struct {
	spinner  spinner.Model
	quitting bool
}

func initialModel() *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &model{spinner: s}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(m.spinner.CmdTick)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			f(tea.Quit)
		}
	default:
		f(m.spinner.Update(msg)...)
	}
}

func (m *model) View(r tea.Renderer) {
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		str += "\n"
	}

	r.Write(str)
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
