package styles

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
