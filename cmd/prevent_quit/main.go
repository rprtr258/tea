package prevent_quit //nolint:revive,stylecheck

// A program demonstrating how to use the WithFilter option to intercept events.

import (
	"context"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/textarea"
	"github.com/rprtr258/tea/styles"
)

var (
	choiceStyle   = styles.Style{}. /*.PaddingLeft(1)*/ Foreground(styles.FgColor("241"))
	saveTextStyle = styles.Style{}.Foreground(styles.FgColor("170"))
	quitViewStyle = styles.Style{} /*.Padding(1) Border(styles.RoundedBorder).BorderForeground(styles.FgColor("170"))*/
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
			save: key.Binding{
				Keys: []string{"ctrl+s"},
				Help: key.Help{"ctrl+s", "save"},
			},
			quit: key.Binding{
				Keys: []string{"esc", "ctrl+c"},
				Help: key.Help{"esc", "quit"},
			},
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
	switch msg := msg.(type) { //nolint:gocritic
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
	m.textarea.Update(msg, f)
}

func (m *model) updatePromptView(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
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

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		if m.hasChanges {
			vb.Styled(quitViewStyle).WriteText(0, 0, styles.JoinHorizontal(
				styles.Top,
				[]string{"You have unsaved changes. Quit without saving?"},
				strings.Split(choiceStyle.Render("[yn]"), "\n"),
			))
			return
		}

		vb.WriteLine("Very important, thank you")
		return
	}

	vb.PaddingTop(1).WriteLine("Type some important things.")
	m.textarea.View(vb.Padding(tea.PaddingOptions{Top: 3}))
	h := m.textarea.Height()
	vb.PaddingTop(3 + h).PaddingLeft(1).Styled(saveTextStyle).WriteLine(m.saveText)
	m.help.ShortHelpView(vb.Styled(saveTextStyle).Padding(tea.PaddingOptions{Top: 4 + h, Left: 1}), []key.Binding{
		m.keymap.save,
		m.keymap.quit,
	})
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
