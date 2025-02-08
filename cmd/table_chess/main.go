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
	labelsY, vbLabelsX := vb.SplitX2(tea.Fixed(3), tea.Fixed(8*3+8+1))
	vbBoard, labelsX := vbLabelsX.SplitY2(tea.Fixed(8*1+8+1), tea.Fixed(2))

	labelStyle := styles.Style{}.Foreground(styles.ANSIColor(241))
	labelsY = labelsY.Styled(labelStyle)
	labelsX = labelsX.Styled(labelStyle)
	for i := 0; i < 8; i++ {
		labelsY.Set(i*2+1, 1, []rune("12345678")[i])
		labelsX.Set(0, i*4+2, []rune("ABCDEFGH")[i])
	}

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

	for i, vb := range tablebox.Box(
		vbBoard,
		tablebox.NormalBorder,
		styles.Style{},
		tablebox.Grid(
			[]tea.Layout{tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto()},
			[]tea.Layout{tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto(), tea.Auto()},
		)...,
	) {
		y, x := i/8, i%8
		vb.Set(0, vb.Width/2, board[y][x])
	}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
