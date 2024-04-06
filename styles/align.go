package styles

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/scuf"
)

// Perform text alignment. If the string is multi-lined, we also make all lines
// the same width by padding them with spaces. If a termenv style is passed,
// use that to style the spaces added.
func alignTextHorizontal(str string, pos Alignment, width int, style []scuf.Modifier) string {
	lines := strings.Split(str, "\n")
	widestLine := getWidestWidth(lines)

	var sb strings.Builder
	for i, line := range lines {
		lineWidth := ansi.PrintableRuneWidth(line)

		// difference from the widest line plus
		// difference from the total width, if set
		if shortAmount := widestLine - lineWidth + max(0, width-widestLine); shortAmount > 0 {
			switch pos {
			case Right:
				line = scuf.String(strings.Repeat(" ", shortAmount), style...) + line
			case Center:
				left := shortAmount / 2
				right := left + shortAmount%2 // note that we put the remainder on the right

				leftSpaces := scuf.String(strings.Repeat(" ", left), style...)
				rightSpaces := scuf.String(strings.Repeat(" ", right), style...)

				line = leftSpaces + line + rightSpaces
			default: // Left
				line += scuf.String(strings.Repeat(" ", shortAmount), style...)
			}
		}

		sb.WriteString(line)
		if i < len(lines)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func alignTextVertical(str string, pos Alignment, height int) string {
	strHeight := strings.Count(str, "\n") + 1
	if height < strHeight {
		return str
	}

	switch pos {
	case Top:
		return str + strings.Repeat("\n", height-strHeight)
	case Center:
		topPadding, bottomPadding := (height-strHeight)/2, (height-strHeight)/2
		if strHeight+topPadding+bottomPadding > height {
			topPadding--
		} else if strHeight+topPadding+bottomPadding < height {
			bottomPadding++
		}
		return strings.Repeat("\n", topPadding) + str + strings.Repeat("\n", bottomPadding)
	case Bottom:
		return strings.Repeat("\n", height-strHeight) + str
	default:
		return str
	}
}
