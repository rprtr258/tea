package plasma

import (
	"context"
	"math"
	"time"

	"github.com/rprtr258/scuf"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

const _fps = 60

type msgFrame struct{}

var animate = tea.Tick(time.Second/_fps, func(_ time.Time) tea.Msg {
	return msgFrame{}
})

type model struct {
	cells         []scuf.Modifier
	height, width int
}

func (m *model) Init(yield func(...tea.Cmd)) {
	yield(animate)
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		yield(tea.Quit)
	case tea.MsgWindowSize:
		m.cells = make([]scuf.Modifier, msg.Width*msg.Height)
		m.height = msg.Height
		m.width = msg.Width
	case msgFrame:
		t := float64(time.Now().UnixMilli()) / 1000
		for y := 0; y < m.height; y++ {
			for x := 0; x < m.width; x++ {
				uvx := float64(x) / float64(m.width)
				uvy := float64(y) / float64(m.height)

				v1 := math.Sin(uvx*5 + t)
				v2 := math.Sin(5*(uvx*math.Sin(t/12)+uvy*math.Cos(t/13)) + t)

				cx := uvx + math.Sin(t/15)*5
				cy := uvy + math.Sin(t/13)*5
				v3 := math.Sin(math.Sqrt(100*(cx*cx+cy*cy)) + t)

				vf := v1 + v2 + v3
				r := uint8(max(0, math.Cos(vf*math.Pi+0*math.Pi/1)-0.5) * 2 * 255)
				g := uint8(max(0, math.Sin(vf*math.Pi+6*math.Pi/3)-0.5) * 2 * 255)
				b := uint8(max(0, math.Sin(vf*math.Pi+4*math.Pi/3)-0.5) * 2 * 255)

				i := y*m.width + x
				m.cells[i] = scuf.BgRGB(r, g, b)
			}
		}
		yield(animate)
	}
}

func (m *model) View(vb tea.Viewbox) {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			vb.Styled(styles.Style{}.Background(styles.Raw(m.cells[y*m.width+x]))).Set(y, x, ' ')
		}
	}
}

func Main(ctx context.Context) error {
	_, err := tea.
		NewProgram(ctx, &model{}).
		WithAltScreen().
		Run()
	return err
}
