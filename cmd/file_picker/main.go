package file_picker //nolint:revive,stylecheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/filepicker"
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

func (m *model) Init(yield func(...tea.Cmd)) {
	m.filepicker.Init(yield)
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			yield(tea.Quit)
			return
		}
	case msgClearError:
		m.err = nil
	}

	yield(m.filepicker.Update(msg)...)

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
		yield(clearErrorAfter(2 * time.Second))
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		return
	}

	vb = vb.PaddingTop(1)

	switch {
	case m.err != nil:
		vb.PaddingLeft(1).WriteLine(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	case m.selectedFile == "":
		vb.PaddingLeft(1).WriteLine("Pick a file:")
	default:
		vb.PaddingLeft(1).WriteLine("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	m.filepicker.View(vb.PaddingTop(2))
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
