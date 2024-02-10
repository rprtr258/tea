package composable_views //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/box"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/components/timer"
	"github.com/rprtr258/tea/styles"
)

// sessionState is used to track which model is focused
type sessionState uint

const (
	_timerView sessionState = iota
	_spinnerView
)

const _defaultTime = time.Minute

// Available spinners
var spinners = []spinner.Spinner{
	spinner.Line,
	spinner.Dot,
	spinner.MiniDot,
	spinner.Jump,
	spinner.Pulse,
	spinner.Points,
	spinner.Globe,
	spinner.Moon,
	spinner.Monkey,
	spinner.Meter,
	spinner.Hamburger,
	spinner.Ellipsis,
	spinner.Circle,
}

var (
	_styleSpinner = styles.Style{}.Foreground(styles.FgColor("69"))
	_styleHelp    = styles.Style{}.Foreground(styles.FgColor("241"))
)

type mainModel struct {
	state   sessionState
	timer   timer.Model
	spinner spinner.Model
	index   int
}

func newModel(timeout time.Duration) *mainModel {
	return &mainModel{
		state:   _timerView,
		timer:   timer.New(timeout),
		spinner: spinner.New(),
	}
}

func (m *mainModel) Init(f func(...tea.Cmd)) {
	// start the timer and spinner on program start
	m.timer.Init(f)
	f(m.spinner.CmdTick)
}

func (m *mainModel) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			f(tea.Quit)
			return
		case "tab":
			m.state = fun.IF(m.state == _timerView, _spinnerView, _timerView)
		case "n":
			if m.state == _timerView {
				m.timer = timer.New(_defaultTime)
				m.timer.Init(f)
			} else {
				m.index = (m.index + 1) % len(spinners)
				m.spinner = spinner.New()
				m.spinner.Style = _styleSpinner
				m.spinner.Spinner = spinners[m.index]
				f(m.spinner.CmdTick)
			}
		}
		switch m.state {
		// update whichever model is focused
		case _spinnerView:
			m.spinner.Update(msg, f)
		default:
			m.timer.Update(msg, f)
		}
	case spinner.MsgTick:
		m.spinner.Update(msg, f)
	case timer.MsgTick:
		m.timer.Update(msg, f)
	}
}

func (m *mainModel) View(vb tea.Viewbox) {
	const (
		_width  = 15
		_height = 5
	)

	vbTimerSpinner, vbHelp := vb.SplitY2(tea.Fixed(1+_height+1), tea.Fixed(1))

	vbTimer, vbSpinner := vbTimerSpinner.Sub(tea.Rectangle{
		Width: (1 + _width + 1) * 2,
	}).SplitX2(tea.Flex(1), tea.Flex(1))

	focusedTimer := m.state == _timerView
	_styleModel := styles.Style{}.Align(styles.Center, styles.Center)
	box.Box(vbTimer.Styled(_styleModel),
		func(vb tea.Viewbox) {
			vb = vb.Padding(tea.PaddingOptions{Top: 2, Bottom: 2, Left: (_width - 4) / 2})
			m.timer.View(vb)
		},
		fun.IF(focusedTimer, box.NormalBorder, box.HiddenBorder),
		box.BorderMaskAll,
		box.Colors(fun.IF(focusedTimer, styles.FgColor("69"), nil)),
		box.Colors(nil),
	)
	box.Box(vbSpinner.Styled(_styleModel),
		func(vb tea.Viewbox) {
			m.spinner.View(vb.Padding(tea.PaddingOptions{Top: 2, Bottom: 2, Left: (_width - utf8.RuneCountInString(spinners[m.index].Frames[0])) / 2}))
		},
		fun.IF(focusedTimer, box.HiddenBorder, box.NormalBorder),
		box.BorderMaskAll,
		box.Colors(fun.IF(focusedTimer, nil, styles.FgColor("69"))),
		box.Colors(nil),
	)
	vbHelp.
		Styled(_styleHelp).
		WriteLine(fmt.Sprintf("tab: focus next • n: new %s • q: exit", fun.IF(focusedTimer, "timer", "spinner")))
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel(_defaultTime)).Run()
	return err
}
