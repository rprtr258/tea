package tea

import (
	"strings"

	"github.com/muesli/termenv"
	"github.com/samber/lo"
)

func ctrlSeq(code string) string {
	return termenv.CSI + code + "m"
}

type framebuffer struct {
	B []rune
	// OPTIMIZE: store ranges of colors instead of color for every pixel
	backgrounds, foregrounds []string
	Height, Width            int
}

// Viewbox is a view of the terminal to render to
type Viewbox struct {
	framebuffer
	ViewHeight, ViewWidth int
	Y, X                  int
}

// NewViewbox creates a new Framebuffer
func NewViewbox(height, width int) Viewbox {
	buf := make([]rune, height*width)
	for i := range buf {
		buf[i] = ' '
	}
	return Viewbox{
		framebuffer: framebuffer{
			B:           buf,
			backgrounds: make([]string, height*width),
			foregrounds: make([]string, height*width),
			Height:      height,
			Width:       width,
		},
		ViewHeight: height,
		ViewWidth:  width,
		Y:          0,
		X:          0,
	}
}

// Render framebuffer to string
func (vb Viewbox) Render() string {
	var sb strings.Builder
	bg := ""
	fg := ""
	for y := 0; y < vb.Height*vb.Width; y += vb.Width {
		if y > 0 {
			sb.WriteRune('\n')
		}

		fullRow := vb.B[y : y+vb.Width]
		for x := 0; x < vb.Width; x++ {
			if vb.backgrounds[y+x] != bg || vb.foregrounds[y+x] != fg {
				bg = vb.backgrounds[y+x]
				fg = vb.foregrounds[y+x]
				sb.WriteString(ctrlSeq(termenv.ResetSeq) + lo.
					Switch[bool, string](true).
					Case(bg == "" && fg == "", "").
					Case(bg == "", ctrlSeq(fg)).
					Case(fg == "", ctrlSeq(bg)).
					Default(ctrlSeq(bg+";"+fg)))
			}
			sb.WriteRune(fullRow[x])
		}
	}
	return sb.String()
}

// Row returns view to current viewbox's row
func (vb Viewbox) Row(y int) Viewbox {
	return Viewbox{
		framebuffer: vb.framebuffer,
		ViewHeight:  1,
		ViewWidth:   vb.Width,
		Y:           y + vb.Y,
		X:           vb.X,
	}
}

// PaddingOptions is padding options
type PaddingOptions struct {
	Top, Bottom, Left, Right int
}

// Padding returns view to current viewbox inner with given paddings and size
// 0 <= top <= bottom < height, 0 <= left <= right < width
func (vb Viewbox) Padding(opt PaddingOptions) Viewbox {
	return Viewbox{
		framebuffer: vb.framebuffer,
		ViewHeight:  vb.ViewHeight - opt.Top - opt.Bottom,
		ViewWidth:   vb.ViewWidth - opt.Left - opt.Right,
		Y:           vb.Y + opt.Top,
		X:           vb.X + opt.Left,
	}
}

// Set writes a rune to the framebuffer in position relative to viewbox
// 0 <= y < height, 0 <= x < width
func (vb Viewbox) Set(y, x int, c rune) {
	vb.B[(vb.Y+y)*vb.Width+vb.X+x] = c
}

// Background colors y'th row bacground to given color from x1 to x2 with
// coordinates relative to viewbox
func (vb Viewbox) Background(y, x1, x2 int, background termenv.Color) {
	for x := x1 + vb.X; x < x2+vb.X; x++ {
		vb.backgrounds[(y+vb.Y)*vb.Width+x] = background.Sequence(true)
	}
}

// Background colors y'th row foreground to given color from x1 to x2 with
// coordinates relative to viewbox
func (vb Viewbox) Foreground(y, x1, x2 int, foreground termenv.Color) {
	for x := x1 + vb.X; x < x2+vb.X; x++ {
		vb.foregrounds[(y+vb.Y)*vb.Width+x] = foreground.Sequence(false)
	}
}
