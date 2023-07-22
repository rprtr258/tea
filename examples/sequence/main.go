package sequence

// A simple example illustrating how to run a series of commands in order.

import (
	"context"
	"log"

	"github.com/rprtr258/tea"
)

type model struct{}

func (m *model) Init() tea.Cmd {
	return tea.Sequence(
		tea.Batch(
			tea.Println("A"),
			tea.Println("B"),
			tea.Println("C"),
		),
		tea.Println("Z"),
		tea.Quit,
	)
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case tea.MsgKey:
		return tea.Quit
	}
	return nil
}

func (m *model) View(r tea.Renderer) {
}

func Main() {
	if _, err := tea.NewProgram(context.Background(), &model{}).Run(); err != nil {
		log.Fatalln("Uh oh:", err.Error())
	}
}
