package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/termenv"
)

// Perform text alignment. If the string is multi-lined, we also make all lines
// the same width by padding them with spaces. If a termenv style is passed,
// use that to style the spaces added.
func alignTextHorizontal(str string, pos Position, width int, style termenv.Style) string {
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
				line = style.Styled(strings.Repeat(" ", shortAmount)) + line
			case Center:
				left := shortAmount / 2
				right := left + shortAmount%2 // note that we put the remainder on the right

				leftSpaces := style.Styled(strings.Repeat(" ", left))
				rightSpaces := style.Styled(strings.Repeat(" ", right))

				line = leftSpaces + line + rightSpaces
			default: // Left
				line += style.Styled(strings.Repeat(" ", shortAmount))
			}
		}

		sb.WriteString(line)
		if i < len(lines)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func alignTextVertical(str string, pos Position, height int) string {
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
