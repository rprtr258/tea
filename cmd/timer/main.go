package timer

import (
	"context"
	"time"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/timer"
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
	m.timer.Init(f)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case timer.MsgTick:
		m.timer.Update(msg, f)
	case timer.MsgStartStop:
		m.timer.Update(msg, f)
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
	if !m.quitting {
		vb = vb.WriteLineX("Exiting in ")
	}
	m.timer.View(vb)
	vb = vb.PaddingLeft(5)
	vb = vb.WriteLineX(fun.IF(
		m.timer.Timedout(),
		"All done!",
		"",
	))
	if !m.quitting {
		m.helpView(vb.Padding(tea.PaddingOptions{Top: 1}))
	}
}

func Main(ctx context.Context) error {
	m := &model{
		timer: timer.NewWithInterval(timeout, time.Millisecond),
		keymap: keymap{
			start: key.Binding{
				Keys:     []string{"s"},
				Help:     key.Help{"s", "start"},
				Disabled: true,
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
				Keys: []string{"q", "ctrl+c"},
				Help: key.Help{"q", "quit"},
			},
		},
		help: help.New(),
	}

	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
