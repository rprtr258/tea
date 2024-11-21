package progress_download //nolint:revive,stylecheck

import (
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/progress"
	"github.com/rprtr258/tea/styles"
)

var helpStyle = styles.Style{}.Foreground(styles.FgColor("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type (
	msgProgress    float64
	msgProgressErr struct{ err error }
)

type model struct {
	pw       *progressWriter
	progress progress.Model
	err      error
}

func (m *model) Init(tea.Context[*model]) {}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	ctxProgress := tea.Of(c, func(m *model) *progress.Model { return &m.progress })
	switch msg := msg.(type) {
	case tea.MsgKey:
		c.Dispatch(tea.Quit)
	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
	case msgProgressErr:
		m.err = msg.err
		c.Dispatch(tea.Quit)
	case msgProgress:
		if msg >= 1.0 {
			// TODO: tea.Tick(time.Millisecond*750)
			c.F(func() tea.Msg2[*model] {
				return func(m *model) {
					<-time.After(time.Millisecond * 750)
					// TODO: ???
				}
			})
			c.Dispatch(tea.Quit)
		}

		m.progress.SetPercent(ctxProgress, float64(msg))
	// MsgFrame is sent when the progress bar wants to animate itself
	case progress.MsgFrame:
		m.progress.Update(ctxProgress, msg)
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.err != nil {
		vb.WriteLine("Error downloading: " + m.err.Error())
		return
	}

	vb = vb.PaddingTop(1).PaddingLeft(padding)
	m.progress.View(vb.Row(0))
	vb.PaddingTop(2).WriteLine(helpStyle("Press any key to quit"))
}
