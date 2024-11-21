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
	running bool
}

// MsgTick is a message that is sent on every timer tick.
type MsgTick struct {
	// Timeout returns whether or not this tick is a timeout tick. You can
	// alternatively listen for MsgTimeout.
	Timeout bool
}

// Model of the timer component.
type Model struct {
	// How long until the timer expires.
	Timeout time.Duration

	// How long to wait before every tick. Defaults to 1 second.
	Interval time.Duration

	running bool
}

type Cmd = func(*Model)

// NewWithInterval creates a new timer with the given timeout and tick interval.
func NewWithInterval(timeout, interval time.Duration) Model {
	return Model{
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
func (m *Model) Init(c tea.Context[*Model]) {
	m.tick(c)
}

// Update handles the timer tick.
func (m *Model) Update(c tea.Context[*Model], msg tea.Msg) {}

// View of the timer component.
func (m *Model) View(vb tea.Viewbox) {
	vb.WriteLine(m.Timeout.String())
}

// CmdStart resumes the timer. Has no effect if the timer has timed out.
func (m *Model) CmdStart(c tea.Context[*Model]) {
	m.cmdStartStop(c, true)
}

// CmdStop pauses the timer. Has no effect if the timer has timed out.
func (m *Model) CmdStop(c tea.Context[*Model]) {
	m.cmdStartStop(c, false)
}

// CmdToggle stops the timer if it's running and starts it if it's stopped.
func (m *Model) CmdToggle(c tea.Context[*Model]) {
	m.cmdStartStop(c, !m.Running())
}

func (m *Model) tick(c tea.Context[*Model]) {
	// msg := MsgTick{Timeout: m.Timedout()}
	// TODO: use tea.Tick(m.Interval)
	c.F(func() tea.Msg2[*Model] {
		return func(m *Model) {
			<-time.After(m.Interval)
			if !m.Running() {
				return
			}

			m.Timeout -= m.Interval
			m.timedout(c)
			m.tick(c)
		}
	})
}

type MsgTimeout struct{}

func (m *Model) timedout(c tea.Context[*Model]) {
	if !m.Timedout() {
		return
	}

	// MsgTimeout must be sent, but no reaction on it
	// MsgTimeout is a message that is sent once when the timer times out.
	//
	// It's a convenience message sent alongside a MsgTick with the Timeout value set to true.
	c.F(func() tea.Msg2[*Model] {
		// msg := MsgTimeout{}
		return func(m *Model) {}
	})
}

func (m *Model) cmdStartStop(c tea.Context[*Model], v bool) {
	msg := MsgStartStop{running: v}
	c.F(func() tea.Msg2[*Model] {
		return func(m *Model) {
			m.running = msg.running
			m.tick(c)
		}
	})
}
