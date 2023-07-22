package exec

import (
	"context"
	"log"
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

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "a":
			m.altscreenActive = !m.altscreenActive
			cmd := tea.EnterAltScreen
			if !m.altscreenActive {
				cmd = tea.ExitAltScreen
			}
			return cmd
		case "e":
			return openEditor()
		case "ctrl+c", "q":
			return tea.Quit
		}
	case msgEditorFinished:
		if msg.err != nil {
			m.err = msg.err
			return tea.Quit
		}
	}
	return nil
}

func (m *model) View(r tea.Renderer) {
	if m.err != nil {
		r.Write("Error: " + m.err.Error() + "\n")
		return
	}

	r.Write("Press 'e' to open your EDITOR.\nPress 'a' to toggle the altscreen\nPress 'q' to quit.\n")
}

func Main() {
	if _, err := tea.NewProgram(context.Background(), &model{}).Run(); err != nil {
		log.Fatalln("Error running program:", err.Error())
	}
}
