package fullscreen

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rprtr258/tea"
)

type model int

type msgTick time.Time

func Main() {
	m := model(5)
	p := tea.NewProgram(context.Background(), &m).WithAltScreen()
	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func (m *model) Update(message tea.Msg) tea.Cmd {
	switch msg := message.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return tea.Quit
		}

	case msgTick:
		*m--
		if *m <= 0 {
			return tea.Quit
		}
		return tick()
	}

	return nil
}

func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf("\n\n     Hi. This program will exit in %d seconds...", *m))
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return msgTick(t)
	})
}
