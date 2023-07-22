package composable_views

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/bubbles/timer"
	"github.com/rprtr258/tea/lipgloss"
)

/*
This example assumes an existing understanding of commands and messages. If you
haven't already read our tutorials on the basics of Bubble Tea and working
with commands, we recommend reading those first.

Find them at:
https://github.com/rprtr258/tea/tree/master/tutorials/commands
https://github.com/rprtr258/tea/tree/master/tutorials/basics
*/

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	spinnerView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state   sessionState
	timer   timer.Model
	spinner spinner.Model
	index   int
}

func newModel(timeout time.Duration) *mainModel {
	return &mainModel{
		state:   timerView,
		timer:   timer.New(timeout),
		spinner: spinner.New(),
	}
}

func (m *mainModel) Init() []tea.Cmd {
	// start the timer and spinner on program start
	return append(m.timer.Init(), m.spinner.Tick)
}

func (m *mainModel) Update(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			return []tea.Cmd{tea.Quit}
		case "tab":
			if m.state == timerView {
				m.state = spinnerView
			} else {
				m.state = timerView
			}
		case "n":
			if m.state == timerView {
				m.timer = timer.New(defaultTime)
				cmds = append(cmds, m.timer.Init()...)
			} else {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
		}
		switch m.state {
		// update whichever model is focused
		case spinnerView:
			cmds = append(cmds, m.spinner.Update(msg)...)
		default:
			cmds = append(cmds, m.timer.Update(msg)...)
		}
	case spinner.MsgTick:
		cmds = append(cmds, m.spinner.Update(msg)...)
	case timer.MsgTick:
		cmds = append(cmds, m.timer.Update(msg)...)
	}
	return cmds
}

func (m *mainModel) View(r tea.Renderer) {
	var s string
	model := m.currentFocusedModel()
	if m.state == timerView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), modelStyle.Render(m.spinner.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), focusedModelStyle.Render(m.spinner.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	r.Write(s)
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "spinner"
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func Main() {
	p := tea.NewProgram(context.Background(), newModel(defaultTime))

	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
