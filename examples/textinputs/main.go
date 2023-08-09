package textinputs

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"fmt"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/cursor"
	"github.com/rprtr258/tea/bubbles/textinput"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

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
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "esc":
			f(tea.Quit)
			return
		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
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
	f(m.updateInputs(msg)...)
}

func (m *model) updateInputs(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		cmds = append(cmds, m.inputs[i].Update(msg)...)
	}

	return cmds
}

func (m *model) View(r tea.Renderer) {
	var sb strings.Builder

	for i := range m.inputs {
		sb.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			sb.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&sb, "\n\n%s\n\n", *button)

	sb.WriteString(helpStyle.Render("cursor mode is "))
	sb.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	sb.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	r.Write(sb.String())
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
