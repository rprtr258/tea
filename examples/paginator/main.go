package paginator

// A simple program demonstrating the paginator component from the Bubbles
// component library.

import (
	"context"
	"fmt"

	"github.com/rprtr258/fun/iter"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/paginator"
	"github.com/rprtr258/tea/lipgloss"
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
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
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

func (m *model) View(r tea.Renderer) {
	r.Write("\n  Paginator Example\n\n")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	for _, item := range m.items[start:end] {
		r.Write("  • ")
		r.Write(item)
		r.Write("\n\n")
	}
	r.Write("  ")
	r.Write(m.paginator.View())
	r.Write("\n\n  h/l ←/→ page • q: quit\n")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
