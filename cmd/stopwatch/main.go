package stopwatch

import (
	"context"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/stopwatch"
)

type model struct {
	stopwatch stopwatch.Model
	keymap    keymap
	help      help.Model
	quitting  bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(m.stopwatch.Init()...)
}

func (m *model) View(vb tea.Viewbox) {
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.
	if !m.quitting {
		vb = vb.WriteLineX("Elapsed: ")
	}
	vb = vb.WriteLineX(m.stopwatch.View())
	if !m.quitting {
		m.help.ShortHelpView(vb.PaddingTop(2), []key.Binding{
			m.keymap.start,
			m.keymap.stop,
			m.keymap.reset,
			m.keymap.quit,
		})
	}
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			f(tea.Quit)
		case key.Matches(msg, m.keymap.reset):
			f(m.stopwatch.CmdReset())
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			f(m.stopwatch.Toggle()...)
		}
	}

	f(m.stopwatch.Update(msg)...)
}

func Main(ctx context.Context) error {
	m := &model{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
		keymap: keymap{
			start: key.Binding{
				Keys: []string{"s"},
				Help: key.Help{"s", "start"},
			},
			stop: key.Binding{
				Keys: []string{"s"},
				Help: key.Help{"s", "stop"},
			},
			reset: key.Binding{
				Keys: []string{"r"},
				Help: key.Help{"r", "reset"},
			},
			quit: key.Binding{
				Keys: []string{"ctrl+c", "q"},
				Help: key.Help{"q", "quit"},
			},
		},
		help: help.New(),
	}

	m.keymap.start.SetEnabled(false)

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
