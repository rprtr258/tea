package tabs

import (
	"context"
	"strings"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
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

// func tabBorderWithBottom(left, middle, right rune) box.Border {
// 	border := box.RoundedBorder
// 	border.BottomLeft = left
// 	border.Bottom = middle
// 	border.BottomRight = right
// 	return border
// }

var (
	// inactiveTabBorder = tabBorderWithBottom('┴', '─', '┴')
	// activeTabBorder   = tabBorderWithBottom('┘', ' ', '└')
	docStyle = styles.Style{}
	// Padding(1, 2, 1, 2)
	// highlightColor   = styles.FgAdaptiveColor("#874BFD", "#7D56F4")
	inactiveTabStyle = styles.Style{}
	// Border(inactiveTabBorder, true).
	// BorderForeground(highlightColor)
	// Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Copy()
	// Border(activeTabBorder, true)
	windowStyle = styles.Style{}.
		// BorderForeground(highlightColor).
		// Padding(2, 0).
		Align(styles.Center)
	// Border(styles.NormalBorder).
	// UnsetBorderTop()
)

func (m *model) View(vb tea.Viewbox) {
	renderedTabs := fun.Map[[]string](
		func(t string, i int) []string {
			isActive := i == m.activeTab
			style := fun.IF(isActive, activeTabStyle, inactiveTabStyle).Copy()
			// border := style.GetBorderStyle()
			// switch {
			// case i == 0: // first
			// 	border.BottomLeft = fun.IF(isActive, '│', '├')
			// case i == len(m.Tabs)-1: // last
			// 	border.BottomRight = fun.IF(isActive, '│', '┤')
			// }
			// style = style.Border(border)
			return strings.Split(style.Render(t), "\n")
		}, m.Tabs...)

	row := styles.JoinHorizontal(styles.Top, renderedTabs...)
	vb.Styled(docStyle).WriteText(0, 0,
		row+
			"\n"+
			windowStyle.
				// Width(styles.Width(row) /*-windowStyle.GetHorizontalFrameSize()*/).
				Render(m.TabContent[m.activeTab]))
}

func Main(ctx context.Context) error {
	tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	m := &model{Tabs: tabs, TabContent: tabContent}
	_, err := tea.NewProgram(ctx, m).Run()
	return err
}
