package stopwatch

import (
	"context"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/help"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/bubbles/stopwatch"
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

func (m *model) View(r tea.Renderer) {
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.
	s := m.stopwatch.View() + "\n"
	if !m.quitting {
		s = "Elapsed: " + s
		s += m.helpView()
	}
	r.Write(s)
}

func (m *model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
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
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	m.keymap.start.SetEnabled(false)

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
