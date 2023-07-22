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

func Main() {
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	logfilePath := os.Getenv("tea_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatalln(err.Error())
		}
	}

	// Initialize our program
	m := model(5)
	p := tea.NewProgram(context.Background(), &m)
	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}

// A model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type model int

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m *model) Init() tea.Cmd {
	return tick
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case tea.MsgKey:
		return tea.Quit
	case msgTick:
		*m--
		if *m <= 0 {
			return tea.Quit
		}
		return tick
	}
	return nil
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m *model) View(r tea.Renderer) {
	r.Write(fmt.Sprintf("Hi. This program will exit in %d seconds. To quit sooner press any key.\n", *m))
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type msgTick time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return msgTick{}
}
