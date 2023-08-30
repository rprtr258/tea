package timer

import (
	"context"
	"time"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/help"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/bubbles/timer"
)

const timeout = time.Second * 5

type model struct {
	timer    timer.Model
	keymap   keymap
	help     help.Model
	quitting bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(m.timer.Init()...)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case timer.MsgTick:
		f(m.timer.Update(msg)...)
	case timer.MsgStartStop:
		f(m.timer.Update(msg)...)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
	case timer.MsgTimeout:
		m.quitting = true
		f(tea.Quit)
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			f(tea.Quit)
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = timeout
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			f(m.timer.CmdToggle())
		}
	}
}

func (m *model) helpView(vb tea.Viewbox) {
	m.help.ShortHelpView(vb.Padding(tea.PaddingOptions{Top: 1}), []key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m *model) View(vb tea.Viewbox) {
	// For a more detailed timer view you could read m.timer.Timeout to get the
	// remaining time as a time.Duration and skip calling m.timer.View() entirely.
	x := 0
	if !m.quitting {
		x = vb.WriteLine(0, 0, "Exiting in ")
	}
	vb.WriteLine(0, x, fun.IF(
		m.timer.Timedout(),
		"All done!",
		m.timer.View(),
	))
	if !m.quitting {
		m.helpView(vb.Padding(tea.PaddingOptions{Top: 1}))
	}
}

func Main(ctx context.Context) error {
	m := &model{
		timer: timer.NewWithInterval(timeout, time.Millisecond),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
				key.WithDisabled(),
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
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
