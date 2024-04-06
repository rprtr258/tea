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
	vb = vb.PaddingTop(1)
	if m.err != nil {
		vb.WriteLineX("We had some trouble: ").WriteLineX(m.err.Error())
		return
	}

	vb = vb.WriteLineX("Checking ").WriteLineX(_url).WriteLineX(" ... ")
	if m.status > 0 {
		vb.WriteLineX(strconv.Itoa(m.status)).WriteLineX(" ").WriteLineX(http.StatusText(m.status)).WriteLineX("!")
	}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}
