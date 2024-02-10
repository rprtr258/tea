package textinputs

// A simple example demonstrating the use of multiple text input components.

import (
	"context"
	"fmt"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/cursor"
	"github.com/rprtr258/tea/components/textinput"
	"github.com/rprtr258/tea/styles"
)

var (
	focusedStyle        = styles.Style{}.Foreground(styles.FgColor("205"))
	blurredStyle        = styles.Style{}.Foreground(styles.FgColor("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = styles.Style{}
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = styles.Style{}.Foreground(styles.FgColor("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel() *model {
	m := &model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textinput.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "esc":
			f(tea.Quit)
			return
		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.ModeHide {
				m.cursorMode = cursor.ModeBlink
			}
			for i := range m.inputs {
				f(m.inputs[i].Cursor.SetMode(m.cursorMode)...)
			}
			return
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				f(tea.Quit)
				return
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					f(m.inputs[i].Focus()...)
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return
		}
	}

	// Handle character input and blinking
	m.updateInputs(msg, f)
}

func (m *model) updateInputs(msg tea.Msg, f func(...tea.Cmd)) {
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i].Update(msg, f)
	}
}

func (m *model) View(vb tea.Viewbox) {
	for i := range m.inputs {
		m.inputs[i].View(vb)
		vb = vb.PaddingTop(1)
	}

	vb = vb.PaddingTop(1)
	vb.WriteLine(fun.IF(m.focusIndex == len(m.inputs), focusedButton, blurredButton))
	vb = vb.PaddingTop(1)
	vb = vb.Styled(helpStyle).WriteLineX("cursor mode is ")
	vb = vb.Styled(cursorModeHelpStyle).WriteLineX(m.cursorMode.String())
	vb.Styled(helpStyle).WriteLine(" (ctrl+r to change style)")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
