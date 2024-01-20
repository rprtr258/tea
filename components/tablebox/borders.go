package tablebox

import "github.com/rprtr258/tea/components/box"

type FullBorder struct {
	box.Border
	TRB  rune
	TLB  rune
	LRT  rune
	LRB  rune
	LRTB rune
}

var (
	noBorder = FullBorder{}

	// NormalBorder is a standard-type border with a normal weight and 90 degree corners.
	NormalBorder = FullBorder{
		Border: box.NormalBorder,
		TRB:    '├',
		TLB:    '┤',
		LRT:    '┬',
		LRB:    '┴',
		LRTB:   '┼',
	}

	// RoundedBorder is a border with rounded corners.
	RoundedBorder = FullBorder{
		Border: box.RoundedBorder,
		TRB:    '├',
		TLB:    '┤',
		LRT:    '┬',
		LRB:    '┴',
		LRTB:   '┼',
	}

	// BlockBorder is a border that takes the whole block.
	BlockBorder = FullBorder{
		Border: box.BlockBorder,
		TLB:    '█',
		TRB:    '█',
		LRT:    '█',
		LRB:    '█',
		LRTB:   '█',
	}

	// HiddenBorder is a border that renders as a series of single-cell
	// spaces. It's useful for cases when you want to remove a standard border but
	// maintain layout positioning. This said, you can still apply a background
	// color to a hidden border.
	HiddenBorder = FullBorder{
		Border: box.HiddenBorder,
		TRB:    ' ',
		TLB:    ' ',
		LRT:    ' ',
		LRB:    ' ',
		LRTB:   ' ',
	}
)
