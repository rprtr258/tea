package mouse

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
)

type model struct {
	init       bool
	mouseEvent tea.MouseEvent
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			f(tea.Quit)
		}
	case tea.MsgMouse:
		m.init = true
		m.mouseEvent = tea.MouseEvent(msg)
	}
}

func (m *model) View(r tea.Renderer) {
	s := "Do mouse stuff. When you're done press q to quit.\n\n"

	if m.init {
		e := m.mouseEvent
		s += fmt.Sprintf("(X: %d, Y: %d) %s", e.X, e.Y, e)
	}

	r.Write(s)
}

func Main(ctx context.Context) error {
	_, err := tea.
		NewProgram(ctx, &model{}).
		WithAltScreen().
		WithMouseAllMotion().
		Run()
	return err
}
