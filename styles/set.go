package styles

import "github.com/rprtr258/scuf"

// Set a value on the underlying rules map.
func (s *Style) set(key propKey, value any) Style {
	if s.rules == nil {
		s.rules = map[propKey]any{}
	}
	switch v := value.(type) {
	case int:
		// We don't allow negative integers on any of our values, so just keep
		// them at zero or above. We could use uints instead, but the
		// conversions are a little tedious, so we're sticking with ints for
		// sake of usability.
		s.rules[key] = max(0, v)
	default:
		s.rules[key] = v
	}
	return *s
}

// Bold sets a bold formatting rule.
func (s Style) Bold(v bool) Style {
	return s.set(_keyBold, v)
}

// Italic sets an italic formatting rule. In some terminal emulators this will
// render with "reverse" coloring if not italic font variant is available.
func (s Style) Italic() Style {
	return s.set(_keyItalic, true)
}

// Underline sets an underline rule. By default, underlines will not be drawn on
// whitespace like margins and padding. To change this behavior set
// UnderlineSpaces.
func (s Style) Underline() Style {
	return s.set(_keyUnderline, true)
}

// Strikethrough sets a strikethrough rule. By default, strikes will not be
// drawn on whitespace like margins and padding. To change this behavior set
// StrikethroughSpaces.
func (s Style) Strikethrough(v bool) Style {
	return s.set(_keyStrikethrough, v)
}

// Reverse sets a rule for inverting foreground and background colors.
func (s Style) Reverse(v bool) Style {
	return s.set(_keyReverse, v)
}

// Blink sets a rule for blinking foreground text.
func (s Style) Blink() Style {
	return s.set(_keyBlink, true)
}

// Faint sets a rule for rendering the foreground color in a dimmer shade.
func (s Style) Faint() Style {
	return s.set(_keyFaint, true)
}

// Foreground sets a foreground color.
//
//	// Sets the foreground to blue
//	s := styles.Style{}.Foreground(styles.Color("#0000ff"))
//
//	// Removes the foreground color
//	s.Foreground(styles.NoColor)
func (s Style) Foreground(c scuf.Modifier) Style {
	return s.set(_keyForeground, c)
}

// Background sets a background color.
func (s Style) Background(c scuf.Modifier) Style {
	return s.set(_keyBackground, c)
}

// Align is a shorthand method for setting horizontal and vertical alignment.
//
// With one argument, the position value is applied to the horizontal alignment.
//
// With two arguments, the value is applied to the vertical and horizontal
// alignments, in that order.
func (s Style) Align(p ...Alignment) Style {
	if len(p) > 0 {
		s.set(_keyAlighHorizontal, p[0])
	}
	if len(p) > 1 {
		s.set(_keyAlignVertical, p[1])
	}
	return s
}

// AlignHorizontal sets a horizontal text alignment rule.
func (s Style) AlignHorizontal(p Alignment) Style {
	return s.set(_keyAlighHorizontal, p)
}

// AlignVertical sets a vertical text alignment rule.
func (s Style) AlignVertical(p Alignment) Style {
	return s.set(_keyAlignVertical, p)
}

// ColorWhitespace determines whether or not the background color should be
// applied to the padding. This is true by default as it's more than likely the
// desired and expected behavior, but it can be disabled for certain graphic
// effects.
func (s Style) ColorWhitespace(v bool) Style {
	return s.set(_keyColorWhitespace, v)
}

// UnderlineSpaces determines whether to underline spaces between words. By
// default, this is true. Spaces can also be underlined without underlining the
// text itself.
func (s Style) UnderlineSpaces(v bool) Style {
	return s.set(_keyUnderlineSpaces, v)
}

// StrikethroughSpaces determines whether to apply strikethroughs to spaces
// between words. By default, this is true. Spaces can also be struck without
// underlining the text itself.
func (s Style) StrikethroughSpaces(v bool) Style {
	return s.set(_keyStrikethroughSpaces, v)
}

// whichSidesInt is a helper method for setting values on sides of a block based
// on the number of arguments. It follows the CSS shorthand rules for blocks
// like margin, padding. and borders. Here are how the rules work:
//
// 0 args:  do nothing
// 1 arg:   all sides
// 2 args:  top -> bottom
// 3 args:  top -> horizontal -> bottom
// 4 args:  top -> right -> bottom -> left
// 5+ args: do nothing.
func whichSidesInt(i ...int) (top, right, bottom, left int, ok bool) { //nolint:nonamedreturns
	switch len(i) {
	case 1:
		top = i[0]
		bottom = i[0]
		left = i[0]
		right = i[0]
		ok = true
	case 2:
		top = i[0]
		bottom = i[0]
		left = i[1]
		right = i[1]
		ok = true
	case 3:
		top = i[0]
		left = i[1]
		right = i[1]
		bottom = i[2]
		ok = true
	case 4:
		top = i[0]
		right = i[1]
		bottom = i[2]
		left = i[3]
		ok = true
	}
	return top, right, bottom, left, ok
}

// whichSidesColor is like whichSides, except it operates on a series of
// boolean values. See the comment on whichSidesInt for details on how this
// works.
func whichSidesColor(i ...scuf.Modifier) (top, right, bottom, left scuf.Modifier, ok bool) { //nolint:nonamedreturns
	switch len(i) {
	case 1:
		top = i[0]
		bottom = i[0]
		left = i[0]
		right = i[0]
		ok = true
	case 2:
		top = i[0]
		bottom = i[0]
		left = i[1]
		right = i[1]
		ok = true
	case 3:
		top = i[0]
		left = i[1]
		right = i[1]
		bottom = i[2]
		ok = true
	case 4:
		top = i[0]
		right = i[1]
		bottom = i[2]
		left = i[3]
		ok = true
	}
	return top, right, bottom, left, ok
}
