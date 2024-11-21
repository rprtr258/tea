package package_manager //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/progress"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/styles"
)

var packages = [...]string{
	"babys-monads",
	"bad-kitty",
	"cashew-apple",
	"chai",
	"coffee-CUPS",
	"currykit",
	"currywurst-devel",
	"eggy",
	"fullenglish",
	"hojicha",
	"jalapeño",
	"libesszet",
	"libgardening",
	"libpurring",
	"libtacos",
	"libyuzu",
	"licorice-utils",
	"molasses-utils",
	"old-socks-devel",
	"party-gherkin",
	"rock-lobster",
	"schnurrkit",
	"snow-peas",
	"spicerack",
	"standmixer",
	"vegeutils",
	"xkohlrabi",
	"xmodmeow",
	"zeichenorientierte-benutzerschnittstellen",
}

func getPackages() []string {
	pkgs := packages

	rand.Shuffle(len(pkgs), func(i, j int) {
		pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
	})

	for k := range pkgs {
		major, minor, patch := rand.Intn(10), rand.Intn(10), rand.Intn(10) //nolint:gosec
		pkgs[k] += fmt.Sprintf("-%d.%d.%d", major, minor, patch)
	}
	return pkgs[:]
}

type model struct {
	packages []string
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

type cmd = tea.Msg2[*model]

var (
	currentPkgNameStyle = styles.Style{}.
				Foreground(styles.FgColor("211"))
	doneStyle = styles.Style{}
	// Margin(1, 2)
	checkMark = styles.Style{}.
			Foreground(styles.FgColor("42")).
			SetString("✓")
)

func newModel() *model {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(styles.Style{}.Foreground(styles.FgColor("63"))),
	)
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

func (m *model) Init(c tea.Context[*model]) {
	downloadAndInstall(c, m.packages[m.index])
	m.spinner.CmdTick(tea.Of(c, func(*model) *spinner.Model { return &m.spinner }))
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	ctxProgress := tea.Of(c, func(*model) *progress.Model { return &m.progress })
	ctxSpinner := tea.Of(c, func(*model) *spinner.Model { return &m.spinner })
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.width, m.height = msg.Width, msg.Height
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			c.Dispatch(tea.Quit)
			return
		}
	case msgInstalledPkg:
		if m.index >= len(m.packages)-1 {
			// Everything's been installed. We're done!
			m.done = true
			c.Dispatch(tea.Quit)
			return
		}

		// Update progress bar
		m.progress.SetPercent(ctxProgress, float64(m.index)/float64(len(m.packages)-1))

		m.index++
		c.Dispatch(tea.Printf("%s %s", checkMark, m.packages[m.index])) // print success message above our program
		downloadAndInstall(c, m.packages[m.index])                      // download the next package
	case spinner.MsgTick:
		m.spinner.Update(ctxSpinner, msg)
	case progress.MsgFrame:
		m.progress.Update(ctxProgress, msg)
	}
}

func (m *model) View(vb tea.Viewbox) {
	n := len(m.packages)

	if m.done {
		vb.Styled(doneStyle).WriteLine(fmt.Sprintf("Done! Installed %d packages.", n))
		return
	}

	// w := styles.Width(fmt.Sprintf("%d", n))
	// pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n-1)

	// m.spinner.View(vb)
	// spin := " "
	// x := vb.WriteLine(0, 3, spin)

	// prog := m.progress.View()
	// cellsAvail := max(0, m.width-styles.Width(spin+prog+pkgCount))

	// pkgName := currentPkgNameStyle.Render(m.packages[m.index])
	// info := styles.Style{}.MaxWidth(cellsAvail).Render("Installing " + pkgName)
	// x = vb.WriteLine(0, x, info)

	// x += max(0, m.width-styles.Width(spin+info+prog+pkgCount))

	// x = vb.WriteLine(0, x, prog)
	// vb.WriteLine(0, x, pkgCount)
}

type msgInstalledPkg string

func downloadAndInstall(c tea.Context[*model], pkg string) {
	// This is where you'd do i/o stuff to download and install packages. In
	// our case we're just pausing for a moment to simulate the process.
	d := time.Millisecond * time.Duration(rand.Intn(500)) //nolint:gosec
	// TODO: tea.Tick(d)
	c.F(func() cmd {
		return func(m *model) {
			<-time.After(d)
			m.Update(c, msgInstalledPkg(pkg))
		}
	})
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram2(ctx, newModel()).Run()
	return err
}
