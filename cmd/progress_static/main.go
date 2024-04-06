package progress_static //nolint:revive,stylecheck

// A simple example that shows how to render a progress bar in a "pure"
// fashion. In this example we bump the progress by 25% every second,
// maintaining the progress state on our top level model using the progress bar
// model's ViewAs method only for rendering.
//
// The signature for ViewAs is:
//
//     func (m *model) ViewAs(percent float64) string
//
// So it takes a float between 0 and 1, and renders the progress bar
// accordingly. When using the progress bar in this "pure" fashion and there's
// no need to call an Update method.
//
// The progress bar is also able to animate itself, however. For details see
// the progress-animated example.

import (
	"context"
	"time"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/progress"
	"github.com/rprtr258/tea/styles"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = styles.Style{}.Foreground(styles.FgColor("#626262")).Render

type msgTick time.Time

var cmdTick = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return msgTick(t)
})

type model struct {
	percent  float64
	progress progress.Model
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(cmdTick)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
	case msgTick:
		m.percent += 0.25
		if m.percent > 1.0 {
			m.percent = 1.0
			f(tea.Quit)
			return
		}
		f(cmdTick)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb = vb.PaddingTop(1).PaddingLeft(padding)
	m.progress.ViewAs(vb.Row(0), m.percent)
	vb.PaddingTop(2).WriteLine(helpStyle("Press any key to quit"))
}

func Main(ctx context.Context) error {
	colorA, _ := colorful.Hex("#FF7CCB")
	colorB, _ := colorful.Hex("#FDFF8C")
	prog := progress.New(progress.WithScaledGradient(colorA, colorB))

	_, err := tea.NewProgram(ctx, &model{progress: prog}).Run()
	return err
}
