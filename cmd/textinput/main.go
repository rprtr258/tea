package textinput

// A simple program demonstrating the text input component.

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/textinput"
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

func (m *model) Init(f func(...tea.Cmd)) {
	f(textinput.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			f(tea.Quit)
			return
		}
	}

	m.textInput.Update(msg, f)
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine("What’s your favorite Pokémon?")
	vb = vb.PaddingTop(2)
	m.textInput.View(vb)
	vb = vb.PaddingTop(2)
	vb.WriteLine("(esc to quit)")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
