package box

// Border contains a series of values which comprise the various parts of a border.
type Border struct {
	Top         rune
	Bottom      rune
	Left        rune
	Right       rune
	TopLeft     rune
	TopRight    rune
	BottomRight rune
	BottomLeft  rune
}

var (
	noBorder = Border{}

	// NormalBorder is a standard-type border with a normal weight and 90 degree corners.
	NormalBorder = Border{
		Top:         '─',
		Bottom:      '─',
		Left:        '│',
		Right:       '│',
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
	}

	// RoundedBorder is a border with rounded corners.
	RoundedBorder = Border{
		Top:         '─',
		Bottom:      '─',
		Left:        '│',
		Right:       '│',
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
	}

	// BlockBorder is a border that takes the whole block.
	BlockBorder = Border{
		Top:         '█',
		Bottom:      '█',
		Left:        '█',
		Right:       '█',
		TopLeft:     '█',
		TopRight:    '█',
		BottomLeft:  '█',
		BottomRight: '█',
	}

	// OuterHalfBlockBorder is a half-block border that sits outside the frame.
	OuterHalfBlockBorder = Border{
		Top:         '▀',
		Bottom:      '▄',
		Left:        '▌',
		Right:       '▐',
		TopLeft:     '▛',
		TopRight:    '▜',
		BottomLeft:  '▙',
		BottomRight: '▟',
	}

	// InnerHalfBlockBorder is a half-block border that sits inside the frame.
	InnerHalfBlockBorder = Border{
		Top:         '▄',
		Bottom:      '▀',
		Left:        '▐',
		Right:       '▌',
		TopLeft:     '▗',
		TopRight:    '▖',
		BottomLeft:  '▝',
		BottomRight: '▘',
	}

	// ThickBorder is a border that's thicker than the one returned by NormalBorder.
	ThickBorder = Border{
		Top:         '━',
		Bottom:      '━',
		Left:        '┃',
		Right:       '┃',
		TopLeft:     '┏',
		TopRight:    '┓',
		BottomLeft:  '┗',
		BottomRight: '┛',
	}

	// DoubleBorder is a border comprised of two thin strokes.
	DoubleBorder = Border{
		Top:         '═',
		Bottom:      '═',
		Left:        '║',
		Right:       '║',
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
	}

	// HiddenBorder is a border that renders as a series of single-cell
	// spaces. It's useful for cases when you want to remove a standard border but
	// maintain layout positioning. This said, you can still apply a background
	// color to a hidden border.
	HiddenBorder = Border{
		Top:         ' ',
		Bottom:      ' ',
		Left:        ' ',
		Right:       ' ',
		TopLeft:     ' ',
		TopRight:    ' ',
		BottomLeft:  ' ',
		BottomRight: ' ',
	}
)
