package styles

import (
	"maps"
	"strings"
	"unicode"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
)

// Property for a key.
type propKey int

// Available properties.
const (
	_keyBold propKey = iota
	_keyItalic
	_keyUnderline
	_keyStrikethrough
	_keyReverse
	_keyBlink
	_keyFaint
	_keyForeground
	_keyBackground
	_keyAlighHorizontal
	_keyAlignVertical

	_keyColorWhitespace

	_keyUnderlineSpaces
	_keyStrikethroughSpaces
)

// Style contains a set of rules that comprise a style as a whole.
type Style struct {
	value string
	rules map[propKey]any
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
	o := Style{rules: map[propKey]any{}}
	maps.Copy(o.rules, s.rules)
	o.value = s.value
	return o
}

// Inherit overlays the style in the argument onto this style by copying each explicitly
// set value from the argument style onto this style if it is not already explicitly set.
// Existing set values are kept intact and not overwritten.
//
// Margins, padding, and underlying string values are not inherited.
func (s Style) Inherit(i Style) Style {
	if s.rules == nil {
		s.rules = map[propKey]any{}
	}
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

		// horizontalAlign = s.GetAlignHorizontal()
		// verticalAlign   = s.GetAlignVertical()

		colorWhitespace = s.getAsBool(_keyColorWhitespace, true)

		underlineSpaces     = underline && s.getAsBool(_keyUnderlineSpaces, true)
		strikethroughSpaces = strikethrough && s.getAsBool(_keyStrikethroughSpaces, true)
	)

	str := joinString(strs...)

	if len(s.rules) == 0 {
		return str
	}

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
	if len(fg) != 0 {
		te = append(te, fg)
		if styleWhitespace {
			teWhitespace = append(teWhitespace, fg)
		}
		if useSpaceStyler {
			teSpace = append(teSpace, fg)
		}
	}

	if len(bg) != 0 {
		te = append(te, bg)
		if colorWhitespace {
			teWhitespace = append(teWhitespace, bg)
		}
		if useSpaceStyler {
			teSpace = append(teSpace, bg)
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
	str = strings.ReplaceAll(str, "\n", "")

	// Word wrap
	// if !inline && width > 0 {
	// 	wrapAt := width //- leftPadding - rightPadding
	// 	str = wordwrap.String(str, wrapAt)
	// 	str = wrap.String(str, wrapAt) // force-wrap long strings
	// }

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
	// if height > 0 {
	// 	str = alignTextVertical(str, verticalAlign, height)
	// }

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	// {
	// 	numLines := strings.Count(str, "\n")

	// 	if numLines != 0 || width != 0 {
	// 		st := fun.IF(colorWhitespace || styleWhitespace, teWhitespace, nil)
	// 		str = alignTextHorizontal(str, horizontalAlign, width, st)
	// 	}
	// }

	// Truncate according to MaxHeight
	// if maxHeight > 0 {
	// 	lines := strings.Split(str, "\n")
	// 	str = strings.Join(lines[:min(maxHeight, len(lines))], "\n")
	// }

	// // Truncate according to MaxWidth
	// if maxWidth > 0 {
	// 	lines := strings.Split(str, "\n")
	// 	for i, line := range lines {
	// 		lines[i] = truncate.String(line, uint(maxWidth))
	// 	}

	// 	str = strings.Join(lines, "\n")
	// }

	return str
}
