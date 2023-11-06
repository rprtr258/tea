package shader

import (
	"time"

	"github.com/rprtr258/scuf"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

type Model struct {
	Shader func(y, x, h, w int, t time.Time) (rune, scuf.Modifier)
}

func (Model) Init(func(...tea.Cmd)) {}

func (Model) Update(tea.Msg, func(...tea.Cmd)) {}

func (m Model) View(vb tea.Viewbox) {
	// TODO: move to tea
	t := time.Now()

	for y := 0; y < vb.Height; y++ {
		for x := 0; x < vb.Width; x++ {
			r, color := m.Shader(y, x, vb.Height, vb.Width, t)
			vb.Styled(styles.Style{}.Background(color)).Set(y, x, r)
		}
	}
}
