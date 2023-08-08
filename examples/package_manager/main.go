package package_manager

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/progress"
	"github.com/rprtr258/tea/bubbles/spinner"
	"github.com/rprtr258/tea/lipgloss"
)

type model struct {
	packages []string
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newModel() *model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return &model{
		packages: getPackages(),
		spinner:  s,
		progress: progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
			progress.WithoutPercentage(),
		),
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(
		downloadAndInstall(m.packages[m.index]),
		m.spinner.CmdTick,
	)
}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.width, m.height = msg.Width, msg.Height
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			f(tea.Quit)
			return
		}
	case msgInstalledPkg:
		if m.index >= len(m.packages)-1 {
			// Everything's been installed. We're done!
			m.done = true
			f(tea.Quit)
			return
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.packages)-1))

		m.index++
		f(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.packages[m.index]), // print success message above our program
			downloadAndInstall(m.packages[m.index]),             // download the next package
		)
		return
	case spinner.MsgTick:
		f(m.spinner.Update(msg)...)
		return
	case progress.MsgFrame:
		f(m.progress.Update(msg)...)
		return
	}
}

func (m *model) View(r tea.Renderer) {
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		r.Write(doneStyle.Render(fmt.Sprintf("Done! Installed %d packages.\n", n)))
		return
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n-1)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.packages[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Installing " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	r.Write(spin + info + gap + prog + pkgCount)
}

type msgInstalledPkg string

func downloadAndInstall(pkg string) tea.Cmd {
	// This is where you'd do i/o stuff to download and install packages. In
	// our case we're just pausing for a moment to simulate the process.
	d := time.Millisecond * time.Duration(rand.Intn(500)) //nolint:gosec
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return msgInstalledPkg(pkg)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
