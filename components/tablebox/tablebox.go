package tablebox

import (
	"iter"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

type node struct {
	// TODO: add padding: n - from left, -n - from right, 0 - in center
	title string

	// optional
	size     tea.Layout
	children []node
}

func SpanTitle(title string, size tea.Layout, children ...node) node {
	return node{title: title, size: size, children: children}
}

func Title(title string) node {
	return node{title: title, size: 0, children: nil}
}

func Span(size tea.Layout, children ...node) node {
	return node{title: "", size: size, children: children}
}

func boxHelper(
	vb tea.Viewbox,
	border FullBorder,
	style styles.Style,
	layout ...node, // TODO: x axis first, maybe do y axis first instead?
) iter.Seq[tea.Viewbox] {
	return func(yield func(tea.Viewbox) bool) {
		if layout == nil {
			yield(vb)
			return
		}

		ws := tea.EvalLayout(vb.Width-len(layout)+1, fun.Map[tea.Layout](func(n node) tea.Layout { return n.size }, layout...)...)
		keypointsX := make([]int, len(ws)+1) // TODO: remove keypoints array ?
		keypointsX[0] = -1
		for i, w := range ws {
			keypointsX[i+1] = keypointsX[i] + w + 1
		}
		for ix, w := range ws {
			n := layout[ix]

			if ix > 0 {
				vb.
					Sub(tea.Rectangle{
						Top:    0,
						Left:   keypointsX[ix],
						Height: vb.Height,
						Width:  1,
					}).
					Styled(style).
					Fill(border.Left)
			}

			if n.children == nil {
				if !yield(vb.Sub(tea.Rectangle{
					Top:    0,
					Left:   keypointsX[ix] + 1,
					Height: vb.Height,
					Width:  w,
				})) {
					return
				}
			} else {
				hs := tea.EvalLayout(vb.Height-len(n.children)+1, fun.Map[tea.Layout](func(n node) tea.Layout { return n.size }, n.children...)...)
				keypointsY := make([]int, len(hs)+1) // TODO: remove keypoints array ?
				keypointsY[0] = -1
				for i, h := range hs {
					keypointsY[i+1] = keypointsY[i] + h + 1
				}

				for iy, h := range hs {
					if iy > 0 {
						vb.
							Sub(tea.Rectangle{
								Top:    keypointsY[iy],
								Left:   0,
								Height: 1,
								Width:  w,
							}).
							Styled(style).
							Fill(border.Top)
					}

					boxHelper(
						vb.Sub(tea.Rectangle{
							Top:    keypointsY[iy] + 1,
							Left:   keypointsX[ix] + 1,
							Height: h,
							Width:  w,
						}),
						border, style,
						n.children[iy].children...,
					)(yield)
				}
			}
		}
	}
}

func Grid(hs, ws []tea.Layout) []node {
	rows := fun.Map[node](func(h tea.Layout) node { return Span(h) }, hs...)
	return fun.Map[node](func(w tea.Layout) node { return Span(w, rows...) }, ws...)
}

// Box draws model inside a box
func Box(
	vb tea.Viewbox,
	border FullBorder,
	style styles.Style,
	layout ...node, // TODO: x axis first, maybe do y axis first instead?
) iter.Seq2[int, tea.Viewbox] {
	// TODO: fill borders correctly and show titles
	// 2 way conjuntions (angles): TL, TR, BL, BR
	/// top left
	vb.Pixel(0, 0).Styled(style).Fill(border.TopLeft)
	/// top right
	vb.Pixel(0, vb.Width-1).Styled(style).Fill(border.TopRight)
	/// bottom left
	vb.Pixel(vb.Height-1, 0).Styled(style).Fill(border.BottomLeft)
	/// bottom right
	vb.Pixel(vb.Height-1, vb.Width-1).Styled(style).Fill(border.BottomRight)

	vb.Sub(tea.Rectangle{Top: 0, Left: 1, Height: 1, Width: vb.Width - 2}).Styled(style).Fill(border.Top)
	vb.Sub(tea.Rectangle{Top: vb.Height - 1, Left: 1, Height: 1, Width: vb.Width - 2}).Styled(style).Fill(border.Bottom)
	vb.Sub(tea.Rectangle{Top: 1, Left: 0, Height: vb.Height - 2, Width: 1}).Styled(style).Fill(border.Left)
	vb.Sub(tea.Rectangle{Top: 1, Left: vb.Width - 1, Height: vb.Height - 2, Width: 1}).Styled(style).Fill(border.Right)

	vb = vb.Padding(tea.PaddingOptions{1, 1, 1, 1})

	return func(yield func(int, tea.Viewbox) bool) {
		k := 0
		for vb := range boxHelper(vb, border, style, layout...) {
			if !yield(k, vb) {
				return
			}
			k++
		}

		if border == noBorder { // TODO: no paddings in that case
			return
		}

		// // 4 way conjunctions (full crosses): LRTB
		// for iy := 0; iy < len(hs)-1; iy++ {
		// 	for ix := 0; ix < len(ws)-1; ix++ {
		// 		vb.Pixel(keypointsY[iy+1], keypointsX[ix+1]).Styled(style).Fill(border.LRTB)
		// 	}
		// }

		// // 3 way conjunctions: LRT, LRB, TLB, TRB
		// /// LRT
		// for ix := 0; ix < len(ws)-1; ix++ {
		// 	vb.Pixel(0, keypointsX[ix+1]).Styled(style).Fill(border.LRT)
		// }
		// /// LRB
		// for ix := 0; ix < len(ws)-1; ix++ {
		// 	vb.Pixel(vb.Height-1, keypointsX[ix+1]).Styled(style).Fill(border.LRB)
		// }
		// /// TRB
		// for iy := 0; iy < len(hs)-1; iy++ {
		// 	vb.Pixel(keypointsY[iy+1], 0).Styled(style).Fill(border.TRB)
		// }
		// /// TLB
		// for iy := 0; iy < len(hs)-1; iy++ {
		// 	vb.Pixel(keypointsY[iy+1], vb.Width-1).Styled(style).Fill(border.TLB)
		// }
	}
}
