package tablebox

import (
	"github.com/rprtr258/scuf"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

type BorderMask int8

const (
	BorderMaskTop = 1 << iota
	BorderMaskBottom
	BorderMaskLeft
	BorderMaskRight
	BorderMaskHorizontals
	BorderMaskVerticals
	BorderMaskAll = BorderMaskTop | BorderMaskBottom |
		BorderMaskLeft | BorderMaskRight |
		BorderMaskHorizontals | BorderMaskVerticals
)

func (mask BorderMask) GetTop() bool {
	return mask&BorderMaskTop != 0
}
func (mask BorderMask) GetBottom() bool {
	return mask&BorderMaskBottom != 0
}
func (mask BorderMask) GetLeft() bool {
	return mask&BorderMaskLeft != 0
}
func (mask BorderMask) GetRight() bool {
	return mask&BorderMaskRight != 0
}

func Colors(cols ...scuf.Modifier) [4]scuf.Modifier {
	switch len(cols) {
	case 1:
		return [4]scuf.Modifier{
			cols[0],
			cols[0],
			cols[0],
			cols[0],
		}
	case 2:
		return [4]scuf.Modifier{
			cols[0],
			cols[1],
			cols[0],
			cols[1],
		}
	case 3:
		return [4]scuf.Modifier{
			cols[0],
			cols[1],
			cols[2],
			cols[1],
		}
	case 4:
		return [4]scuf.Modifier{
			cols[0],
			cols[1],
			cols[2],
			cols[3],
		}
	}
	return [4]scuf.Modifier{
		nil,
		nil,
		nil,
		nil,
	}
}

// Box draws model inside a box
func Box(
	vb tea.Viewbox,
	h, w int,
	inside func(vb tea.Viewbox, y, x int),
	border FullBorder,
	borders BorderMask,
	fg [4]scuf.Modifier,
	bg [4]scuf.Modifier,
) {
	u := make([]tea.Layout, h) // TODO: remove this cringe
	for i := range u {
		u[i] = 1
	}

	hs := tea.EvalLayout(vb.Height-h-1, u...)
	ws := tea.EvalLayout(vb.Width-w-1, u...)

	for iy := 0; iy < h; iy++ {
		for ix := 0; ix < w; ix++ {
			inside(vb.Sub(tea.Rectangle{
				Top:    iy + 1,
				Left:   ix + 1,
				Height: hs[iy],
				Width:  ws[ix],
			}), iy, ix)
		}
	}

	// If no border is set or all borders are been disabled, abort.
	if border == noBorder || borders == 0 {
		return
	}

	topFG, rightFG, bottomFG, leftFG := fg[0], fg[1], fg[2], fg[3]
	topBG, rightBG, bottomBG, leftBG := bg[0], bg[1], bg[2], bg[3]

	// Figure out which corners we should actually be using based on which sides are set to show.
	// TODO: support other options
	if borders.GetTop() {
		if !borders.GetLeft() {
			border.TopLeft = 0
		}
		if !borders.GetRight() {
			border.TopRight = 0
		}
	}
	if borders.GetBottom() {
		if !borders.GetLeft() {
			border.BottomLeft = 0
		}
		if !borders.GetRight() {
			border.BottomRight = 0
		}
	}

	styleTop := styles.Style{}.Foreground(topFG).Background(topBG)
	styleBottom := styles.Style{}.Foreground(bottomFG).Background(bottomBG)
	// angles
	if borders.GetTop() && borders.GetLeft() {
		vb.Styled(styleTop).Set(0, 0, border.TopLeft)
	}
	if borders.GetTop() && borders.GetRight() {
		vb.Styled(styleTop).Set(0, vb.Width-1, border.TopRight)
	}
	if borders.GetBottom() && borders.GetLeft() {
		vb.Styled(styleBottom).Set(vb.Height-1, 0, border.BottomLeft)
	}
	if borders.GetBottom() && borders.GetRight() {
		vb.Styled(styleBottom).Set(vb.Height-1, vb.Width-1, border.BottomRight)
	}
	// borders
	if borders.GetTop() {
		for i := 1; i < vb.Width-1; i++ {
			vb.Styled(styleTop).Set(0, i, border.Top)
		}
	}
	if borders.GetBottom() {
		for i := 1; i < vb.Width-1; i++ {
			vb.Styled(styleBottom).Set(vb.Height-1, i, border.Bottom)
		}
	}
	if borders.GetLeft() {
		styleLeft := styles.Style{}.Foreground(leftFG).Background(leftBG)
		for i := 1; i < vb.Height-1; i++ {
			vb.Styled(styleLeft).Set(i, 0, border.Left)
		}
	}
	if borders.GetRight() {
		styleRight := styles.Style{}.Foreground(rightFG).Background(rightBG)
		for i := 1; i < vb.Height-1; i++ {
			vb.Styled(styleRight).Set(i, vb.Width-1, border.Right)
		}
	}
}
