package fullscreen

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"context"
	"fmt"
	"time"

	"github.com/rprtr258/tea"
)

type model int

type msgTick time.Time

func (m *model) Init(f func(...tea.Cmd)) {
	f(tick(), tea.EnterAltScreen)
}

func (m *model) Update(message tea.Msg, f func(...tea.Cmd)) {
	switch msg := message.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			f(tea.Quit)
		}
	case msgTick:
		*m--
		if *m <= 0 {
			f(tea.Quit)
			return
		}
		f(tick())
	}
}

func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf("\n\n     Hi. This program will exit in %d seconds...", *m))
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return msgTick(t)
	})
}

func Main(ctx context.Context) error {
	m := model(5)
	_, err := tea.
		NewProgram(ctx, &m).
		WithAltScreen().
		Run()
	return err
}
