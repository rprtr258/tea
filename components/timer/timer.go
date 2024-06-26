// Package timer provides a simple timeout component.
package timer

import (
	"time"

	"github.com/rprtr258/tea"
)

// Authors note with regard to start and stop commands:
//
// Technically speaking, sending commands to start and stop the timer in this
// case is extraneous. To stop the timer we'd just need to set the 'running'
// property on the model to false which cause logic in the update function to
// stop responding to TickMsgs. To start the model we'd set 'running' to true
// and fire off a MsgTick. Helper functions would look like:
//
//     func (m *model) Start() tea.Cmd
//     func (m *model) Stop()
//
// The danger with this approach, however, is that order of operations becomes
// important with helper functions like the above. Consider the following:
//
//     // Would not work
//     return m, m.timer.Start()
//
//	   // Would work
//     cmd := m.timer.start()
//     return m, cmd
//
// Thus, because of potential pitfalls like the ones above, we've introduced
// the extraneous MsgStartStop to simplify the mental model when using this
// package. Bear in mind that the practice of sending commands to simply
// communicate with other parts of your application, such as in this package,
// is still not recommended.

// MsgStartStop is used to start and stop the timer.
type MsgStartStop struct {
	ID      int
	running bool
}

// MsgTick is a message that is sent on every timer tick.
type MsgTick struct {
	// ID is the identifier of the timer that sends the message. This makes
	// it possible to determine which timer a tick belongs to when there
	// are multiple timers running.
	//
	// Note, however, that a timer will reject ticks from other timers, so
	// it's safe to flow all MsgTick-s through all timers and have them still
	// behave appropriately.
	ID int

	// Timeout returns whether or not this tick is a timeout tick. You can
	// alternatively listen for MsgTimeout.
	Timeout bool
}

// MsgTimeout is a message that is sent once when the timer times out.
//
// It's a convenience message sent alongside a MsgTick with the Timeout value set to true.
type MsgTimeout struct {
	ID int
}

var _id = 0

// Model of the timer component.
type Model struct {
	id int

	// How long until the timer expires.
	Timeout time.Duration

	// How long to wait before every tick. Defaults to 1 second.
	Interval time.Duration

	running bool
}

// NewWithInterval creates a new timer with the given timeout and tick interval.
func NewWithInterval(timeout, interval time.Duration) Model {
	_id++
	return Model{
		id:       _id,
		Timeout:  timeout,
		Interval: interval,
		running:  true,
	}
}

// New creates a new timer with the given timeout and default 1s interval.
func New(timeout time.Duration) Model {
	return NewWithInterval(timeout, time.Second)
}

// Running returns whether or not the timer is running. If the timer has timed
// out this will always return false.
func (m *Model) Running() bool {
	if m.Timedout() || !m.running {
		return false
	}
	return true
}

// Timedout returns whether or not the timer has timed out.
func (m *Model) Timedout() bool {
	return m.Timeout <= 0
}

// Init starts the timer.
func (m *Model) Init(f func(...tea.Cmd)) {
	f(m.tick())
}

// Update handles the timer tick.
func (m *Model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case MsgStartStop:
		if msg.ID != 0 && msg.ID != m.id {
			return
		}
		m.running = msg.running
		f(m.tick())
	case MsgTick:
		if !m.Running() || msg.ID != 0 && msg.ID != m.id {
			break
		}

		m.Timeout -= m.Interval
		f(m.timedout()...)
		f(m.tick())
	}
}

// View of the timer component.
func (m *Model) View(vb tea.Viewbox) {
	vb.WriteLine(m.Timeout.String())
}

// CmdStart resumes the timer. Has no effect if the timer has timed out.
func (m *Model) CmdStart() tea.Cmd {
	return m.cmdStartStop(true)
}

// CmdStop pauses the timer. Has no effect if the timer has timed out.
func (m *Model) CmdStop() tea.Cmd {
	return m.cmdStartStop(false)
}

// CmdToggle stops the timer if it's running and starts it if it's stopped.
func (m *Model) CmdToggle() tea.Cmd {
	return m.cmdStartStop(!m.Running())
}

func (m *Model) tick() tea.Cmd {
	id := m.id
	return tea.Tick(m.Interval, func(_ time.Time) tea.Msg {
		return MsgTick{ID: id, Timeout: m.Timedout()}
	})
}

func (m *Model) timedout() []tea.Cmd {
	if !m.Timedout() {
		return nil
	}

	return []tea.Cmd{func() tea.Msg {
		return MsgTimeout{ID: m.id}
	}}
}

func (m *Model) cmdStartStop(v bool) tea.Cmd {
	return func() tea.Msg {
		return MsgStartStop{ID: m.id, running: v}
	}
}
