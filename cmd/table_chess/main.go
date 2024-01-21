package table_chess

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
		vb.MaxWidth(33).MaxHeight(17),
		[]tea.Layout{tea.Fixed(1), tea.Fixed(1), tea.Fixed(1), tea.Fixed(1), tea.Fixed(1), tea.Fixed(1), tea.Fixed(1), tea.Fixed(1)},
		[]tea.Layout{tea.Fixed(3), tea.Fixed(3), tea.Fixed(3), tea.Fixed(3), tea.Fixed(3), tea.Fixed(3), tea.Fixed(3), tea.Fixed(3)}, // TODO: all flex(1)
		func(vb tea.Viewbox, y, x int) {
			board := [8][8]rune{
				{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'},
				{'♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'},
				{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
				{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
				{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
				{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
				{'♙', '♙', '♙', '♙', '♙', '♙', '♙', '♙'},
				{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'},
			}

			vb.Set(0, vb.Width/2, board[y][x])
		},
		tablebox.NormalBorder,
		styles.Style{}.Foreground(styles.ANSIColor(238)),
	)
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
