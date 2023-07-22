package textinput

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"context"
	"fmt"
	"log"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/textinput"
)

type model struct {
	textInput textinput.Model
}

func initialModel() *model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &model{
		textInput: ti,
	}
}

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{textinput.Blink}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return []tea.Cmd{tea.Quit}
		}
	}

	return m.textInput.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n")
}

func Main() {
	p := tea.NewProgram(context.Background(), initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
