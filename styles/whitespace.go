package styles

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/scuf"
)

// whitespace is a whitespace renderer.
type whitespace struct {
	re    *Renderer
	style []scuf.Modifier
	chars string
}

// newWhitespace creates a new whitespace renderer. The order of the options
// matters, if you're using WithWhitespaceRenderer, make sure it comes first as
// other options might depend on it.
func newWhitespace(r *Renderer, opts ...WhitespaceOption) *whitespace {
	w := &whitespace{
		re:    r,
		style: nil,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Render whitespaces.
func (w whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	r := []rune(w.chars)
	j := 0
	b := strings.Builder{}

	// Cycle through runes and print them into the whitespace.
	for i := 0; i < width; {
		b.WriteRune(r[j])
		j++
		if j >= len(r) {
			j = 0
		}
		i += ansi.PrintableRuneWidth(string(r[j]))
	}

	// Fill any extra gaps white spaces. This might be necessary if any runes
	// are more than one cell wide, which could leave a one-rune gap.
	short := width - ansi.PrintableRuneWidth(b.String())
	if short > 0 {
		b.WriteString(strings.Repeat(" ", short))
	}

	return scuf.String(b.String(), w.style...)
}

// WhitespaceOption sets a styling rule for rendering whitespace.
type WhitespaceOption func(*whitespace)

// WithWhitespaceForeground sets the color of the characters in the whitespace.
func WithWhitespaceForeground(fg scuf.Modifier) WhitespaceOption {
	return func(w *whitespace) {
		w.style = append(w.style, fg)
	}
}

// WithWhitespaceBackground sets the background color of the whitespace.
func WithWhitespaceBackground(bg scuf.Modifier) WhitespaceOption {
	return WithWhitespaceForeground(bg)
}

// WithWhitespaceChars sets the characters to be rendered in the whitespace.
func WithWhitespaceChars(s string) WhitespaceOption {
	return func(w *whitespace) {
		w.chars = s
	}
}
