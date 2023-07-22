package pipe

// An example illustrating how to pipe in data to a Bubble Tea application.
// More so, this serves as proof that Bubble Tea will automatically listen for
// keystrokes when input is not a TTY, such as when data is piped or redirected
// in.

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
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

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.MsgKey); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			return tea.Quit
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

func Main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		log.Fatalln("Try piping in some text.")
	}

	reader := bufio.NewReader(os.Stdin)

	var sb strings.Builder
	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = sb.WriteRune(r)
		if err != nil {
			log.Fatalln("Error getting input:", err.Error())
		}
	}

	model := newModel(strings.TrimSpace(sb.String()))

	if _, err := tea.NewProgram(context.Background(), model).Run(); err != nil {
		log.Fatalln("Couldn't start program:", err.Error())
	}
}
