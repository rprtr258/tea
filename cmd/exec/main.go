package exec

import (
	"context"
	"os"
	"os/exec"

	"github.com/rprtr258/tea"
)

type msgEditorFinished struct{ err error }

func openEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return msgEditorFinished{err}
	})
}

type model struct {
	altscreenActive bool
	err             error
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "a":
			m.altscreenActive = !m.altscreenActive
			cmd := tea.EnterAltScreen
			if !m.altscreenActive {
				cmd = tea.ExitAltScreen
			}
			f(cmd)
		case "e":
			f(openEditor())
		case "ctrl+c", "q":
			f(tea.Quit)
		}
	case msgEditorFinished:
		if msg.err != nil {
			m.err = msg.err
			f(tea.Quit)
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.err != nil {
		vb.WriteLine("Error: " + m.err.Error())
		return
	}

	vb.WriteLine("Press 'e' to open your EDITOR.")
	vb.PaddingTop(1).WriteLine("Press 'a' to toggle the altscreen")
	vb.PaddingTop(2).WriteLine("Press 'q' to quit.")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
