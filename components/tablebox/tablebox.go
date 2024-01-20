package tablebox

import (
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

// Box draws model inside a box
func Box(
	vb tea.Viewbox,
	h, w []tea.Layout,
	inside func(vb tea.Viewbox, y, x int),
	border FullBorder,
	style styles.Style,
) {
	hs := tea.EvalLayout(vb.Height-len(h)-1, h...)
	ws := tea.EvalLayout(vb.Width-len(w)-1, w...)

	keypointsX := make([]int, len(ws)+1)
	for i, w := range ws {
		keypointsX[i+1] = keypointsX[i] + w + 1
	}
	keypointsY := make([]int, len(hs)+1)
	for i, h := range hs {
		keypointsY[i+1] = keypointsY[i] + h + 1
	}

	for iy, h := range hs {
		for ix, w := range ws {
			inside(vb.Sub(tea.Rectangle{
				Top:    keypointsY[iy] + 1,
				Left:   keypointsX[ix] + 1,
				Height: h,
				Width:  w,
			}), iy, ix)
		}
	}

	if border == noBorder {
		return
	}

	// borders: LR, TB
	/// LR
	for _, y := range keypointsY {
		vb.
			Sub(tea.Rectangle{
				Top:    y,
				Left:   1,
				Height: 1,
				Width:  vb.Width - 2,
			}).
			Styled(style).
			Fill(border.Top) // TODO: on bottom should be border.Bottom, in the middle should be what?
	}
	/// LR
	for _, x := range keypointsX {
		vb.
			Sub(tea.Rectangle{
				Top:    1,
				Left:   x,
				Height: vb.Height - 2,
				Width:  1,
			}).
			Styled(style).
			Fill(border.Left) // TODO: analogical problem as with Top-Middle-Bottom
	}

	// 4 way conjunctions (full crosses): LRTB
	for iy := 0; iy < len(hs)-1; iy++ {
		for ix := 0; ix < len(ws)-1; ix++ {
			vb.Pixel(keypointsY[iy+1], keypointsX[ix+1]).Styled(style).Fill(border.LRTB)
		}
	}

	// 3 way conjunctions: LRT, LRB, TLB, TRB
	/// LRT
	for ix := 0; ix < len(ws)-1; ix++ {
		vb.Pixel(0, keypointsX[ix+1]).Styled(style).Fill(border.LRT)
	}
	/// LRB
	for ix := 0; ix < len(ws)-1; ix++ {
		vb.Pixel(vb.Height-1, keypointsX[ix+1]).Styled(style).Fill(border.LRB)
	}
	/// TRB
	for iy := 0; iy < len(hs)-1; iy++ {
		vb.Pixel(keypointsY[iy+1], 0).Styled(style).Fill(border.TRB)
	}
	/// TLB
	for iy := 0; iy < len(hs)-1; iy++ {
		vb.Pixel(keypointsY[iy+1], vb.Width-1).Styled(style).Fill(border.TLB)
	}

	// 2 way conjuntions (angles): TL, TR, BL, BR
	/// top left
	vb.Pixel(0, 0).Styled(style).Fill(border.TopLeft)
	/// top right
	vb.Pixel(0, vb.Width-1).Styled(style).Fill(border.TopRight)
	/// bottom left
	vb.Pixel(vb.Height-1, 0).Styled(style).Fill(border.BottomLeft)
	/// bottom right
	vb.Pixel(vb.Height-1, vb.Width-1).Styled(style).Fill(border.BottomRight)
}
