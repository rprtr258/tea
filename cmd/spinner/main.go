package spinner

// A simple program demonstrating spinner component.

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/styles"
)

type model struct {
	spinner spinner.Model
}

func initialModel() *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Style{}.Foreground(styles.FgColor("205"))
	return &model{
		spinner: s,
	}
}

func (m *model) Init(c tea.Context[*model]) {
	ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
	m.spinner.CmdTick(ctxSpinner)
}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			c.Dispatch(tea.Quit)
		}
	default:
		ctxSpinner := tea.Of(c, func(m *model) *spinner.Model { return &m.spinner })
		m.spinner.Update(ctxSpinner, msg)
	}
}

func (m *model) View(vb tea.Viewbox) {
	m.spinner.View(vb.Padding(tea.PaddingOptions{
		Left: 1,
		Top:  2,
	}))
	vb.PaddingLeft(5).PaddingTop(2).WriteLine(" Loading forever...press q to quit")
	// if m.quitting {
	// 	r.Write("\n")
	// }
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram2(ctx, initialModel()).Run()
	return err
}
