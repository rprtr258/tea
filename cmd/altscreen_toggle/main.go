package altscreen_toggle //nolint:revive,stylecheck

import (
	"context"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"

	"github.com/rprtr258/tea"
)

var (
	keyword = func(s string) string {
		return scuf.String(s, scuf.FgANSI(204), scuf.BgANSI(235))
	}
	help = func(s string) string {
		return scuf.String(s, scuf.FgANSI(241))
	}
)

type model struct {
	altscreen bool
	quitting  bool
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	if msg, ok := msg.(tea.MsgKey); ok { //nolint:gocritic
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			f(tea.Quit)
		case " ":
			f(fun.IF(m.altscreen, tea.ExitAltScreen, tea.EnterAltScreen))
			m.altscreen = !m.altscreen
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		vb.WriteLine("Bye!")
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

	vb = vb.PaddingTop(2)
	vb.WriteLine("  You're in " + keyword(mode))

	vb.PaddingTop(4).WriteLine(help("  space: switch modes • q: exit"))
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
