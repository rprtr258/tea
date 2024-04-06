package pipe

// An example illustrating how to pipe in data to a Tea application.
// More so, this serves as proof that Tea will automatically listen for
// keystrokes when input is not a TTY, such as when data is piped or redirected
// in.

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/textinput"
	"github.com/rprtr258/tea/styles"
)

type model struct {
	userInput textinput.Model
}

func newModel(initialValue string) *model {
	i := textinput.New()
	i.Prompt = ""
	i.Cursor.Style = styles.Style{}.Foreground(styles.FgColor("63"))
	i.Width = 48
	i.SetValue(initialValue)
	i.CursorEnd()
	i.Focus()

	return &model{
		userInput: i,
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(textinput.Blink)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	if key, ok := msg.(tea.MsgKey); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyEnter:
			f(tea.Quit)
			return
		}
	}

	m.userInput.Update(msg, f)
}

func (m *model) View(vb tea.Viewbox) {
	const _x = "You piped in: "
	vb = vb.PaddingTop(1)
	vb.WriteLine(_x)
	m.userInput.View(vb.PaddingLeft(len(_x)))
	vb.PaddingTop(2).WriteLine("Press ^C to exit")
}

func Main(ctx context.Context) error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return errors.New("try piping in some text")
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
