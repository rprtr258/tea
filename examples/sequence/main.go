package sequence

// A simple example illustrating how to run a series of commands in order.

import (
	"context"

	"github.com/rprtr258/tea"
)

type model struct{}

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{
		tea.Println("A"),
		tea.Println("B"),
		tea.Println("C"),
		tea.Println("Z"),
		tea.Quit,
	}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg.(type) {
	case tea.MsgKey:
		return []tea.Cmd{tea.Quit}
	}
	return nil
}

func (m *model) View(r tea.Renderer) {}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
