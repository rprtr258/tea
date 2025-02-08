package tea

import (
	"bytes"
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/rprtr258/fun"
	"github.com/rprtr258/tea/styles"
)

// func ctrlSeq(code string) string {
// 	return termenv.CSI + code + "m"
// }

type framebuffer struct {
	Height, Width int
	B             []rune
	// OPTIMIZE: store ranges of colors instead of color for every pixel
	// backgrounds, foregrounds []string
	styles []styles.Style
}

// Viewbox is a view of the terminal to render to
type Viewbox struct {
	Height, Width int
	Y, X          int

	fb    framebuffer
	style styles.Style
}

// NewViewbox creates a new Framebuffer
func NewViewbox(height, width int) Viewbox {
	buf := make([]rune, height*width)
	for i := range buf {
		buf[i] = ' '
	}

	styless := make([]styles.Style, height*width)
	for i := range styless {
		styless[i] = styles.Style{}
	}

	return Viewbox{
		fb: framebuffer{
			Height: height,
			Width:  width,
			B:      buf,
			// backgrounds: make([]string, height*width),
			// foregrounds: make([]string, height*width),
			styles: styless,
		},
		Height: height,
		Width:  width,
		Y:      0,
		X:      0,
		style:  styles.Style{},
	}
}

func (vb Viewbox) String() string {
	return fmt.Sprintf(
		"Viewbox{Height: %d, Width: %d, Y: %d, X: %d}",
		vb.Height, vb.Width, vb.Y, vb.X,
	)
}

func (vb *Viewbox) clear() {
	for i := range vb.fb.B {
		vb.fb.B[i] = ' '
	}

	for i := range vb.fb.styles {
		vb.fb.styles[i] = styles.Style{}
	}

	vb.style = styles.Style{}
}

var sb = func() bytes.Buffer {
	var sb bytes.Buffer
	sb.Grow(32 * 105)
	return sb
}()

// Render framebuffer to string
// TODO: optimize
func (vb Viewbox) Render() []byte {
	sb.Reset()
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
		Width:  vb.Width,
		Y:      y + vb.Y,
		X:      vb.X,
		style:  vb.style,
	}
}

// PaddingOptions is padding options
type PaddingOptions struct {
	Top, Left     int
	Bottom, Right int
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

type Rectangle struct {
	Top, Left     int
	Height, Width int
}

func (vb Viewbox) Sub(rect Rectangle) Viewbox {
	return Viewbox{
		fb:     vb.fb,
		Height: fun.Clamp(fun.IF(rect.Height == 0, vb.Height, rect.Height), 0, vb.Height),
		Width:  fun.Clamp(fun.IF(rect.Width == 0, vb.Width, rect.Width), 0, vb.Width),
		Y:      vb.Y + fun.Clamp(rect.Top, 0, vb.Height-1),
		X:      vb.X + fun.Clamp(rect.Left, 0, vb.Width-1),
		style:  vb.style,
	}
}

func (vb Viewbox) Pixel(y, x int) Viewbox {
	return vb.Sub(Rectangle{
		Top:    y,
		Left:   x,
		Height: 1,
		Width:  1,
	})
}

func (vb Viewbox) PaddingTop(top int) Viewbox {
	return vb.Sub(Rectangle{
		Top:    top,
		Height: vb.Height - top,
		Width:  vb.Width,
	})
}

func (vb Viewbox) PaddingLeft(left int) Viewbox {
	return vb.Sub(Rectangle{
		Left:   left,
		Height: vb.Height,
		Width:  vb.Width - left,
	})
}

func (vb Viewbox) MaxHeight(height int) Viewbox {
	return vb.Sub(Rectangle{
		Height: min(vb.Height, height),
		Width:  vb.Width,
	})
}

func (vb Viewbox) MaxWidth(width int) Viewbox {
	return vb.Sub(Rectangle{
		Height: vb.Height,
		Width:  min(vb.Width, width),
	})
}

func (vb Viewbox) Styled(style styles.Style) Viewbox {
	vb = Viewbox{
		fb:     vb.fb,
		Height: vb.Height,
		Width:  vb.Width,
		Y:      vb.Y,
		X:      vb.X,
		style:  style.Inherit(vb.style),
	}
	if bg := style.GetBackground(); bg != nil {
		for y := 0; y < vb.Height; y++ {
			for x := 0; x < vb.Width; x++ {
				i := (vb.Y+y)*vb.fb.Width + vb.X + x
				vb.fb.styles[i] = vb.fb.styles[i].Background(bg)
			}
		}
	}
	return vb
}

// WriteLine starting from x=offset without wrapping, returns end offset
func (vb Viewbox) WriteLine(line string) int {
	x := 0
	for _, c := range line {
		if x >= vb.Width {
			return x
		}

		x += vb.Set(0, x, c)
	}
	return x
}

func (vb Viewbox) WriteLineX(line string) Viewbox {
	return vb.PaddingLeft(vb.WriteLine(line))
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
func (vb Viewbox) Set(y, x int, c rune) int {
	if y < 0 || y >= vb.Height || x < 0 || x >= vb.Width || len(vb.fb.B) == 0 {
		return 0
	}

	if c == '\n' {
		panic(fmt.Sprintf("unexpected newline: %d, %d", y, x))
	}

	i := (vb.Y+y)*vb.fb.Width + vb.X + x
	vb.fb.B[i] = c
	vb.fb.styles[i] = vb.style

	if runewidth.RuneWidth(c) == 2 {
		i++
		vb.fb.B[i] = 0
		vb.fb.styles[i] = styles.Style{}
		return 2
	}

	return 1
}

func (vb Viewbox) Fill(c rune) {
	for y := 0; y < vb.Height; y++ {
		for x := 0; x < vb.Width; x++ {
			vb.Set(y, x, c)
		}
	}
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

type Layout int

func Auto() Layout { // TODO: implement
	return Flex(1)
}

func Fixed(n int) Layout {
	return Layout(-n)
}

func Flex(n int) Layout {
	return Layout(n)
}

func EvalLayout(x int, ls ...Layout) []int {
	res := make([]int, len(ls))
	totalFlex := 0
	flexes := 0
	for i, l := range ls {
		if l < 0 {
			x += int(l)
			res[i] = -int(l)
		} else {
			totalFlex += int(l)
			flexes++
		}
	}
	for i, l := range ls {
		if l < 0 {
			continue
		}

		if flexes == 1 {
			res[i] = x
			break
		}

		res[i] = int(float64(x*int(l)) / float64(totalFlex))
		x -= res[i]
		flexes--
	}
	return res
}

func (vb Viewbox) SplitY(ls ...Layout) []Viewbox {
	heights := EvalLayout(vb.Height, ls...)
	res := make([]Viewbox, len(heights))
	y := vb.Y
	for i, h := range heights {
		res[i] = Viewbox{
			fb:     vb.fb,
			Height: h,
			Width:  vb.Width,
			Y:      y,
			X:      vb.X,
			style:  vb.style,
		}
		y += h
	}
	return res
}

func (vb Viewbox) SplitY2(l1, l2 Layout) (_, _ Viewbox) {
	r := vb.SplitY(l1, l2)
	return r[0], r[1]
}

func (vb Viewbox) SplitY3(l1, l2, l3 Layout) (_, _, _ Viewbox) {
	r := vb.SplitY(l1, l2, l3)
	return r[0], r[1], r[2]
}

func (vb Viewbox) SplitY4(l1, l2, l3, l4 Layout) (_, _, _, _ Viewbox) {
	r := vb.SplitY(l1, l2, l3, l4)
	return r[0], r[1], r[2], r[3]
}

func (vb Viewbox) SplitY5(l1, l2, l3, l4, l5 Layout) (_, _, _, _, _ Viewbox) {
	r := vb.SplitY(l1, l2, l3, l4, l5)
	return r[0], r[1], r[2], r[3], r[4]
}

func (vb Viewbox) SplitX(ls ...Layout) []Viewbox {
	widths := EvalLayout(vb.Width, ls...)
	res := make([]Viewbox, len(widths))
	x := vb.X
	for i, w := range widths {
		res[i] = Viewbox{
			fb:     vb.fb,
			Height: vb.Height,
			Width:  w,
			Y:      vb.Y,
			X:      x,
			style:  vb.style,
		}
		x += w
	}
	return res
}

func (vb Viewbox) SplitX2(l1, l2 Layout) (_, _ Viewbox) {
	r := vb.SplitX(l1, l2)
	return r[0], r[1]
}

func (vb Viewbox) SplitX3(l1, l2, l3 Layout) (_, _, _ Viewbox) {
	r := vb.SplitX(l1, l2, l3)
	return r[0], r[1], r[2]
}
