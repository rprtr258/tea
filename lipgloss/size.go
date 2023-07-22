package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

// Width returns the cell width of characters in the string. ANSI sequences are
// ignored and characters wider than one cell (such as Chinese characters and
// emojis) are appropriately measured.
//
// You should use this instead of len(string) len([]rune(string) as neither
// will give you accurate results.
func Width(str string) int {
	width := 0
	for _, l := range strings.Split(str, "\n") {
		w := ansi.PrintableRuneWidth(l)
		if w > width {
			width = w
		}
	}
	return width
}

// Height returns height of a string in cells. This is done simply by
// counting \n characters. If your strings use \r\n for newlines you should
// convert them to \n first, or simply write a separate function for measuring
// height.
func Height(str string) int {
	return strings.Count(str, "\n") + 1
}
