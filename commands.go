package tea

import (
	"time"
)

// Every is a command that ticks in sync with the system clock. So, if you
// wanted to tick with the system clock every second, minute or hour you
// could use this. It's also handy for having different things tick in sync.
//
// Because we're ticking with the system clock the tick will likely not run for
// the entire specified duration. For example, if we're ticking for one minute
// and the clock is at 12:34:20 then the next tick will happen at 12:35:00, 40
// seconds later.
//
// To produce the command, pass a duration and a function which returns
// a message containing the time at which the tick occurred.
//
//	type MsgTick time.Time
//
//	cmd := Every(time.Second, func(t time.Time) Msg {
//	   return MsgTick(t)
//	})
//
// Beginners' note: Every sends a single message and won't automatically
// dispatch messages at an interval. To do that, you'll want to return another
// Every command after receiving your tick message. For example:
//
//	type MsgTick time.Time
//
//	// Send a message every second.
//	func tickEvery() Cmd {
//	    return Every(time.Second, func(t time.Time) Msg {
//	        return MsgTick(t)
//	    })
//	}
//
//	func (m *model) Init() Cmd {
//	    // Start ticking.
//	    return tickEvery()
//	}
//
//	func (m *model) Update(msg Msg) Cmd {
//	    switch msg.(type) {
//	    case MsgTick:
//	        // Return your Every command again to loop.
//	        return tickEvery()
//	    }
//	    return nil
//	}
//
// Every is analogous to Tick in the Elm Architecture.
func Every(d time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		now := time.Now()
		return fn(<-time.After(d - now.Sub(now.Truncate(d))))
	}
}

// Tick produces a command at an interval independent of the system clock at
// the given duration. That is, the timer begins precisely when invoked,
// and runs for its entire duration.
//
// To produce the command, pass a duration and a function which returns
// a message containing the time at which the tick occurred.
//
//	type MsgTick time.Time
//
//	cmd := Tick(time.Second, func(t time.Time) Msg {
//	   return MsgTick(t)
//	})
//
// Beginners' note: Tick sends a single message and won't automatically
// dispatch messages at an interval. To do that, you'll want to return another
// Tick command after receiving your tick message. For example:
//
//	type MsgTick time.Time
//
//	func doTick() Cmd {
//	    return Tick(time.Second, func(t time.Time) Msg {
//	        return MsgTick(t)
//	    })
//	}
//
//	func (m *model) Init() Cmd {
//	    // Start ticking.
//	    return doTick()
//	}
//
//	func (m *model) Update(msg Msg) Cmd {
//	    switch msg.(type) {
//	    case MsgTick:
//	        // Return your Tick command again to loop.
//	        return doTick()
//	    }
//	    return nil
//	}
func Tick(d time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		return fn(<-time.After(d))
	}
}
