package styles

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	s.bold = false
	return s
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	s.italic = false
	return s
}

// UnsetUnderline removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	s.underline = false
	return s
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	s.strikethrough = false
	return s
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	s.reverse = false
	return s
}

// UnsetBlink removes the blink style rule, if set.
func (s Style) UnsetBlink() Style {
	s.blink = false
	return s
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	s.faint = false
	return s
}

// UnsetForeground removes the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	s.foreground = nil
	return s
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	s.background = nil
	return s
}

// UnsetAlign removes the horizontal and vertical text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	delete(s.rules, _keyAlighHorizontal)
	delete(s.rules, _keyAlignVertical)
	return s
}

// UnsetAlignHorizontal removes the horizontal text alignment style rule, if set.
func (s Style) UnsetAlignHorizontal() Style {
	delete(s.rules, _keyAlighHorizontal)
	return s
}

// UnsetAlignVertical removes the vertical text alignment style rule, if set.
func (s Style) UnsetAlignVertical() Style {
	delete(s.rules, _keyAlignVertical)
	return s
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	delete(s.rules, _keyColorWhitespace)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnderlineSpaces.
func (s Style) UnsetUnderlineSpaces() Style {
	delete(s.rules, _keyUnderlineSpaces)
	return s
}

// UnsetStrikethroughSpaces removes the value set by StrikethroughSpaces.
func (s Style) UnsetStrikethroughSpaces() Style {
	delete(s.rules, _keyStrikethroughSpaces)
	return s
}

// UnsetString sets the underlying string value to the empty string.
func (s Style) UnsetString() Style {
	s.value = ""
	return s
}
