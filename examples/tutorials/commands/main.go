package commands

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/rprtr258/tea"
)

const _url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(_url) //nolint:noctx
	if err != nil {
		return msgErr(err)
	}
	defer res.Body.Close() // nolint:errcheck

	return msgStatus(res.StatusCode)
}

type (
	msgStatus int
	msgErr    error
)

func (m *model) Init(f func(...tea.Cmd)) {
	f(checkServer)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case msgStatus:
		m.status = int(msg)
		f(tea.Quit)
	case msgErr:
		m.err = msg
		f(tea.Quit)
	case tea.MsgKey:
		if msg.Type == tea.KeyCtrlC {
			f(tea.Quit)
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.err != nil {
		x := vb.WriteLine(1, 0, "We had some trouble: ")
		vb.WriteLine(1, x, m.err.Error())
		return
	}

	x := vb.WriteLine(1, 0, "Checking ")
	x = vb.WriteLine(1, x, _url)
	x = vb.WriteLine(1, x, " ... ")
	if m.status > 0 {
		x = vb.WriteLine(1, x, strconv.Itoa(m.status))
		x = vb.WriteLine(1, x, " ")
		x = vb.WriteLine(1, x, http.StatusText(m.status))
		vb.WriteLine(1, x, "!")
	}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
