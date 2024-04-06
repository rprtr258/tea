package simple

// A simple program that counts down from 5 and then exits.

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rprtr258/tea"
)

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type msgTick time.Time

func cmdTick() tea.Msg {
	time.Sleep(time.Second)
	return msgTick{}
}

// A model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type model int

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m *model) Init(f func(...tea.Cmd)) {
	f(cmdTick)
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg.(type) {
	case tea.MsgKey:
		f(tea.Quit)
	case msgTick:
		*m--
		if *m <= 0 {
			f(tea.Quit)
			return
		}

		f(cmdTick)
	}
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine("Hi. This program will exit in " + fmt.Sprint(*m) + " seconds. To quit sooner press any key.")
}

func Main(ctx context.Context) error {
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	if logfilePath := os.Getenv("tea_LOG"); logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatalln(err.Error())
		}
	}

	// Initialize our program
	m := model(5)
	_, err := tea.NewProgram(ctx, &m).Run()
	return err
}
