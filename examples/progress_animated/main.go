package progress_animated

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"context"
	"strings"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/progress"
	"github.com/rprtr258/tea/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

type msgTick time.Time

type model struct {
	progress progress.Model
}

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{tickCmd()}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		return []tea.Cmd{tea.Quit}

	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return nil

	case msgTick:
		if m.progress.Percent() == 1.0 {
			return []tea.Cmd{tea.Quit}
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.25)
		return []tea.Cmd{tickCmd(), cmd}

	// MsgFrame is sent when the progress bar wants to animate itself
	case progress.MsgFrame:
		return m.progress.Update(msg)

	default:
		return nil
	}
}

func (m *model) View(r tea.Renderer) {
	pad := strings.Repeat(" ", padding)
	r.Write("\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit"))
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return msgTick(t)
	})
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func Main(ctx context.Context) error {
	m := &model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
