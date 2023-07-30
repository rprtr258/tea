package pipe

// An example illustrating how to pipe in data to a Bubble Tea application.
// More so, this serves as proof that Bubble Tea will automatically listen for
// keystrokes when input is not a TTY, such as when data is piped or redirected
// in.

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/textinput"
	"github.com/rprtr258/tea/lipgloss"
)

type model struct {
	userInput textinput.Model
}

func newModel(initialValue string) *model {
	i := textinput.New()
	i.Prompt = ""
	i.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	i.Width = 48
	i.SetValue(initialValue)
	i.CursorEnd()
	i.Focus()

	return &model{
		userInput: i,
	}
}

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{textinput.Blink}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	if key, ok := msg.(tea.MsgKey); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			return []tea.Cmd{tea.Quit}
		}
	}

	return m.userInput.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf(
		"\nYou piped in: %s\n\nPress ^C to exit",
		m.userInput.View(),
	))
}

func Main(ctx context.Context) error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return errors.New("Try piping in some text.")
	}

	reader := bufio.NewReader(os.Stdin)

	var sb strings.Builder
	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}

		sb.WriteRune(r)
	}

	model := newModel(strings.TrimSpace(sb.String()))

	_, err = tea.NewProgram(ctx, model).Run()
	return err
}
