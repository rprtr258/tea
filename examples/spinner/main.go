package spinner

// A simple program demonstrating spinner component.

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

type model struct {
	spinner spinner.Model
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
			f(tea.Quit)
		}
	default:
		f(m.spinner.Update(msg)...)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine(2, 1, m.spinner.View()+" Loading forever...press q to quit")
	// if m.quitting {
	// 	r.Write("\n")
	// }
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
