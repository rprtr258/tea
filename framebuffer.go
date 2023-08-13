package tea

import (
	"strings"

	"github.com/muesli/termenv"
	"github.com/samber/lo"
)

func ctrlSeq(code string) string {
	return termenv.CSI + code + "m"
}

type screen struct {
	B []rune
	// OPTIMIZE: store ranges of colors instead of color for every pixel
	backgrounds, foregrounds []string
	Height, Width            int
}

// FrameBuffer is a view of the terminal to render to
type FrameBuffer struct {
	screen
	ViewHeight, ViewWidth int
	Y, X                  int
}

// NewFramebuffer creates a new Framebuffer
func NewFramebuffer(height, width int) FrameBuffer {
	buf := make([]rune, height*width)
	for i := range buf {
		buf[i] = ' '
	}
	return FrameBuffer{
		screen: screen{
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
func (fb FrameBuffer) Render() string {
	// OPTIMIZE: strings.Builder
	rows := make([]string, fb.Height)
	bg := ""
	fg := ""
	for y := 0; y < fb.Height; y++ {
		fullRow := fb.B[y*fb.Width : (y+1)*fb.Width]
		newRow := ""
		for x := 0; x < fb.Width; x++ {
			coloring := ""
			if fb.backgrounds[y*fb.Width+x] != bg || fb.foregrounds[y*fb.Width+x] != fg {
				bg = fb.backgrounds[y*fb.Width+x]
				fg = fb.foregrounds[y*fb.Width+x]
				coloring = ctrlSeq(termenv.ResetSeq) + lo.
					Switch[bool, string](true).
					Case(bg == "" && fg == "", "").
					Case(bg == "", ctrlSeq(fg)).
					Case(fg == "", ctrlSeq(bg)).
					Default(ctrlSeq(bg+";"+fg))
			}
			newRow += coloring + string([]rune{fullRow[x]})
		}

		rows[y] = newRow
	}
	return strings.Join(rows, "\n")
}

// Row returns view to current viewbox's row
func (fb FrameBuffer) Row(y int) FrameBuffer {
	return FrameBuffer{
		screen:     fb.screen,
		ViewHeight: 1,
		ViewWidth:  fb.Width,
		Y:          y + fb.Y,
		X:          fb.X,
	}
}

// PaddingOptions is padding options
type PaddingOptions struct {
	Top, Bottom, Left, Right int
}

// Padding returns view to current viewbox inner with given paddings and size
// 0 <= top <= bottom < height, 0 <= left <= right < width
func (fb FrameBuffer) Padding(opt PaddingOptions) FrameBuffer {
	return FrameBuffer{
		screen:     fb.screen,
		ViewHeight: fb.ViewHeight - opt.Top - opt.Bottom,
		ViewWidth:  fb.ViewWidth - opt.Left - opt.Right,
		Y:          fb.Y + opt.Top,
		X:          fb.X + opt.Left,
	}
}

// Set writes a rune to the framebuffer in position relative to viewbox
// 0 <= y < height, 0 <= x < width
func (fb FrameBuffer) Set(y, x int, c rune) {
	fb.B[(fb.Y+y)*fb.Width+fb.X+x] = c
}

// Background colors y'th row bacground to given color from x1 to x2 with
// coordinates relative to viewbox
func (fb FrameBuffer) Background(y, x1, x2 int, background termenv.Color) {
	for x := x1 + fb.X; x < x2+fb.X; x++ {
		fb.backgrounds[(y+fb.Y)*fb.Width+x] = background.Sequence(true)
	}
}

// Background colors y'th row foreground to given color from x1 to x2 with
// coordinates relative to viewbox
func (fb FrameBuffer) Foreground(y, x1, x2 int, foreground termenv.Color) {
	for x := x1 + fb.X; x < x2+fb.X; x++ {
		fb.foregrounds[(y+fb.Y)*fb.Width+x] = foreground.Sequence(false)
	}
}
