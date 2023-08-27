package basics

import (
	"context"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/tea"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func initialModel() *model {
	return &model{
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			f(tea.Quit)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
}

func (m *model) View(r tea.Renderer) {
	r.Write("What should we buy at the market?\n\n")

	for i, choice := range m.choices {
		r.Write(fun.IF(m.cursor == i, ">", " "))

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		r.Write(" [")
		r.Write(checked)
		r.Write("] ")

		r.Write(choice)
		r.Write("\n")
	}

	r.Write("\nPress q to quit.\n")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
