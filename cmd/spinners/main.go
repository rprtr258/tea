package spinners

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/styles"
)

var (
	// Available spinners
	spinners = [...]spinner.Spinner{
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

	textStyle    = styles.Style{}.Foreground(styles.FgColor("252")).Render
	spinnerStyle = styles.Style{}.Foreground(styles.FgColor("69"))
	helpStyle    = styles.Style{}.Foreground(styles.FgColor("241")).Render
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
		m.spinner.Update(msg, f)
	}
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New(
		spinner.WithSpinner(spinners[m.index]),
		spinner.WithStyle(spinnerStyle),
	)
}

func (m *model) View(vb tea.Viewbox) {
	m.spinner.View(vb.PaddingLeft(1))
	vb = vb.PaddingLeft(5 + 1)
	vb.WriteLine(textStyle("Spinning..."))
	vb.PaddingTop(2).WriteLine(helpStyle("h/l, ←/→: change spinner • q: exit"))
}

func Main(ctx context.Context) error {
	m := &model{}
	m.resetSpinner()

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
