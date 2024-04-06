package paginator

// A simple program demonstrating paginator component.

import (
	"context"
	"fmt"

	"github.com/rprtr258/fun/iter"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/paginator"
	"github.com/rprtr258/tea/styles"
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
	p.ActiveDotStyle = styles.Style{}.Foreground(styles.FgAdaptiveColor("235", "252"))
	p.InactiveDotStyle = styles.Style{}.Foreground(styles.FgAdaptiveColor("250", "238"))
	p.InactiveDot = '•'
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
	switch msg := msg.(type) { //nolint:gocritic
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
	vb = vb.PaddingTop(1)
	vb.WriteLine("  Paginator Example")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i, item := range m.items[start:end] {
		vb.PaddingTop(2 + 2*i).WriteLine("  • " + item)
	}
	y := 2 + 2*(end-start)
	m.paginator.View(vb.PaddingTop(y).PaddingLeft(2))
	vb.PaddingTop(y + 1).WriteLine("  h/l ←/→ page • q: quit")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
