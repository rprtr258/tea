package tabs

import (
	"context"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			f(tea.Quit)
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
		}
	}
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle()
	// Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: lipgloss.FgColor("#874BFD"), Dark: lipgloss.FgColor("#7D56F4")}
	inactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor)
		// Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Copy().
			Border(activeTabBorder, true)
	windowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
		// Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder).
		UnsetBorderTop()
)

func (m *model) View(vb tea.Viewbox) {
	renderedTabs := fun.Map[string](
		m.Tabs,
		func(t string, i int) string {
			isActive := i == m.activeTab
			style := fun.IF(isActive, activeTabStyle, inactiveTabStyle).Copy()
			border := style.GetBorderStyle()
			switch {
			case i == 0: // first
				border.BottomLeft = fun.IF(isActive, "│", "├")
			case i == len(m.Tabs)-1: // last
				border.BottomRight = fun.IF(isActive, "│", "┤")
			}
			style = style.Border(border)
			return style.Render(t)
		})

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	vb.Styled(docStyle).WriteText(0, 0,
		row+
			"\n"+
			windowStyle.
				Width(lipgloss.Width(row) /*-windowStyle.GetHorizontalFrameSize()*/).
				Render(m.TabContent[m.activeTab]))
}

func Main(ctx context.Context) error {
	tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	m := &model{Tabs: tabs, TabContent: tabContent}
	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
