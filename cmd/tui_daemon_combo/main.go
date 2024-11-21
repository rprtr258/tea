package tui_daemon_combo //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/styles"
)

var (
	_helpStyle = styles.Style{}.Foreground(styles.FgColor("241"))
	_emojis    = []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
)

type result struct {
	duration time.Duration
	emoji    rune
}

// msgProcessFinished is sent when a pretend process completes.
type msgProcessFinished time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess(c tea.Context[*model]) {
	c.F(func() tea.Msg2[*model] {
		pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec // not needed
		time.Sleep(pause)
		msg := msgProcessFinished(pause)
		return func(m *model) {
			m.Update(c, msg)
		}
	})
}

type model struct {
	spinner  spinner.Model
	results  []result
	quitting bool
}

func newModel() *model {
	const _showLastResults = 5

	return &model{
		spinner:  spinner.New(spinner.WithStyle(styles.Style{}.Foreground(styles.FgColor("206")))),
		results:  make([]result, _showLastResults),
		quitting: false,
	}
}

func (m *model) Init(c tea.Context[*model]) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	m.spinner.CmdTick(ctxSpinner)
	runPretendProcess(c)
}

func randomEmoji() rune {
	return _emojis[rand.Intn(len(_emojis))] //nolint:gosec // not needed
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	switch msg := msg.(type) {
	case tea.MsgKey:
		m.quitting = true
		c.Dispatch(tea.Quit)
	case spinner.MsgTick:
		m.spinner.Update(ctxSpinner, msg)
	case msgProcessFinished:
		d := time.Duration(msg)
		res := result{emoji: randomEmoji(), duration: d}
		m.results = append(m.results[1:], res)
		runPretendProcess(c)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb = vb.PaddingTop(1).PaddingLeft(1)

	m.spinner.View(vb)
	vb.PaddingLeft(2).WriteLine(" Doing some work...")

	vb = vb.PaddingTop(2)
	for i, res := range m.results {
		vb.PaddingTop(i).WriteLine(fun.IF(
			res.duration == 0,
			"........................",
			fmt.Sprintf("%c Job finished in %s", res.emoji, res.duration),
		))
	}

	vb = vb.PaddingTop(len(m.results) + 1)
	vb.Styled(_helpStyle).WriteLine("Press any key to exit")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram2(ctx, newModel()).Run()
	return err
}
