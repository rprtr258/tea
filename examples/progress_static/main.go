package progress_static

// A simple example that shows how to render a progress bar in a "pure"
// fashion. In this example we bump the progress by 25% every second,
// maintaining the progress state on our top level model using the progress bar
// model's ViewAs method only for rendering.
//
// The signature for ViewAs is:
//
//     func (m *model) ViewAs(percent float64) string
//
// So it takes a float between 0 and 1, and renders the progress bar
// accordingly. When using the progress bar in this "pure" fashion and there's
// no need to call an Update method.
//
// The progress bar is also able to animate itself, however. For details see
// the progress-animated example.

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

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type msgTick time.Time

var cmdTick = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return msgTick(t)
})

type model struct {
	percent  float64
	progress progress.Model
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(cmdTick)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
	case msgTick:
		m.percent += 0.25
		if m.percent > 1.0 {
			m.percent = 1.0
			f(tea.Quit)
			return
		}
		f(cmdTick)
	}
}

func (m *model) View(r tea.Renderer) {
	pad := strings.Repeat(" ", padding)
	r.Write("\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit"))
}

func Main(ctx context.Context) error {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	_, err := tea.NewProgram(ctx, &model{progress: prog}).Run()
	return err
}
