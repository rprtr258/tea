package paginator

// A simple program demonstrating paginator component.

import (
	"context"
	"fmt"

	"github.com/rprtr258/fun/iter"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/paginator"
)

func newModel() *model {
	items := iter.Map(
		iter.FromRange(0, 100, 1),
		func(i int) string {
			return fmt.Sprintf("Item %d", i)
		}).ToSlice()

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	// p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	// p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(items))

	return &model{
		paginator: p,
		items:     items,
	}
}

type model struct {
	items     []string
	paginator paginator.Model
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			f(tea.Quit)
			return
		}
	}
	f(m.paginator.Update(msg)...)
}

func (m *model) View(vb tea.Viewbox) {
	vb.WriteLine(1, 0, "  Paginator Example")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i, item := range m.items[start:end] {
		vb.WriteLine(3+2*i, 0, "  • "+item)
	}
	y := 3 + 2*(end-start)
	m.paginator.View(vb.Padding(tea.PaddingOptions{Top: y, Left: 2}))
	vb.WriteLine(y+2, 0, "  h/l ←/→ page • q: quit")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
