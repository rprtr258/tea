package http

// A simple program that makes a GET request and prints the response status.

import (
	"context"
	"fmt"
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

func (m *model) Init(f func(...tea.Cmd)) {
	f(checkServer)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			f(tea.Quit)
		}
	case msgStatus:
		m.status = int(msg)
		f(tea.Quit)
	case msgErr:
		m.err = msg.err
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

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
