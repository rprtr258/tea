package altscreen_toggle

import (
	"context"
	"fmt"

	"github.com/muesli/termenv"
	"github.com/rprtr258/tea"
)

var (
	color   = termenv.EnvColorProfile().Color
	keyword = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	help    = termenv.Style{}.Foreground(color("241")).Styled
)

type model struct {
	altscreen bool
	quitting  bool
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			f(tea.Quit)
		case " ":
			if m.altscreen {
				f(tea.ExitAltScreen)
			} else {
				f(tea.EnterAltScreen)
			}
			m.altscreen = !m.altscreen
		}
	}
}

func (m *model) View(r tea.Renderer) {
	if m.quitting {
		r.Write("Bye!\n")
		return
	}

	const (
		altscreenMode = " altscreen mode "
		inlineMode    = " inline mode "
	)

	var mode string
	if m.altscreen {
		mode = altscreenMode
	} else {
		mode = inlineMode
	}

	r.Write(fmt.Sprintf("\n\n  You're in %s\n\n\n%s", keyword(mode), help("  space: switch modes • q: exit\n")))
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
