package lipgloss

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	delete(s.rules, _boldKey)
	return s
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	delete(s.rules, _italicKey)
	return s
}

// UnsetUnderline removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	delete(s.rules, _underlineKey)
	return s
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	delete(s.rules, _strikethroughKey)
	return s
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	delete(s.rules, _reverseKey)
	return s
}

// UnsetBlink removes the blink style rule, if set.
func (s Style) UnsetBlink() Style {
	delete(s.rules, blinkKey)
	return s
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	delete(s.rules, faintKey)
	return s
}

// UnsetForeground removes the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	delete(s.rules, foregroundKey)
	return s
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	delete(s.rules, backgroundKey)
	return s
}

// UnsetWidth removes the width style rule, if set.
func (s Style) UnsetWidth() Style {
	delete(s.rules, widthKey)
	return s
}

// UnsetHeight removes the height style rule, if set.
func (s Style) UnsetHeight() Style {
	delete(s.rules, heightKey)
	return s
}

// UnsetAlign removes the horizontal and vertical text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	delete(s.rules, alignHorizontalKey)
	delete(s.rules, alignVerticalKey)
	return s
}

// UnsetAlignHorizontal removes the horizontal text alignment style rule, if set.
func (s Style) UnsetAlignHorizontal() Style {
	delete(s.rules, alignHorizontalKey)
	return s
}

// UnsetAlignVertical removes the vertical text alignment style rule, if set.
func (s Style) UnsetAlignVertical() Style {
	delete(s.rules, alignVerticalKey)
	return s
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	delete(s.rules, colorWhitespaceKey)
	return s
}

// UnsetBorderStyle removes the border style rule, if set.
func (s Style) UnsetBorderStyle() Style {
	delete(s.rules, borderStyleKey)
	return s
}

// UnsetBorderTop removes the border top style rule, if set.
func (s Style) UnsetBorderTop() Style {
	delete(s.rules, borderTopKey)
	return s
}

// UnsetBorderRight removes the border right style rule, if set.
func (s Style) UnsetBorderRight() Style {
	delete(s.rules, borderRightKey)
	return s
}

// UnsetBorderBottom removes the border bottom style rule, if set.
func (s Style) UnsetBorderBottom() Style {
	delete(s.rules, borderBottomKey)
	return s
}

// UnsetBorderLeft removes the border left style rule, if set.
func (s Style) UnsetBorderLeft() Style {
	delete(s.rules, borderLeftKey)
	return s
}

// UnsetBorderForeground removes all border foreground color styles, if set.
func (s Style) UnsetBorderForeground() Style {
	delete(s.rules, borderTopForegroundKey)
	delete(s.rules, borderRightForegroundKey)
	delete(s.rules, borderBottomForegroundKey)
	delete(s.rules, borderLeftForegroundKey)
	return s
}

// UnsetBorderTopForeground removes the top border foreground color rule,
// if set.
func (s Style) UnsetBorderTopForeground() Style {
	delete(s.rules, borderTopForegroundKey)
	return s
}

// UnsetBorderRightForeground removes the right border foreground color rule,
// if set.
func (s Style) UnsetBorderRightForeground() Style {
	delete(s.rules, borderRightForegroundKey)
	return s
}

// UnsetBorderBottomForeground removes the bottom border foreground color
// rule, if set.
func (s Style) UnsetBorderBottomForeground() Style {
	delete(s.rules, borderBottomForegroundKey)
	return s
}

// UnsetBorderLeftForeground removes the left border foreground color rule,
// if set.
func (s Style) UnsetBorderLeftForeground() Style {
	delete(s.rules, borderLeftForegroundKey)
	return s
}

// UnsetBorderBackground removes all border background color styles, if
// set.
func (s Style) UnsetBorderBackground() Style {
	delete(s.rules, borderTopBackgroundKey)
	delete(s.rules, borderRightBackgroundKey)
	delete(s.rules, borderBottomBackgroundKey)
	delete(s.rules, borderLeftBackgroundKey)
	return s
}

// UnsetBorderTopBackgroundColor removes the top border background color rule,
// if set.
func (s Style) UnsetBorderTopBackgroundColor() Style {
	delete(s.rules, borderTopBackgroundKey)
	return s
}

// UnsetBorderRightBackground removes the right border background color
// rule, if set.
func (s Style) UnsetBorderRightBackground() Style {
	delete(s.rules, borderRightBackgroundKey)
	return s
}

// UnsetBorderBottomBackground removes the bottom border background color
// rule, if set.
func (s Style) UnsetBorderBottomBackground() Style {
	delete(s.rules, borderBottomBackgroundKey)
	return s
}

// UnsetBorderLeftBackground removes the left border color rule, if set.
func (s Style) UnsetBorderLeftBackground() Style {
	delete(s.rules, borderLeftBackgroundKey)
	return s
}

// UnsetInline removes the inline style rule, if set.
func (s Style) UnsetInline() Style {
	delete(s.rules, _inlineKey)
	return s
}

// UnsetMaxWidth removes the max width style rule, if set.
func (s Style) UnsetMaxWidth() Style {
	delete(s.rules, _maxWidthKey)
	return s
}

// UnsetMaxHeight removes the max height style rule, if set.
func (s Style) UnsetMaxHeight() Style {
	delete(s.rules, _maxHeightKey)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnderlineSpaces.
func (s Style) UnsetUnderlineSpaces() Style {
	delete(s.rules, _underlineSpacesKey)
	return s
}

// UnsetStrikethroughSpaces removes the value set by StrikethroughSpaces.
func (s Style) UnsetStrikethroughSpaces() Style {
	delete(s.rules, _strikethroughSpacesKey)
	return s
}

// UnsetString sets the underlying string value to the empty string.
func (s Style) UnsetString() Style {
	s.value = ""
	return s
}
