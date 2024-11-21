package send_msg //nolint:revive,stylecheck

// A simple example that shows how to send messages to a Tea program
// from outside the program using Program.Send(Msg).

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/styles"
)

var (
	spinnerStyle  = styles.Style{}.Foreground(styles.FgColor("63"))
	helpStyle     = styles.Style{}.Foreground(styles.FgColor("241")) // .Margin(1, 0)
	dotStyle      = helpStyle.Copy()                                 // .UnsetMargins()
	durationStyle = dotStyle.Copy()
	// appStyle      = styles.Style{} // .Margin(1, 2, 0, 2)
)

type msgResult struct {
	duration time.Duration
	food     string
}

func (r msgResult) String() string {
	if r.duration == 0 {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("🍔 Ate %s %s", r.food,
		durationStyle.Render(r.duration.String()))
}

type model struct {
	spinner  spinner.Model
	results  []msgResult
	quitting bool
}

func newModel() *model {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
	return &model{
		spinner: s,
		results: make([]msgResult, numLastResults),
	}
}

func (m *model) Init(c tea.Context[*model]) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	m.spinner.CmdTick(ctxSpinner)
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	switch msg := msg.(type) {
	case tea.MsgKey:
		m.quitting = true
		c.Dispatch(tea.Quit)
	case msgResult:
		m.results = append(m.results[1:], msg)
	case spinner.MsgTick:
		m.spinner.Update(ctxSpinner, msg)
	}
}

func (m *model) View(vb tea.Viewbox) {
	if m.quitting {
		vb.WriteLine("That’s all for today!")
	} else {
		m.spinner.View(vb)
		vb.PaddingLeft(4).WriteLine(" Eating food...")
	}

	for i, res := range m.results {
		vb.PaddingTop(2 + i).WriteLine(res.String())
	}

	if !m.quitting {
		vb.PaddingTop(2 + len(m.results)).Styled(helpStyle).WriteLine("Press any key to exit")
	}
}

func randomFood() string {
	food := []string{
		"an apple", "a pear", "a gherkin", "a party gherkin",
		"a kohlrabi", "some spaghetti", "tacos", "a currywurst", "some curry",
		"a sandwich", "some peanut butter", "some cashews", "some ramen",
	}
	return food[rand.Intn(len(food))] // nolint:gosec
}

func Main(ctx context.Context) error {
	p := tea.NewProgram2(ctx, newModel())

	// Simulate activity
	go func() {
		for {
			pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
			<-time.After(pause)

			// Send the Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(msgResult{food: randomFood(), duration: pause})
		}
	}()

	_, err := p.Run()
	return err
}
