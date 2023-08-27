package cellbuffer

// A simple example demonstrating how to draw and animate on a cellular grid.
// Note that the cellbuffer implementation in this example does not support
// double-width runes.

import (
	"context"
	"time"

	"github.com/charmbracelet/harmonica"

	"github.com/rprtr258/tea"
)

const (
	_fps       = 60
	_frequency = 7.5
	_damping   = 0.15
	_asterisk  = '*'
)

func sqr(x float64) float64 {
	return x * x
}

func drawEllipse(cb *cellbuffer, xc, yc, rx, ry float64) {
	x, y := 0.0, ry

	d1 := sqr(ry) - sqr(rx)*ry + 0.25*sqr(rx)
	dx := 2 * sqr(ry) * x
	dy := 2 * sqr(rx) * y

	for dx < dy {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d1 < 0 {
			x++
			dx += 2 * sqr(ry)
			d1 += dx + sqr(ry)
		} else {
			x++
			y--
			dx += 2 * sqr(ry)
			dy -= 2 * sqr(rx)
			d1 += dx - dy + sqr(ry)
		}
	}

	d2 := sqr(ry)*sqr(x+0.5) + sqr(rx)*sqr(y-1) - sqr(rx)*sqr(ry)

	for y >= 0 {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d2 > 0 {
			y--
			dy -= 2 * sqr(rx)
			d2 += sqr(rx) - dy
		} else {
			y--
			x++
			dx += 2 * sqr(ry)
			dy -= 2 * sqr(rx)
			d2 += dx - dy + sqr(rx)
		}
	}
}

type cellbuffer struct {
	cells  []byte
	stride int
}

func (c *cellbuffer) init(w, h int) {
	if w == 0 {
		return
	}

	c.stride = w
	c.cells = make([]byte, w*h)
	c.wipe()
}

func (c cellbuffer) set(x, y int) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = _asterisk
}

func (c *cellbuffer) wipe() {
	for i := range c.cells {
		c.cells[i] = ' '
	}
}

func (c cellbuffer) width() int {
	return c.stride
}

func (c cellbuffer) height() int {
	return (len(c.cells) + c.stride - 1) / c.stride
}

func (c cellbuffer) ready() bool {
	return len(c.cells) > 0
}

func (c cellbuffer) Render(r tea.Renderer) {
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			r.Write("\n")
		}
		r.Write(string(c.cells[i]))
	}
}

type msgFrame struct{}

var animate = tea.Tick(time.Second/_fps, func(_ time.Time) tea.Msg {
	return msgFrame{}
})

type model struct {
	cells                cellbuffer
	spring               harmonica.Spring
	targetX, targetY     float64
	x, y                 float64
	xVelocity, yVelocity float64
}

func (m *model) Init(yield func(...tea.Cmd)) {
	yield(animate)
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		yield(tea.Quit)
	case tea.MsgWindowSize:
		if !m.cells.ready() {
			m.targetX, m.targetY = float64(msg.Width)/2, float64(msg.Height)/2
		}
		m.cells.init(msg.Width, msg.Height)
	case tea.MsgMouse:
		if !m.cells.ready() {
			return
		}

		m.targetX, m.targetY = float64(msg.X), float64(msg.Y)
	case msgFrame:
		if !m.cells.ready() {
			return
		}

		m.cells.wipe()
		m.x, m.xVelocity = m.spring.Update(m.x, m.xVelocity, m.targetX)
		m.y, m.yVelocity = m.spring.Update(m.y, m.yVelocity, m.targetY)
		drawEllipse(&m.cells, m.x, m.y, 16, 8)
		yield(animate)
	}
}

func (m *model) View(r tea.Renderer) {
	m.cells.Render(r)
}

func Main(ctx context.Context) error {
	_, err := tea.
		NewProgram(ctx, &model{
			spring: harmonica.NewSpring(harmonica.FPS(_fps), _frequency, _damping),
		}).
		WithAltScreen().
		WithMouseCellMotion().
		Run()
	return err
}
