package box

import (
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

// GetBorder returns the style's border style (type Border) and value for the
// top, right, bottom, and left in that order. If no value is set for the
// border style, Border{} is returned. For all other unset values false is
// returned.
func GetBorder(mask BorderMask) (top, right, bottom, left bool) { //nolint:nonamedreturns
	return mask&BorderMaskTop != 0,
		mask&BorderMaskRight != 0,
		mask&BorderMaskBottom != 0,
		mask&BorderMaskLeft != 0
}

func Colors(cols ...styles.TerminalColor) [4]styles.TerminalColor {
	switch len(cols) {
	case 1:
		return [4]styles.TerminalColor{
			cols[0],
			cols[0],
			cols[0],
			cols[0],
		}
	case 2:
		return [4]styles.TerminalColor{
			cols[0],
			cols[1],
			cols[0],
			cols[1],
		}
	case 3:
		return [4]styles.TerminalColor{
			cols[0],
			cols[1],
			cols[2],
			cols[1],
		}
	case 4:
		return [4]styles.TerminalColor{
			cols[0],
			cols[1],
			cols[2],
			cols[3],
		}
	}
	return [4]styles.TerminalColor{
		styles.NoColor,
		styles.NoColor,
		styles.NoColor,
		styles.NoColor,
	}
}

// Box draws model inside a box
func Box(
	vb tea.Viewbox,
	inside func(tea.Viewbox),
	border Border,
	borders BorderMask,
	fg [4]styles.TerminalColor,
	bg [4]styles.TerminalColor,
) {
	if inside != nil {
		inside(vb.Padding(tea.PaddingOptions{
			Top:    1,
			Left:   1,
			Bottom: 1,
			Right:  1,
		}))
	}

	// If no border is set or all borders are been disabled, abort.
	if border == noBorder || borders == 0 {
		return
	}

	topSet, rightSet, bottomSet, leftSet := GetBorder(borders)
	topFG, rightFG, bottomFG, leftFG := fg[0], fg[1], fg[2], fg[3]
	topBG, rightBG, bottomBG, leftBG := bg[0], bg[1], bg[2], bg[3]

	// Figure out which corners we should actually be using based on which
	// sides are set to show.
	if topSet {
		if !leftSet {
			border.TopLeft = 0
		}
		if !rightSet {
			border.TopRight = 0
		}
	}
	if bottomSet {
		if !leftSet {
			border.BottomLeft = 0
		}
		if !rightSet {
			border.BottomRight = 0
		}
	}

	styleTop := styles.Style{}.Foreground(topFG).Background(topBG)
	styleBottom := styles.Style{}.Foreground(bottomFG).Background(bottomBG)
	// angles
	if topSet && leftSet {
		vb.Styled(styleTop).Set(0, 0, border.TopLeft)
	}
	if topSet && rightSet {
		vb.Styled(styleTop).Set(0, vb.Width-1, border.TopRight)
	}
	if bottomSet && leftSet {
		vb.Styled(styleBottom).Set(vb.Height-1, 0, border.BottomLeft)
	}
	if bottomSet && rightSet {
		vb.Styled(styleBottom).Set(vb.Height-1, vb.Width-1, border.BottomRight)
	}
	// borders
	if topSet {
		for i := 1; i < vb.Width-1; i++ {
			vb.Styled(styleTop).Set(0, i, border.Top)
		}
	}
	if bottomSet {
		for i := 1; i < vb.Width-1; i++ {
			vb.Styled(styleBottom).Set(vb.Height-1, i, border.Bottom)
		}
	}
	if leftSet {
		styleLeft := styles.Style{}.Foreground(leftFG).Background(leftBG)
		for i := 1; i < vb.Height-1; i++ {
			vb.Styled(styleLeft).Set(i, 0, border.Left)
		}
	}
	if rightSet {
		styleRight := styles.Style{}.Foreground(rightFG).Background(rightBG)
		for i := 1; i < vb.Height-1; i++ {
			vb.Styled(styleRight).Set(i, vb.Width-1, border.Right)
		}
	}
}
