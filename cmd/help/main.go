package help

import (
	"context"
	"log"
	"os"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up, Down, Left, Right key.Binding
	Help                  key.Binding
	Quit                  key.Binding
}

var keys = keyMap{
	Up: key.Binding{
		Keys: []string{"up", "k"},
		Help: key.Help{"↑/k", "move up"},
	},
	Down: key.Binding{
		Keys: []string{"down", "j"},
		Help: key.Help{"↓/j", "move down"},
	},
	Left: key.Binding{
		Keys: []string{"left", "h"},
		Help: key.Help{"←/h", "move left"},
	},
	Right: key.Binding{
		Keys: []string{"right", "l"},
		Help: key.Help{"→/l", "move right"},
	},
	Help: key.Binding{
		Keys: []string{"?"},
		Help: key.Help{"?", "toggle help"},
	},
	Quit: key.Binding{
		Keys: []string{"q", "esc", "ctrl+c"},
		Help: key.Help{"q", "quit"},
	},
}

type model struct {
	keys       keyMap
	keysMap    help.KeyMap
	help       help.Model
	inputStyle styles.Style
	lastKey    string
	quitting   bool
}

func newModel() *model {
	return &model{
		keys: keys,
		keysMap: help.KeyMap{
			// ShortHelp returns keybindings to be shown in the mini help view.
			ShortHelp: []key.Binding{keys.Help, keys.Quit},
			// FullHelp returns keybindings for the expanded help view.
			FullHelp: [][]key.Binding{
				{keys.Up, keys.Down, keys.Left, keys.Right}, // first column
				{keys.Help, keys.Quit},                      // second column
			},
		},
		help:       help.New(),
		inputStyle: styles.Style{}.Foreground(styles.FgColor("#FF75B7")),
	}
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Left):
			m.lastKey = "←"
		case key.Matches(msg, m.keys.Right):
			m.lastKey = "→"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			f(tea.Quit)
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		vb.WriteLine("Bye!")
		return
	}

	status := fun.IF(
		m.lastKey == "",
		"Waiting for input...",
		"You chose: "+m.inputStyle.Render(m.lastKey),
	)

	vb = vb.PaddingTop(1)
	vb.WriteLine(status)
	vb = vb.PaddingTop(1)

	m.help.View(vb, m.keysMap)
}

func Main(ctx context.Context) error {
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "help")
		if err != nil {
			log.Fatalln("Couldn't open a file for logging:", err.Error())
		}
		defer f.Close() // nolint:errcheck
	}

	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
