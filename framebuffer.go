package tea

import (
	"bytes"
	"fmt"

	"github.com/rprtr258/tea/lipgloss"
)

// func ctrlSeq(code string) string {
// 	return termenv.CSI + code + "m"
// }

type framebuffer struct {
	Height, Width int
	B             []rune
	// OPTIMIZE: store ranges of colors instead of color for every pixel
	// backgrounds, foregrounds []string
	styles []lipgloss.Style
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

	styles := make([]lipgloss.Style, height*width)
	for i := range styles {
		styles[i] = lipgloss.NewStyle()
	}

	return Viewbox{
		fb: framebuffer{
			Height: height,
			Width:  width,
			B:      buf,
			// backgrounds: make([]string, height*width),
			// foregrounds: make([]string, height*width),
			styles: styles,
		},
		Height: height,
		Width:  width,
		Y:      0,
		X:      0,
		style:  lipgloss.NewStyle(),
	}
}

func (vb *Viewbox) Clear() {
	for i := range vb.fb.B {
		vb.fb.B[i] = ' '
	}

	for i := range vb.fb.styles {
		vb.fb.styles[i] = lipgloss.NewStyle()
	}

	vb.style = lipgloss.NewStyle()
}

// Render framebuffer to string
// TODO: optimize
func (vb Viewbox) Render() []byte {
	var sb bytes.Buffer
	// bg := ""
	// fg := ""
	for y := 0; y < vb.fb.Height*vb.fb.Width; y += vb.fb.Width {
		if y > 0 {
			sb.WriteRune('\n')
		}

		// fullRow := vb.fb.B[y : y+vb.fb.Width]
		for x := 0; x < vb.fb.Width; x++ {
			i := y + x

			sb.WriteString(vb.fb.styles[i].Render(string([]rune{vb.fb.B[i]})))

			// 	if vb.fb.backgrounds[y+x] != bg || vb.fb.foregrounds[y+x] != fg {
			// 		bg = vb.fb.backgrounds[y+x]
			// 		fg = vb.fb.foregrounds[y+x]
			// 		sb.WriteString(ctrlSeq(termenv.ResetSeq) + lo.
			// 			Switch[bool, string](true).
			// 			Case(bg == "" && fg == "", "").
			// 			Case(bg == "", ctrlSeq(fg)).
			// 			Case(fg == "", ctrlSeq(bg)).
			// 			Default(ctrlSeq(bg+";"+fg)))
			// 	}
			// 	sb.WriteRune(fullRow[x])
		}
	}
	return sb.Bytes()
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

func (vb Viewbox) PaddingTop(top int) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: vb.Height - top,
		Width:  vb.Width,
		Y:      vb.Y + top,
		X:      vb.X,
		style:  vb.style,
	}
}

func (vb Viewbox) PaddingLeft(left int) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: vb.Height,
		Width:  vb.Width - left,
		Y:      vb.Y,
		X:      vb.X + left,
		style:  vb.style,
	}
}

func (vb Viewbox) MaxHeight(height int) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: min(vb.Height, height),
		Width:  vb.Width,
		Y:      vb.Y,
		X:      vb.X,
		style:  vb.style,
	}
}

func (vb Viewbox) MaxWidth(width int) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: vb.Height,
		Width:  min(vb.Width, width),
		Y:      vb.Y,
		X:      vb.X,
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

// WriteLine starting from x=offset without wrapping, returns end offset
func (vb Viewbox) WriteLine(y, offsetX int, line string) int {
	x := offsetX
	for _, c := range line {
		if x >= vb.Width {
			return x
		}

		vb.Set(y, x, c)
		x++
	}
	return x
}

// WriteText starting from y, x with wrapping, returns end position
func (vb Viewbox) WriteText(y, x int, text string) (int, int) {
	for _, c := range text {
		if c == '\n' { // TODO: remove
			x = 0
			y++
			continue
		}

		vb.Set(y, x, c)
		x++
		if x >= vb.Width {
			x = 0
			y++
		}
	}
	return y, x
}

// Set writes a rune to the framebuffer in position relative to viewbox
// 0 <= y < height, 0 <= x < width
func (vb Viewbox) Set(y, x int, c rune) {
	if y < 0 || y >= vb.Height || x < 0 || x >= vb.Width {
		return
	}

	if c == '\n' {
		panic(fmt.Sprintf("unexpected newline: %d, %d", y, x))
	}

	i := (vb.Y+y)*vb.fb.Width + vb.X + x
	vb.fb.B[i] = c
	vb.fb.styles[i] = vb.style
}

// // background colors y'th row bacground to given color from x1 to x2 with
// // coordinates relative to viewbox
// func (vb Viewbox) background(y, x1, x2 int, background termenv.Color) {
// 	for x := x1 + vb.X; x < x2+vb.X; x++ {
// 		vb.fb.backgrounds[(y+vb.Y)*vb.fb.Width+x] = background.Sequence(true)
// 	}
// }

// // Background colors y'th row foreground to given color from x1 to x2 with
// // coordinates relative to viewbox
// func (vb Viewbox) foreground(y, x1, x2 int, foreground termenv.Color) {
// 	for x := x1 + vb.X; x < x2+vb.X; x++ {
// 		vb.fb.foregrounds[(y+vb.Y)*vb.fb.Width+x] = foreground.Sequence(false)
// 	}
// }
