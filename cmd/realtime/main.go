package realtime

// A simple example that shows how to send activity to Tea in real-time
// through a channel.

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type msgResponse struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Tea will run this
// asynchronously.
func cmdListenForActivity(c tea.Context[*model], sub chan struct{}) {
	c.F(func() tea.Msg2[*model] {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	})
}

// A command that waits for the activity on a channel.
func cmdWaitForActivity(c tea.Context[*model], sub chan struct{}) {
	c.F(func() tea.Msg2[*model] {
		return func(m *model) {
			msg := msgResponse(<-sub)
			m.Update(c, msg)
		}
	})
}

type model struct {
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
}

func (m *model) Init(c tea.Context[*model]) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	m.spinner.CmdTick(ctxSpinner)
	cmdWaitForActivity(c, m.sub)   // wait for activity
	cmdListenForActivity(c, m.sub) // generate activity
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	switch msg.(type) {
	case tea.MsgKey:
		c.Dispatch(tea.Quit)
	case msgResponse:
		m.responses++                // record external activity
		cmdWaitForActivity(c, m.sub) // wait for next event
	case spinner.MsgTick:
		ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
		m.spinner.Update(ctxSpinner, msg)
	}
}

func (m *model) View(vb tea.Viewbox) {
	m.spinner.View(vb)
	vb.PaddingLeft(4).PaddingTop(1).WriteLine(" " + " Events received: " + fmt.Sprint(m.responses))
	vb.PaddingTop(3).WriteLine(" Press any key to exit")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram2(ctx, &model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	}).Run()
	return err
}
