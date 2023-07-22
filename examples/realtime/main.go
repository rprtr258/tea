package realtime

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type msgResponse struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return msgResponse(<-sub)
	}
}

type model struct {
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
	quitting  bool
}

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{
		m.spinner.Tick,
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity
	}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg.(type) {
	case tea.MsgKey:
		m.quitting = true
		return []tea.Cmd{tea.Quit}
	case msgResponse:
		m.responses++                            // record external activity
		return []tea.Cmd{waitForActivity(m.sub)} // wait for next event
	case spinner.MsgTick:
		return m.spinner.Update(msg)
	default:
		return nil
	}
}

func (m *model) View(r tea.Renderer) {
	s := fmt.Sprintf("\n %s Events received: %d\n\n Press any key to exit\n", m.spinner.View(), m.responses)
	if m.quitting {
		s += "\n"
	}
	r.Write(s)
}

func Main() {
	p := tea.NewProgram(context.Background(), &model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	})

	if _, err := p.Run(); err != nil {
		log.Fatalln("could not start program:", err.Error())
	}
}
