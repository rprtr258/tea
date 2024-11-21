package progress_animated //nolint:revive,stylecheck

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"context"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/progress"
	"github.com/rprtr258/tea/styles"
)

const (
	_padding  = 2
	_maxWidth = 80
)

type msgTick time.Time

type model struct {
	progress progress.Model
}

func (m *model) Init(c tea.Context[*model]) {
	m._cmdTick(c)
}

func (m *model) _cmdTick(c tea.Context[*model]) {
	// TODO: tea.Tick(time.Second*1)
	c.F(func() tea.Msg2[*model] {
		return func(m *model) {
			t := <-time.After(time.Second * 1)
			msg := msgTick(t)
			m.Update(c, msg)
		}
	})
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	ctxProgress := tea.Of(c, func(m *model) *progress.Model { return &m.progress })
	switch msg := msg.(type) {
	case tea.MsgKey:
		c.Dispatch(tea.Quit)
	case tea.MsgWindowSize:
		m.progress.Width = msg.Width - _padding*2 - 4
		if m.progress.Width > _maxWidth {
			m.progress.Width = _maxWidth
		}
	case msgTick:
		if m.progress.Percent() == 1.0 {
			c.Dispatch(tea.Quit)
			return
		}

		// Note that you can also use progress.Model.SetPercent to set the percentage value explicitly, too.
		m._cmdTick(c)
		m.progress.IncrPercent(ctxProgress, 0.25)
	// MsgFrame is sent when the progress bar wants to animate itself
	case progress.MsgFrame:
		m.progress.Update(ctxProgress, msg)
	}
}

var helpStyle = styles.Style{}.Foreground(styles.FgColor("#626262"))

func (m *model) View(vb tea.Viewbox) {
	vb = vb.PaddingTop(1).PaddingLeft(_padding)
	m.progress.View(vb)
	vb.PaddingTop(2).Styled(helpStyle).WriteLine("Press any key to quit")
}

func Main(ctx context.Context) error {
	m := &model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	_, err := tea.NewProgram2(ctx, m).Run()
	return err
}
