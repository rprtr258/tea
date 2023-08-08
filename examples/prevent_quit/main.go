package prevent_quit

// A program demonstrating how to use the WithFilter option to intercept events.

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/help"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/bubbles/textarea"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))
	saveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
)

type model struct {
	textarea   textarea.Model
	help       help.Model
	keymap     keymap
	saveText   string
	hasChanges bool
	quitting   bool
}

type keymap struct {
	save key.Binding
	quit key.Binding
}

func initialModel() *model {
	ti := textarea.New()
	ti.Placeholder = "Only the best words"
	ti.Focus()

	return &model{
		textarea: ti,
		help:     help.New(),
		keymap: keymap{
			save: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "save"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textarea.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	if m.quitting {
		m.updatePromptView(msg, f)
		return
	}

	m.updateTextView(msg, f)
}

func (m *model) updateTextView(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		m.saveText = ""
		switch {
		case key.Matches(msg, m.keymap.save):
			m.saveText = "Changes saved!"
			m.hasChanges = false
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			f(tea.Quit)
			return
		case msg.Type == tea.KeyRunes:
			m.saveText = ""
			m.hasChanges = true
			fallthrough
		default:
			if !m.textarea.Focused() {
				f(m.textarea.Focus()...)
			}
		}
	}
	f(m.textarea.Update(msg)...)
}

func (m *model) updatePromptView(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		// For simplicity's sake, we'll treat any key besides "y" as "no"
		if key.Matches(msg, m.keymap.quit) || msg.String() == "y" {
			m.hasChanges = false
			f(tea.Quit)
			return
		}
		m.quitting = false
	}
}

func (m *model) View(r tea.Renderer) {
	if m.quitting {
		if m.hasChanges {
			text := lipgloss.JoinHorizontal(lipgloss.Top, "You have unsaved changes. Quit without saving?", choiceStyle.Render("[yn]"))
			r.Write(quitViewStyle.Render(text))
			return
		}

		r.Write("Very important, thank you\n")
		return
	}

	helpView := m.help.ShortHelpView([]key.Binding{
		m.keymap.save,
		m.keymap.quit,
	})

	r.Write(fmt.Sprintf(
		"\nType some important things.\n\n%s\n\n %s\n %s",
		m.textarea.View(),
		saveTextStyle.Render(m.saveText),
		helpView,
	) + "\n\n")
}

func Main(ctx context.Context) error {
	_, err := tea.
		NewProgram(ctx, initialModel()).
		WithFilter(func(m *model, msg tea.Msg) tea.Msg {
			if _, ok := msg.(tea.MsgQuit); !ok {
				return msg
			}

			if !m.hasChanges {
				return msg
			}

			return nil
		}).
		Run()
	return err
}
