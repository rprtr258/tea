package lipgloss

import (
	"maps"
	"strings"
	"unicode"

	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
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

// TODO: make func
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

	te := []scuf.Modifier{}
	if bold {
		te = append(te, scuf.ModBold)
	}
	if italic {
		te = append(te, scuf.ModItalic)
	}
	if underline {
		te = append(te, scuf.ModUnderline)
	}
	teWhitespace := []scuf.Modifier{}
	if reverse {
		teWhitespace = append(teWhitespace, scuf.ModReverse)
		te = append(te, scuf.ModReverse)
	}
	if blink {
		te = append(te, scuf.ModBlink)
	}
	if faint {
		te = append(te, scuf.ModFaint)
	}

	// Do we need to style spaces separately?
	useSpaceStyler := underlineSpaces || strikethroughSpaces

	// Do we need to style whitespace (padding and space outside
	// paragraphs) separately?
	styleWhitespace := reverse

	teSpace := []scuf.Modifier{}
	if fg != nil {
		te = append(te, fg.color())
		if styleWhitespace {
			teWhitespace = append(teWhitespace, fg.color())
		}
		if useSpaceStyler {
			teSpace = append(teSpace, fg.color())
		}
	}

	if bg != nil {
		te = append(te, bg.color())
		if colorWhitespace {
			teWhitespace = append(teWhitespace, bg.color())
		}
		if useSpaceStyler {
			teSpace = append(teSpace, bg.color())
		}
	}

	if underline {
		te = append(te, scuf.ModUnderline)
	}
	if strikethrough {
		te = append(te, scuf.ModCrossout)
	}

	if underlineSpaces {
		teSpace = append(teSpace, scuf.ModUnderline)
	}
	if strikethroughSpaces {
		teSpace = append(teSpace, scuf.ModCrossout)
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
					runeStyle := fun.IF(unicode.IsSpace(r), teSpace, te)
					sb.WriteString(scuf.String(string(r), runeStyle...))
				}
			} else {
				sb.WriteString(scuf.String(line, te...))
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
			st := fun.IF(colorWhitespace || styleWhitespace, teWhitespace, nil)
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
