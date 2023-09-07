package altscreen_toggle

import (
	"context"

	termenv "github.com/rprtr258/col"
	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
)

var (
	color   = termenv.EnvColorProfile().Color
	keyword = termenv.S().Foreground(color("204")).Background(color("235")).Render
	help    = termenv.S().Foreground(color("241")).Render
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

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		vb.WriteLine(0, 0, "Bye!")
		return
	}

	const (
		altscreenMode = " altscreen mode "
		inlineMode    = " inline mode "
	)

	mode := fun.IF(
		m.altscreen,
		altscreenMode,
		inlineMode,
	)

	x := vb.WriteLine(2, 0, "  You're in ")
	vb.WriteLine(2, x, keyword(mode))

	vb.WriteLine(6, 0, help("  space: switch modes • q: exit"))
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
