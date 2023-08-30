package basics

import (
	"context"

	"github.com/rprtr258/tea"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	vb       tea.Viewbox
}

func initialModel() *model {
	return &model{
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		cursor:   0,

		vb: tea.NewViewbox(0, 0),
	}
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgWindowSize:
		// TODO: m.vb.SetSize(msg.Width, msg.Height)
		m.vb = tea.NewViewbox(msg.Height, msg.Width)
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			yield(tea.Quit)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
}

func (m *model) View(r tea.Renderer) {
	m.vb.Clear()

	m.vb.WriteLine(0, 0, "What should we buy at the market?")

	vbChoices := m.vb.Padding(tea.PaddingOptions{Top: 2})
	vbChoices.Set(m.cursor, 0, '>')
	for i, choice := range m.choices {
		// 0123456789...
		// > [x] Buy carrots
		vbChoices.Set(i, 2, '[')
		if _, ok := m.selected[i]; ok {
			vbChoices.Set(i, 3, 'x')
		}
		vbChoices.Set(i, 4, ']')
		vbChoices.WriteLine(i, 6, choice)
	}

	vbChoices.WriteLine(len(m.choices)+1, 0, "Press q to quit.")

	r.Write(m.vb.Render())
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
