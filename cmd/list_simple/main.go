package list_simple //nolint:revive,stylecheck

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/list"
	"github.com/rprtr258/tea/styles"
)

var (
	titleStyle        = styles.Style{} // .MarginLeft(2)
	itemStyle         = styles.Style{} // .PaddingLeft(4)
	selectedItemStyle = styles.Style{}. /*.PaddingLeft(2)*/ Foreground(styles.FgColor("170"))
	paginationStyle   = list.DefaultStyle.PaginationStyle // .PaddingLeft(4)
	quitTextStyle     = styles.Style{}                    // .Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                                     { return 1 }
func (d itemDelegate) Spacing() int                                    { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model[item]) []tea.Cmd { return nil }
func (d itemDelegate) Render(vb tea.Viewbox, m *list.Model[item], index int, i item) {
	str := fmt.Sprintf("%d. %s", index+1, i)

	var style styles.Style
	if index == m.Index() {
		style = selectedItemStyle
		vb.Styled(selectedItemStyle).WriteLine("> ")
	} else {
		style = itemStyle
	}
	vb.PaddingLeft(2).Styled(style).WriteLine(str)
}

type model struct {
	list     list.Model[item]
	choice   string
	quitting bool
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.list.SetWidth(msg.Width)
		return
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			f(tea.Quit)
			return
		case "enter":
			i, ok := m.list.SelectedItem()
			if ok {
				m.choice = string(i)
			}
			f(tea.Quit)
			return
		}
	}

	m.list.Update(msg, f)
}

func (m *model) View(vb tea.Viewbox) {
	vb = vb.Padding(tea.PaddingOptions{
		Top:  1,
		Left: 2,
	})
	vb = vb.Sub(tea.Rectangle{
		Width:  vb.Width,
		Height: 13,
	})
	switch {
	case m.choice != "":
		vb.Styled(quitTextStyle).WriteLine(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	case m.quitting:
		vb.Styled(quitTextStyle).WriteLine("Not hungry? That’s cool.")
	default:
		m.list.View(vb.PaddingTop(1))
	}
}

func Main(ctx context.Context) error {
	items := []item{
		"Ramen",
		"Tomato Soup",
		"Hamburgers",
		"Cheeseburgers",
		"Currywurst",
		"Okonomiyaki",
		"Pasta",
		"Fillet Mignon",
		"Caviar",
		"Just Wine",
	}

	const (
		listHeight   = 14
		defaultWidth = 20
	)

	l := list.New[item](items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	_, err := tea.NewProgram(ctx, &model{list: l}).Run()
	return err
}
