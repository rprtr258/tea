package styles

import (
	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/scuf"
)

// GetBold returns the style's bold value. If no value is set false is returned.
func (s Style) GetBold() bool {
	return s.bold
}

// GetItalic returns the style's italic value. If no value is set false is
// returned.
func (s Style) GetItalic() bool {
	return s.getAsBool(_keyItalic, false)
}

// GetUnderline returns the style's underline value. If no value is set false is
// returned.
func (s Style) GetUnderline() bool {
	return s.getAsBool(_keyUnderline, false)
}

// GetStrikethrough returns the style's strikethrough value. If no value is set false
// is returned.
func (s Style) GetStrikethrough() bool {
	return s.getAsBool(_keyStrikethrough, false)
}

// GetReverse returns the style's reverse value. If no value is set false is
// returned.
func (s Style) GetReverse() bool {
	return s.getAsBool(_keyReverse, false)
}

// GetBlink returns the style's blink value. If no value is set false is
// returned.
func (s Style) GetBlink() bool {
	return s.getAsBool(_keyBlink, false)
}

// GetFaint returns the style's faint value. If no value is set false is
// returned.
func (s Style) GetFaint() bool {
	return s.getAsBool(_keyFaint, false)
}

// GetForeground returns the style's foreground color. If no value is set
// NoColor{} is returned.
func (s Style) GetForeground() scuf.Modifier {
	return s.getAsColor(_keyForeground)
}

// GetBackground returns the style's background color. If no value is set
// NoColor{} is returned.
func (s Style) GetBackground() scuf.Modifier {
	return s.getAsColor(_keyBackground)
}

// GetAlign returns the style's implicit horizontal alignment setting.
// If no alignment is set Position.Left is returned.
func (s Style) GetAlign() Alignment {
	v := s.getAsPosition(_keyAlighHorizontal)
	if v == Alignment(0) {
		return Left
	}
	return v
}

// GetAlignHorizontal returns the style's implicit horizontal alignment setting.
// If no alignment is set Position.Left is returned.
func (s Style) GetAlignHorizontal() Alignment {
	v := s.getAsPosition(_keyAlighHorizontal)
	if v == Alignment(0) {
		return Left
	}
	return v
}

// GetAlignVertical returns the style's implicit vertical alignment setting.
// If no alignment is set Position.Top is returned.
func (s Style) GetAlignVertical() Alignment {
	return s.getAsPosition(_keyAlignVertical)
}

// GetColorWhitespace returns the style's whitespace coloring setting. If no
// value is set false is returned.
func (s Style) GetColorWhitespace() bool {
	return s.getAsBool(_keyColorWhitespace, false)
}

// GetUnderlineSpaces returns whether or not the style is set to underline
// spaces. If not value is set false is returned.
func (s Style) GetUnderlineSpaces() bool {
	return s.getAsBool(_keyUnderlineSpaces, false)
}

// GetStrikethroughSpaces returns whether or not the style is set to strikethrough
// spaces. If not value is set false is returned.
func (s Style) GetStrikethroughSpaces() bool {
	return s.getAsBool(_keyStrikethroughSpaces, false)
}

// Returns whether or not the given property is set.
func (s Style) isSet(k propKey) bool {
	_, exists := s.rules[k]
	return exists
}

func (s Style) getAsBool(k propKey, defaultVal bool) bool {
	v, ok := s.rules[k]
	if !ok {
		return defaultVal
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return defaultVal
}

func (s Style) getAsColor(k propKey) scuf.Modifier {
	v, ok := s.rules[k]
	if !ok {
		return nil
	}

	c, ok := v.(scuf.Modifier)
	if !ok {
		return nil
	}

	return c
}

func (s Style) getAsInt(k propKey) int {
	v, ok := s.rules[k]
	if !ok {
		return 0
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func (s Style) getAsPosition(k propKey) Alignment {
	v, ok := s.rules[k]
	if !ok {
		return Alignment(0)
	}
	if p, ok := v.(Alignment); ok {
		return p
	}
	return Alignment(0)
}

// getWidestWidth returns size of the widest line
func getWidestWidth(lines []string) int {
	widest := 0
	for _, l := range lines {
		widest = max(widest, ansi.PrintableRuneWidth(l))
	}
	return widest
}
