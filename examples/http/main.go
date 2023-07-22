package http

// A simple program that makes a GET request and prints the response status.

import (
	"context"
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

type (
	msgStatus int
	msgErr    struct{ err error }
)

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{checkServer}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return []tea.Cmd{tea.Quit}
		default:
			return nil
		}

	case msgStatus:
		m.status = int(msg)
		return []tea.Cmd{tea.Quit}

	case msgErr:
		m.err = msg.err
		return nil

	default:
		return nil
	}
}

func (m *model) View(r tea.Renderer) {
	s := fmt.Sprintf("Checking %s...", url)
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
	}
	r.Write(s + "\n")
}

func checkServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(url)
	if err != nil {
		return msgErr{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return msgStatus(res.StatusCode)
}

func Main() {
	p := tea.NewProgram(context.Background(), &model{})
	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
