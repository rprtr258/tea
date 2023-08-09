package file_picker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/filepicker"
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type msgClearError struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return msgClearError{}
	})
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(m.filepicker.Init())
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			f(tea.Quit)
			return
		}
	case msgClearError:
		m.err = nil
	}

	f(m.filepicker.Update(msg)...)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		f(clearErrorAfter(2 * time.Second))
	}
}

func (m *model) View(r tea.Renderer) {
	if m.quitting {
		r.Write("")
		return
	}

	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	r.Write(s.String())
}

func Main(ctx context.Context) error {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := &model{
		filepicker: fp,
	}

	mm, err := tea.NewProgram(ctx, m).WithOutput(os.Stderr).Run()
	if err != nil {
		return err
	}

	fmt.Printf(
		"\n  You selected: %s\n\n",
		m.filepicker.Styles.Selected.Render(mm.selectedFile),
	)
	return nil
}
