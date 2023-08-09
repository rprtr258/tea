package cellbuffer

// A simple example demonstrating how to draw and animate on a cellular grid.
// Note that the cellbuffer implementation in this example does not support
// double-width runes.

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/rprtr258/tea"
)

const (
	fps       = 60
	frequency = 7.5
	damping   = 0.15
	asterisk  = "*"
)

func drawEllipse(cb *cellbuffer, xc, yc, rx, ry float64) {
	var (
		dx, dy, d1, d2 float64
		x              float64
		y              = ry
	)

	d1 = ry*ry - rx*rx*ry + 0.25*rx*rx
	dx = 2 * ry * ry * x
	dy = 2 * rx * rx * y

	for dx < dy {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d1 < 0 {
			x++
			dx = dx + (2 * ry * ry)
			d1 = d1 + dx + (ry * ry)
		} else {
			x++
			y--
			dx = dx + (2 * ry * ry)
			dy = dy - (2 * rx * rx)
			d1 = d1 + dx - dy + (ry * ry)
		}
	}

	d2 = ((ry * ry) * ((x + 0.5) * (x + 0.5))) + ((rx * rx) * ((y - 1) * (y - 1))) - (rx * rx * ry * ry)

	for y >= 0 {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d2 > 0 {
			y--
			dy = dy - (2 * rx * rx)
			d2 = d2 + (rx * rx) - dy
		} else {
			y--
			x++
			dx = dx + (2 * ry * ry)
			dy = dy - (2 * rx * rx)
			d2 = d2 + dx - dy + (rx * rx)
		}
	}
}

type cellbuffer struct {
	cells  []string
	stride int
}

func (c *cellbuffer) init(w, h int) {
	if w == 0 {
		return
	}
	c.stride = w
	c.cells = make([]string, w*h)
	c.wipe()
}

func (c cellbuffer) set(x, y int) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = asterisk
}

func (c *cellbuffer) wipe() {
	for i := range c.cells {
		c.cells[i] = " "
	}
}

func (c cellbuffer) width() int {
	return c.stride
}

func (c cellbuffer) height() int {
	h := len(c.cells) / c.stride
	if len(c.cells)%c.stride != 0 {
		h++
	}
	return h
}

func (c cellbuffer) ready() bool {
	return len(c.cells) > 0
}

func (c cellbuffer) String() string {
	var sb strings.Builder
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			sb.WriteRune('\n')
		}
		sb.WriteString(c.cells[i])
	}
	return sb.String()
}

type msgFrame struct{}

var animate = tea.Tick(time.Second/fps, func(_ time.Time) tea.Msg {
	return msgFrame{}
})

type model struct {
	cells                cellbuffer
	spring               harmonica.Spring
	targetX, targetY     float64
	x, y                 float64
	xVelocity, yVelocity float64
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(animate)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
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
		f(animate)
	}
}

func (m *model) View(r tea.Renderer) {
	r.Write(m.cells.String())
}

func Main(ctx context.Context) error {
	_, err := tea.
		NewProgram(ctx, &model{
			spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
		}).
		WithAltScreen().
		WithMouseCellMotion().
		Run()
	return err
}
