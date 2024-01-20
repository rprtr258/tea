package styles

import (
	"math"
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/fun"
)

// JoinHorizontal is a utility function for horizontally joining two
// potentially multi-lined strings along a vertical axis. The first argument is
// the position, with 0 being all the way at the top and 1 being all the way
// at the bottom.
//
// If you just want to align to the left, right or center you may as well just
// use the helper constants Top, Center, and Bottom.
//
// Example:
//
//	blockB := "...\n...\n..."
//	blockA := "...\n...\n...\n...\n..."
//
//	// Join 20% from the top
//	str := styles.JoinHorizontal(0.2, blockA, blockB)
//
//	// Join on the top edge
//	str := styles.JoinHorizontal(styles.Top, blockA, blockB)
func JoinHorizontal(pos Alignment, blocks ...[]string) string { // TODO: instead do viewboxes layouts
	if len(blocks) == 0 {
		return ""
	}
	if len(blocks) == 1 {
		return strings.Join(blocks[0], "\n")
	}

	// Max line widths for the above text blocks
	maxWidths := fun.Map[int](getWidestWidth, blocks...)

	// Height of the tallest block
	var maxHeight int
	for _, block := range blocks {
		maxHeight = max(maxHeight, len(block))
	}

	// Add extra lines to make each side the same height
	for i := range blocks {
		if len(blocks[i]) >= maxHeight {
			continue
		}

		extraLines := make([]string, maxHeight-len(blocks[i]))

		switch pos {
		case Top:
			blocks[i] = append(blocks[i], extraLines...)

		case Bottom:
			blocks[i] = append(extraLines, blocks[i]...) //nolint:makezero

		default: // Somewhere in the middle
			n := len(extraLines)
			split := int(math.Round(float64(n) * pos.value()))
			top := n - split
			bottom := n - top

			blocks[i] = append(extraLines[top:], blocks[i]...)
			blocks[i] = append(blocks[i], extraLines[bottom:]...)
		}
	}

	// Merge lines
	var sb strings.Builder
	for i := range blocks[0] { // remember, all blocks have the same number of members now
		for j, block := range blocks {
			sb.WriteString(block[i])

			// Also make lines the same length
			sb.WriteString(strings.Repeat(" ", maxWidths[j]-ansi.PrintableRuneWidth(block[i])))
		}
		if i < len(blocks[0])-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

// JoinVertical is a utility function for vertically joining two potentially
// multi-lined strings along a horizontal axis. The first argument is the
// position, with 0 being all the way to the left and 1 being all the way to
// the right.
//
// If you just want to align to the left, right or center you may as well just
// use the helper constants Left, Center, and Right.
//
// Example:
//
//	blockB := "...\n...\n..."
//	blockA := "...\n...\n...\n...\n..."
//
//	// Join 20% from the top
//	str := styles.JoinVertical(0.2, blockA, blockB)
//
//	// Join on the right edge
//	str := styles.JoinVertical(styles.Right, blockA, blockB)
func JoinVertical(alignment Alignment, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	blocks := fun.Map[[]string](
		func(str string) []string {
			return strings.Split(str, "\n")
		}, strs...)

	var maxWidth int
	for _, block := range blocks {
		maxWidth = max(maxWidth, getWidestWidth(block))
	}

	var sb strings.Builder
	for i, block := range blocks {
		for j, line := range block {
			w := maxWidth - ansi.PrintableRuneWidth(line)

			switch alignment {
			case Left:
				sb.WriteString(line)
				sb.WriteString(strings.Repeat(" ", w))

			case Right:
				sb.WriteString(strings.Repeat(" ", w))
				sb.WriteString(line)

			default: // Somewhere in the middle
				if w < 1 {
					sb.WriteString(line)
					break
				}

				split := int(math.Round(float64(w) * alignment.value()))
				right := w - split
				left := w - right

				sb.WriteString(strings.Repeat(" ", left))
				sb.WriteString(line)
				sb.WriteString(strings.Repeat(" ", right))
			}

			// Write a newline as long as we're not on the last line of the last block.
			if i != len(blocks)-1 || j != len(block)-1 {
				sb.WriteRune('\n')
			}
		}
	}

	return sb.String()
}
