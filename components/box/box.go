package box

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
	BorderMaskAll = BorderMaskTop | BorderMaskBottom | BorderMaskLeft | BorderMaskRight
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

// Box draws `inside` inside a box(border)
func Box(
	vb tea.Viewbox,
	inside func(tea.Viewbox),
	border Border,
	borders BorderMask,
	fg [4]scuf.Modifier,
	bg [4]scuf.Modifier,
) {
	if inside != nil {
		inside(vb.Padding(tea.PaddingOptions{
			Top:    1,
			Left:   1,
			Bottom: 1,
			Right:  1,
		}))
	}

	// If no border is set or all borders are disabled, abort
	if border == noBorder || borders == 0 {
		return
	}

	topFG, rightFG, bottomFG, leftFG := fg[0], fg[1], fg[2], fg[3]
	topBG, rightBG, bottomBG, leftBG := bg[0], bg[1], bg[2], bg[3]

	// Figure out which corners we should actually be using based on which sides are set to show.
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
		vb.Pixel(0, 0).Styled(styleTop).Fill(border.TopLeft)
	}
	if borders.GetTop() && borders.GetRight() {
		vb.Pixel(0, vb.Width-1).Styled(styleTop).Fill(border.TopRight)
	}
	if borders.GetBottom() && borders.GetLeft() {
		vb.Pixel(vb.Height-1, 0).Styled(styleBottom).Fill(border.BottomLeft)
	}
	if borders.GetBottom() && borders.GetRight() {
		vb.Pixel(vb.Height-1, vb.Width-1).Styled(styleBottom).Fill(border.BottomRight)
	}
	// borders
	if borders.GetTop() {
		vb.
			Sub(tea.Rectangle{
				Left:   1,
				Height: 1,
				Width:  vb.Width - 2,
			}).
			Styled(styleTop).
			Fill(border.Top)
	}
	if borders.GetBottom() {
		vb.
			Sub(tea.Rectangle{
				Top:    vb.Height - 1,
				Left:   1,
				Height: 1,
				Width:  vb.Width - 2,
			}).
			Styled(styleBottom).
			Fill(border.Bottom)
	}
	if borders.GetLeft() {
		vb.
			Sub(tea.Rectangle{
				Top:    1,
				Height: vb.Height - 2,
				Width:  1,
			}).
			Styled(styles.
				Style{}.
				Foreground(leftFG).
				Background(leftBG)).
			Fill(border.Left)
	}
	if borders.GetRight() {
		vb.
			Sub(tea.Rectangle{
				Left:   vb.Width - 1,
				Top:    1,
				Height: vb.Height - 2,
				Width:  1,
			}).
			Styled(styles.
				Style{}.
				Foreground(rightFG).
				Background(rightBG)).
			Fill(border.Right)
	}
}
