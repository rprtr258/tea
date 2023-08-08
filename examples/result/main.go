package result

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"context"
	"fmt"
	"strings"

	"github.com/rprtr258/tea"
)

var choices = []string{"Taro", "Coffee", "Lychee"}

type model struct {
	cursor int
	choice string
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
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

func (m *model) View(r tea.Renderer) {
	s := strings.Builder{}
	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	r.Write(s.String())
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
