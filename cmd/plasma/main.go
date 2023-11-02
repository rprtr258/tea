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
		t := int(time.Now().UnixMilli() / 100 % 10000)
		for y := 0; y < m.height; y++ {
			for x := 0; x < m.width; x++ {
				// Calculate the pixel's coordinates in the plasma space
				xf := float64(x+t) / float64(m.width)
				yf := float64(y-t) / float64(m.height)

				// Calculate the plasma effect based on the current time
				// You can modify this formula to achieve different effects
				noise := (xf*5 + yf*7) * 4

				// Calculate the color based on the calculated noise
				red := uint8((math.Sin(noise) + 1) * 127)
				green := uint8((math.Cos(noise) + 1) * 127)
				blue := uint8((math.Sin(noise+math.Pi/2) + 1) * 127)

				i := y*m.width + x
				m.cells[i] = scuf.BgRGB(red, green, blue)
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
