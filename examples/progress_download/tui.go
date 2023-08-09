package progress_download

import (
	"strings"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/progress"
	"github.com/rprtr258/tea/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type (
	msgProgress    float64
	msgProgressErr struct{ err error }
)

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type model struct {
	pw       *progressWriter
	progress progress.Model
	err      error
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
	case msgProgressErr:
		m.err = msg.err
		f(tea.Quit)
	case msgProgress:
		if msg >= 1.0 {
			f(finalPause(), tea.Quit)
		}

		f(m.progress.SetPercent(float64(msg)))
	// MsgFrame is sent when the progress bar wants to animate itself
	case progress.MsgFrame:
		f(m.progress.Update(msg)...)
	}
}

func (m *model) View(r tea.Renderer) {
	if m.err != nil {
		r.Write("Error downloading: " + m.err.Error() + "\n")
		return
	}

	pad := strings.Repeat(" ", padding)
	r.Write("\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit"))
}
