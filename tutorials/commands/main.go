package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rprtr258/tea"
)

const url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url) //nolint:noctx
	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error } // TODO: rename all msgs to msgXxx

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m *model) Init() tea.Cmd {
	return checkServer
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return tea.Quit

	case errMsg:
		m.err = msg
		return tea.Quit

	case tea.MsgKey:
		if msg.Type == tea.KeyCtrlC {
			return tea.Quit
		}
	}

	return nil
}

func (m *model) View(r tea.Renderer) {
	if m.err != nil {
		r.Write(fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err))
		return
	}

	s := fmt.Sprintf("Checking %s ... ", url)
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	r.Write("\n" + s + "\n\n")
}

func main() {
	if _, err := tea.NewProgram(&model{}).Run(); err != nil {
		log.Fatalln("Uh oh, there was an error:", err.Error())
	}
}
