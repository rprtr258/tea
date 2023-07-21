package tui_daemon_combo

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/muesli/reflow/indent"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

func Main() {
	var (
		daemonMode bool
		showHelp   bool
	)

	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	p := tea.NewProgram(newModel())
	if daemonMode || !isatty.IsTerminal(os.Stdout.Fd()) {
		// If we're in daemon mode don't render the TUI
		p = p.WithoutRenderer()
	} else {
		// If we're in TUI mode, discard log output
		log.SetOutput(io.Discard)
	}

	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}

type result struct {
	duration time.Duration
	emoji    string
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

func (m *model) Init() tea.Cmd {
	log.Println("Starting work...")
	return tea.Batch(
		m.spinner.Tick,
		runPretendProcess,
	)
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		m.quitting = true
		return tea.Quit
	case spinner.TickMsg:
		return m.spinner.Update(msg)
	case processFinishedMsg:
		d := time.Duration(msg)
		res := result{emoji: randomEmoji(), duration: d}
		log.Printf("%s Job finished in %s", res.emoji, res.duration)
		m.results = append(m.results[1:], res)
		return runPretendProcess
	default:
		return nil
	}
}

func (m *model) View(r tea.Renderer) {
	s := "\n" +
		m.spinner.View() + " Doing some work...\n\n"

	for _, res := range m.results {
		if res.duration == 0 {
			s += "........................\n"
		} else {
			s += fmt.Sprintf("%s Job finished in %s\n", res.emoji, res.duration)
		}
	}

	s += helpStyle("\nPress any key to exit\n")

	if m.quitting {
		s += "\n"
	}

	r.Write(indent.String(s, 1))
	return
}

// processFinishedMsg is sent when a pretend process completes.
type processFinishedMsg time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess() tea.Msg {
	pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
	time.Sleep(pause)
	return processFinishedMsg(pause)
}

func randomEmoji() string {
	emojis := []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
	return string(emojis[rand.Intn(len(emojis))]) // nolint:gosec
}
