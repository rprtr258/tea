package tui_daemon_combo

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	_helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	_emojis    = []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
)

type result struct {
	duration time.Duration
	emoji    rune
}

type model struct {
	spinner  spinner.Model
	results  []result
	quitting bool
}

func newModel() *model {
	const showLastResults = 5

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	return &model{
		spinner: sp,
		results: make([]result, showLastResults),
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	log.Println("Starting work...")
	f(
		m.spinner.CmdTick,
		runPretendProcess,
	)
}

func randomEmoji() rune {
	return _emojis[rand.Intn(len(_emojis))] //nolint:gosec // not needed
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		m.quitting = true
		yield(tea.Quit)
	case spinner.MsgTick:
		yield(m.spinner.Update(msg)...)
	case msgProcessFinished:
		d := time.Duration(msg)
		res := result{emoji: randomEmoji(), duration: d}
		log.Printf("%c Job finished in %s", res.emoji, res.duration)
		m.results = append(m.results[1:], res)
		yield(runPretendProcess)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine(1, 0, m.spinner.View()+" Doing some work...")

	for i, res := range m.results {
		vb.WriteLine(3+i, 0, fun.IF(
			res.duration == 0,
			"........................",
			fmt.Sprintf("%c Job finished in %s", res.emoji, res.duration),
		))
	}

	vb.WriteLine(3+len(m.results), 0, _helpStyle("\nPress any key to exit\n"))

	y := 3 + len(m.results)
	if m.quitting {
		y++
	}

	// TODO: indent
	// r.Write(indent.String(sb.String(), 1))
}

// msgProcessFinished is sent when a pretend process completes.
type msgProcessFinished time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess() tea.Msg {
	pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec // not needed
	time.Sleep(pause)
	return msgProcessFinished(pause)
}

func Main(ctx context.Context) error {
	var daemonMode bool
	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")

	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "show help")

	flag.Parse()

	if showHelp {
		flag.Usage()
		return nil
	}

	p := tea.NewProgram(ctx, newModel())
	if daemonMode || !isatty.IsTerminal(os.Stdout.Fd()) {
		// If we're in daemon mode don't render the TUI
		p = p.WithoutRenderer()
	} else {
		// If we're in TUI mode, discard log output
		log.SetOutput(io.Discard)
	}

	_, err := p.Run()
	return err
}
