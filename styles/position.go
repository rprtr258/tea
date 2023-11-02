package styles

import (
	"math"
	"strings"

	"github.com/muesli/reflow/ansi"
)

// Alignment represents a position along a horizontal or vertical axis. It's in
// situations where an axis is involved, like alignment, joining, placement and
// so on.
//
// A value of 0 represents the start (the left or top) and 1 represents the end
// (the right or bottom). 0.5 represents the center.
//
// There are constants Top, Bottom, Center, Left and Right in this package that
// can be used to aid readability.
type Alignment float64

func (p Alignment) value() float64 {
	return min(1, max(0, float64(p)))
}

// Position aliases.
const (
	Top    Alignment = 0.0
	Bottom Alignment = 1.0
	Center Alignment = 0.5
	Left   Alignment = 0.0
	Right  Alignment = 1.0
)

// Place places a string or text block vertically in an unstyled box of a given
// width or height.
func Place(width, height int, hPos, vPos Alignment, str string, opts ...WhitespaceOption) string {
	return _renderer.Place(width, height, hPos, vPos, str, opts...)
}

// Place places a string or text block vertically in an unstyled box of a given
// width or height.
func (r *Renderer) Place(width, height int, hPos, vPos Alignment, str string, opts ...WhitespaceOption) string {
	return r.PlaceVertical(height, vPos, r.PlaceHorizontal(width, hPos, str, opts...), opts...)
}

// PlaceHorizontal places a string or text block horizontally in an unstyled
// block of a given width. If the given width is shorter than the max width of
// the string (measured by its longest line) this will be a noop.
func PlaceHorizontal(width int, pos Alignment, str string, opts ...WhitespaceOption) string {
	return _renderer.PlaceHorizontal(width, pos, str, opts...)
}

// PlaceHorizontal places a string or text block horizontally in an unstyled
// block of a given width. If the given width is shorter than the max width of
// the string (measured by its longest line) this will be a noöp.
func (r *Renderer) PlaceHorizontal(width int, pos Alignment, str string, opts ...WhitespaceOption) string {
	lines := strings.Split(str, "\n")
	contentWidth := getWidestWidth(lines)
	gap := width - contentWidth

	if gap <= 0 {
		return str
	}

	ws := newWhitespace(r, opts...)

	var sb strings.Builder
	for i, line := range lines {
		// Is this line shorter than the longest line?
		short := max(0, contentWidth-ansi.PrintableRuneWidth(line))
		totalGap := gap + short

		switch pos {
		case Left:
			sb.WriteString(line)
			sb.WriteString(ws.render(totalGap))
		case Right:
			sb.WriteString(ws.render(totalGap))
			sb.WriteString(line)
		default: // somewhere in the middle
			split := int(math.Round(float64(totalGap) * pos.value()))
			left := totalGap - split
			sb.WriteString(ws.render(left))

			sb.WriteString(line)

			right := totalGap - left
			sb.WriteString(ws.render(right))
		}

		if i < len(lines)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

// PlaceVertical places a string or text block vertically in an unstyled block
// of a given height. If the given height is shorter than the height of the
// string (measured by its newlines) then this will be a noop.
func PlaceVertical(height int, pos Alignment, str string, opts ...WhitespaceOption) string {
	return _renderer.PlaceVertical(height, pos, str, opts...)
}

// PlaceVertical places a string or text block vertically in an unstyled block
// of a given height. If the given height is shorter than the height of the
// string (measured by its newlines) then this will be a noöp.
func (r *Renderer) PlaceVertical(height int, pos Alignment, str string, opts ...WhitespaceOption) string {
	contentHeight := strings.Count(str, "\n") + 1
	gap := height - contentHeight

	if gap <= 0 {
		return str
	}

	ws := newWhitespace(r, opts...)

	width := getWidestWidth(strings.Split(str, "\n"))
	emptyLine := ws.render(width)
	b := strings.Builder{}

	switch pos {
	case Top:
		b.WriteString(str)
		b.WriteRune('\n')
		for i := 0; i < gap; i++ {
			b.WriteString(emptyLine)
			if i < gap-1 {
				b.WriteRune('\n')
			}
		}

	case Bottom:
		b.WriteString(strings.Repeat(emptyLine+"\n", gap))
		b.WriteString(str)

	default: // Somewhere in the middle
		split := int(math.Round(float64(gap) * pos.value()))
		top := gap - split
		bottom := gap - top

		b.WriteString(strings.Repeat(emptyLine+"\n", top))
		b.WriteString(str)

		for i := 0; i < bottom; i++ {
			b.WriteRune('\n')
			b.WriteString(emptyLine)
		}
	}

	return b.String()
}
