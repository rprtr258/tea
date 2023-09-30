package realtime

// A simple example that shows how to send activity to Bubble Tea in real-time
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
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func cmdListenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func cmdWaitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return msgResponse(<-sub)
	}
}

type model struct {
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(
		m.spinner.CmdTick,
		cmdWaitForActivity(m.sub),   // wait for activity
		cmdListenForActivity(m.sub), // generate activity
	)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
	case msgResponse:
		m.responses++                // record external activity
		f(cmdWaitForActivity(m.sub)) // wait for next event
	case spinner.MsgTick:
		f(m.spinner.Update(msg)...)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine(1, 0, " "+m.spinner.View()+" Events received: "+fmt.Sprint(m.responses))
	vb.WriteLine(3, 0, " Press any key to exit")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	}).Run()
	return err
}
