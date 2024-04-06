package result

// A simple example that shows how to retrieve a value from a Tea
// program after the Tea has exited.

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
)

var choices = []string{"Taro", "Coffee", "Lychee"}

type model struct {
	cursor int
	choice string
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			f(tea.Quit)
		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			f(tea.Quit)
		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine("What kind of Tea would you like to order?")
	vb = vb.Padding(tea.PaddingOptions{Top: 1})
	for i := 0; i < len(choices); i++ {
		vb.WriteLine("( ) " + choices[i])
		vb = vb.PaddingTop(1)
	}
	vb.Set(m.cursor, 1, '•')
	vb.PaddingTop(1).WriteLine("(press q to quit)")
}

func Main(ctx context.Context) error {
	// Run returns the model as a tea.Model.
	m, err := tea.NewProgram(ctx, &model{}).Run()
	if err != nil {
		return err
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m.choice != "" {
		fmt.Printf("\n---\nYou chose %s!\n", m.choice)
	}

	return nil
}
