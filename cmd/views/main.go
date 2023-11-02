package views

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/rprtr258/tea/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/rprtr258/scuf"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

const (
	_progressBarWidth  = 71
	_progressCharFull  = "█"
	_progressCharEmpty = "░"
)

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return func(s string) string {
		return scuf.String(s, scuf.Modifier(styles.FgColor(color)))
	}
}

// General stuff for styling the view
var (
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(_progressCharEmpty)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp = func(colorA, colorB string, steps float64) []string { // Generate a blend of colors.
		cA, _ := colorful.Hex(colorA)
		cB, _ := colorful.Hex(colorB)

		s := make([]string, 0, int(steps)+1)
		for i := 0.0; i < steps; i++ {
			c := cA.BlendLuv(cB, i/steps)
			s = append(s, colorToHex(c))
		}
		return s
	}("#B14FFF", "#00FFA3", _progressBarWidth)
)

type msgTick struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return msgTick{}
	})
}

type msgFrame struct{}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return msgFrame{}
	})
}

type model struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(tick())
}

// Main update function.
func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.MsgKey); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			f(tea.Quit)
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		m.updateChoices(msg, f)
		return
	}

	m.updateChosen(msg, f)
}

// The main view, which just calls the appropriate sub-view
func (m *model) View(vb tea.Viewbox) {
	if m.Quitting {
		vb.PaddingLeft(2).PaddingTop(1).WriteLine("See you later!")
		return
	}

	var s string
	if !m.Chosen {
		s = m.choicesView()
	} else {
		s = m.chosenView()
	}
	vb.PaddingTop(1).WriteLine(s)
	// TODO: indent
	// r.Write(indent.String("\n"+s+"\n\n", 2))
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func (m *model) updateChoices(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			f(frame())
		}

	case msgTick:
		if m.Ticks == 0 {
			m.Quitting = true
			f(tea.Quit)
			return
		}
		m.Ticks--
		f(tick())
	}
}

// Update loop for the second view after a choice has been made
func (m *model) updateChosen(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg.(type) {
	case msgFrame:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 3
				f(tick())
				return
			}

			f(frame())
		}
	case msgTick:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				f(tea.Quit)
				return
			}

			m.Ticks--
			f(tick())
		}
	}
}

// Sub-views

// The first view, where you're choosing a task
func (m *model) choicesView() string {
	c := m.Choice

	tpl := "What to do today?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Plant carrots", c == 0),
		checkbox("Go to the market", c == 1),
		checkbox("Read something", c == 2),
		checkbox("See friends", c == 3),
	)

	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The second view, after a task has been chosen
func (m *model) chosenView() string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keyword("libgarden"), keyword("vegeutils"))
	case 1:
		msg = fmt.Sprintf("A trip to the market?\n\nOkay, then we should install %s and %s...", keyword("marketkit"), keyword("libshopping"))
	case 2:
		msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keyword("actual library"))
	default:
		msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keyword("social-skills"), keyword("conversationutils"))
	}

	label := "Downloading..."
	if m.Loaded {
		label = fmt.Sprintf("Downloaded. Exiting in %s seconds...", colorFg(strconv.Itoa(m.Ticks), "79"))
	}

	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(percent float64) string {
	w := float64(_progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += scuf.String(_progressCharFull, scuf.FgRGB(scuf.MustParseHexRGB(ramp[i])))
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return scuf.String(val, scuf.FgRGB(scuf.MustParseHexRGB(color)))
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and 1.
func colorFloatToHex(f float64) string {
	s := strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return s
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{0, false, 10, 0, 0, false, false}).Run()
	return err
}
