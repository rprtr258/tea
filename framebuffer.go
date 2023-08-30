package tea

import (
	"strings"

	"github.com/muesli/termenv"
	"github.com/samber/lo"

	"github.com/rprtr258/tea/lipgloss"
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
	Height, Width int
	Y, X          int

	fb    framebuffer
	style lipgloss.Style
}

// NewViewbox creates a new Framebuffer
func NewViewbox(height, width int) Viewbox {
	buf := make([]rune, height*width)
	for i := range buf {
		buf[i] = ' '
	}
	return Viewbox{
		fb: framebuffer{
			B:           buf,
			backgrounds: make([]string, height*width),
			foregrounds: make([]string, height*width),
			Height:      height,
			Width:       width,
		},
		Height: height,
		Width:  width,
		Y:      0,
		X:      0,
		style:  lipgloss.NewStyle(),
	}
}

// Render framebuffer to string
func (vb Viewbox) Render() string {
	var sb strings.Builder
	bg := ""
	fg := ""
	for y := 0; y < vb.fb.Height*vb.fb.Width; y += vb.fb.Width {
		if y > 0 {
			sb.WriteRune('\n')
		}

		fullRow := vb.fb.B[y : y+vb.fb.Width]
		for x := 0; x < vb.fb.Width; x++ {
			if vb.fb.backgrounds[y+x] != bg || vb.fb.foregrounds[y+x] != fg {
				bg = vb.fb.backgrounds[y+x]
				fg = vb.fb.foregrounds[y+x]
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
		fb:     vb.fb,
		Height: 1,
		Width:  vb.fb.Width,
		Y:      y + vb.Y,
		X:      vb.X,
		style:  vb.style,
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
		fb:     vb.fb,
		Height: vb.Height - opt.Top - opt.Bottom,
		Width:  vb.Width - opt.Left - opt.Right,
		Y:      vb.Y + opt.Top,
		X:      vb.X + opt.Left,
		style:  vb.style,
	}
}

func (vb Viewbox) Styled(style lipgloss.Style) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: vb.Height,
		Width:  vb.Width,
		Y:      vb.Y,
		X:      vb.X,
		style:  style,
	}
}

// Set writes a rune to the framebuffer in position relative to viewbox
// 0 <= y < height, 0 <= x < width
func (vb Viewbox) Set(y, x int, c rune) {
	vb.fb.B[(vb.Y+y)*vb.fb.Width+vb.X+x] = c
}

// background colors y'th row bacground to given color from x1 to x2 with
// coordinates relative to viewbox
func (vb Viewbox) background(y, x1, x2 int, background termenv.Color) {
	for x := x1 + vb.X; x < x2+vb.X; x++ {
		vb.fb.backgrounds[(y+vb.Y)*vb.fb.Width+x] = background.Sequence(true)
	}
}

// Background colors y'th row foreground to given color from x1 to x2 with
// coordinates relative to viewbox
func (vb Viewbox) foreground(y, x1, x2 int, foreground termenv.Color) {
	for x := x1 + vb.X; x < x2+vb.X; x++ {
		vb.fb.foregrounds[(y+vb.Y)*vb.fb.Width+x] = foreground.Sequence(false)
	}
}
