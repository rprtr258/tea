package table_pokemon

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/tablebox"
	"github.com/rprtr258/tea/styles"
)

type model struct{}

func (*model) Init(func(...tea.Cmd)) {}

func (*model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c":
			f(tea.Quit)
			return
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	tablebox.Box(
		vb.MaxWidth(80).MaxHeight(32),
		[]tea.Layout{tea.Fixed(1), tea.Fixed(28)}, // TODO: fixed(1), auto
		[]tea.Layout{tea.Fixed(6), tea.Fixed(14), tea.Fixed(12), tea.Fixed(10), tea.Fixed(14), tea.Fixed(17)},
		func(vb tea.Viewbox, y, x int) {
			// vb.Fill('0' + rune(max(y, x)))
		},
		tablebox.NormalBorder,
		styles.Style{}.Foreground(styles.ANSIColor(238)),
	)
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
