package lipgloss

import (
	"maps"
	"strings"
	"unicode"

	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	termenv "github.com/rprtr258/col"
)

// Property for a key.
type propKey int

// Available properties.
const (
	_boldKey propKey = iota
	_italicKey
	_underlineKey
	_strikethroughKey
	_reverseKey
	blinkKey
	faintKey
	foregroundKey
	backgroundKey
	widthKey
	heightKey
	alignHorizontalKey
	alignVerticalKey

	colorWhitespaceKey

	// Border runes.
	borderStyleKey

	// Border edges.
	borderTopKey
	borderRightKey
	borderBottomKey
	borderLeftKey

	// Border foreground colors.
	borderTopForegroundKey
	borderRightForegroundKey
	borderBottomForegroundKey
	borderLeftForegroundKey

	// Border background colors.
	borderTopBackgroundKey
	borderRightBackgroundKey
	borderBottomBackgroundKey
	borderLeftBackgroundKey

	_inlineKey
	_maxWidthKey
	_maxHeightKey
	_underlineSpacesKey
	_strikethroughSpacesKey
)

// NewStyle returns a new, empty Style. While it's syntactic sugar for the
// Style{} primitive, it's recommended to use this function for creating styles
// in case the underlying implementation changes. It takes an optional string
// value to be set as the underlying string value for this style.
func NewStyle() Style {
	return _renderer.NewStyle()
}

// NewStyle returns a new, empty Style. While it's syntactic sugar for the
// Style{} primitive, it's recommended to use this function for creating styles
// in case the underlying implementation changes. It takes an optional string
// value to be set as the underlying string value for this style.
func (r *Renderer) NewStyle() Style {
	return Style{
		r:     r,
		rules: map[propKey]any{},
	}
}

// Style contains a set of rules that comprise a style as a whole.
type Style struct {
	r     *Renderer
	rules map[propKey]any
	value string
}

// joinString joins a list of strings into a single string separated with a
// space.
func joinString(strs ...string) string {
	return strings.Join(strs, " ")
}

// SetString sets the underlying string value for this style. To render once
// the underlying string is set, use the Style.String. This method is
// a convenience for cases when having a stringer implementation is handy, such
// as when using fmt.Sprintf. You can also simply define a style and render out
// strings directly with Style.Render.
func (s Style) SetString(strs ...string) Style {
	s.value = joinString(strs...)
	return s
}

// Value returns the raw, unformatted, underlying string value for this style.
func (s Style) Value() string {
	return s.value
}

// String implements stringer for a Style, returning the rendered result based
// on the rules in this style. An underlying string value must be set with
// Style.SetString prior to using this method.
func (s Style) String() string {
	return s.Render()
}

// Copy returns a copy of this style, including any underlying string values.
func (s Style) Copy() Style {
	o := NewStyle()
	maps.Copy(o.rules, s.rules)
	o.r = s.r
	o.value = s.value
	return o
}

// Inherit overlays the style in the argument onto this style by copying each explicitly
// set value from the argument style onto this style if it is not already explicitly set.
// Existing set values are kept intact and not overwritten.
//
// Margins, padding, and underlying string values are not inherited.
func (s Style) Inherit(i Style) Style {
	for k, v := range i.rules {
		if _, ok := s.rules[k]; ok {
			continue
		}

		s.rules[k] = v
	}
	return s
}

// Render applies the defined style formatting to a given string.
func (s Style) Render(strs ...string) string {
	if s.r == nil {
		s.r = _renderer
	}
	if s.value != "" {
		strs = append([]string{s.value}, strs...)
	}

	var (
		bold          = s.GetBold()
		italic        = s.GetItalic()
		underline     = s.GetUnderline()
		strikethrough = s.GetStrikethrough()
		reverse       = s.GetReverse()
		blink         = s.GetBlink()
		faint         = s.GetFaint()

		fg = s.GetForeground()
		bg = s.GetBackground()

		width           = s.GetWidth()
		height          = s.GetHeight()
		horizontalAlign = s.GetAlignHorizontal()
		verticalAlign   = s.GetAlignVertical()

		// topPadding    = s.GetPaddingTop()
		// rightPadding  = s.GetPaddingRight()
		// bottomPadding = s.GetPaddingBottom()
		// leftPadding   = s.GetPaddingLeft()

		colorWhitespace = s.getAsBool(colorWhitespaceKey, true)
		inline          = s.GetInline()
		maxWidth        = s.GetMaxWidth()
		maxHeight       = s.GetMaxHeight()

		underlineSpaces     = underline && s.getAsBool(_underlineSpacesKey, true)
		strikethroughSpaces = strikethrough && s.getAsBool(_strikethroughSpacesKey, true)
	)

	str := joinString(strs...)

	if len(s.rules) == 0 {
		return str
	}

	// Enable support for ANSI on the legacy Windows cmd.exe console. This is a
	// no-op on non-Windows systems and on Windows runs only once.
	enableLegacyWindowsANSI()

	te := s.r.ColorProfile().S()
	if bold {
		te = te.Bold()
	}
	if italic {
		te = te.Italic()
	}
	if underline {
		te = te.Underline()
	}
	teWhitespace := s.r.ColorProfile().S()
	if reverse {
		teWhitespace = teWhitespace.Reverse()
		te = te.Reverse()
	}
	if blink {
		te = te.Blink()
	}
	if faint {
		te = te.Faint()
	}

	// Do we need to style spaces separately?
	useSpaceStyler := underlineSpaces || strikethroughSpaces

	// Do we need to style whitespace (padding and space outside
	// paragraphs) separately?
	styleWhitespace := reverse

	teSpace := s.r.ColorProfile().S()
	if fg != noColor {
		te = te.Foreground(fg.color(s.r))
		if styleWhitespace {
			teWhitespace = teWhitespace.Foreground(fg.color(s.r))
		}
		if useSpaceStyler {
			teSpace = teSpace.Foreground(fg.color(s.r))
		}
	}

	if bg != noColor {
		te = te.Background(bg.color(s.r))
		if colorWhitespace {
			teWhitespace = teWhitespace.Background(bg.color(s.r))
		}
		if useSpaceStyler {
			teSpace = teSpace.Background(bg.color(s.r))
		}
	}

	if underline {
		te = te.Underline()
	}
	if strikethrough {
		te = te.CrossOut()
	}

	if underlineSpaces {
		teSpace = teSpace.Underline()
	}
	if strikethroughSpaces {
		teSpace = teSpace.CrossOut()
	}

	// Strip newlines in single line mode
	if inline {
		str = strings.ReplaceAll(str, "\n", "")
	}

	// Word wrap
	if !inline && width > 0 {
		wrapAt := width //- leftPadding - rightPadding
		str = wordwrap.String(str, wrapAt)
		str = wrap.String(str, wrapAt) // force-wrap long strings
	}

	// Render core text
	{
		var sb strings.Builder

		lines := strings.Split(str, "\n")
		for i, line := range lines {
			if useSpaceStyler {
				// Look for spaces and apply a different styler
				for _, r := range line {
					runeStyle := te
					if unicode.IsSpace(r) {
						runeStyle = teSpace
					}

					sb.WriteString(runeStyle.Render(string(r)))
				}
			} else {
				sb.WriteString(te.Render(line))
			}

			if i != len(lines)-1 {
				sb.WriteRune('\n')
			}
		}

		str = sb.String()
	}

	// Height
	if height > 0 {
		str = alignTextVertical(str, verticalAlign, height)
	}

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		numLines := strings.Count(str, "\n")

		if numLines != 0 || width != 0 {
			st := termenv.S()
			if colorWhitespace || styleWhitespace {
				st = teWhitespace
			}

			str = alignTextHorizontal(str, horizontalAlign, width, st)
		}
	}

	if !inline {
		str = s.applyBorder(str)
	}

	// Truncate according to MaxHeight
	if maxHeight > 0 {
		lines := strings.Split(str, "\n")
		str = strings.Join(lines[:min(maxHeight, len(lines))], "\n")
	}

	// Truncate according to MaxWidth
	if maxWidth > 0 {
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			lines[i] = truncate.String(line, uint(maxWidth))
		}

		str = strings.Join(lines, "\n")
	}

	return str
}
