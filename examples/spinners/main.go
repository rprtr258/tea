package spinners

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
		spinner.Circle,
	}

	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	index   int
	spinner spinner.Model
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(m.spinner.CmdTick)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			f(tea.Quit)
		case "h", "left":
			m.index--
			if m.index < 0 {
				m.index = len(spinners) - 1
			}
			m.resetSpinner()
			f(m.spinner.CmdTick)
		case "l", "right":
			m.index++
			if m.index >= len(spinners) {
				m.index = 0
			}
			m.resetSpinner()
			f(m.spinner.CmdTick)
		}
	case spinner.MsgTick:
		f(m.spinner.Update(msg)...)
	}
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New(
		spinner.WithSpinner(spinners[m.index]),
		spinner.WithStyle(spinnerStyle),
	)
}

func (m *model) View(r tea.Renderer) {
	var gap string
	if m.index != 1 {
		gap = " "
	}

	r.Write(fmt.Sprintf(
		"\n %s%s%s\n\n%s",
		m.spinner.View(),
		gap,
		textStyle("Spinning..."),
		helpStyle("h/l, ←/→: change spinner • q: exit\n"),
	))
}

func Main(ctx context.Context) error {
	m := &model{}
	m.resetSpinner()

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
