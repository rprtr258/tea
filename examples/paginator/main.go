package paginator

// A simple program demonstrating the paginator component from the Bubbles
// component library.

import (
	"context"
	"fmt"
	"strings"

	"github.com/rprtr258/tea/bubbles/paginator"
	"github.com/rprtr258/tea/lipgloss"
	"github.com/samber/lo"

	"github.com/rprtr258/tea"
)

func newModel() *model {
	uuh := [100]struct{}{}
	items := lo.Map(uuh[:], func(_ struct{}, i int) string {
		return fmt.Sprintf("Item %d", i)
	})

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

func (m *model) Init() []tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return []tea.Cmd{tea.Quit}
		}
	}
	return m.paginator.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	var sb strings.Builder
	sb.WriteString("\n  Paginator Example\n\n")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	for _, item := range m.items[start:end] {
		sb.WriteString("  • " + item + "\n\n")
	}
	sb.WriteString("  " + m.paginator.View())
	sb.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	r.Write(sb.String())
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
