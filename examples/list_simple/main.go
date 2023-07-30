package list_simple

import (
	"context"
	"fmt"
	"io"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/list"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyle.PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyle.HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                                     { return 1 }
func (d itemDelegate) Spacing() int                                    { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model[item]) []tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m *list.Model[item], index int, i item) {
	str := fmt.Sprintf("%d. %s", index+1, i)

	var res string
	if index == m.Index() {
		res = selectedItemStyle.Render("> " + str)
	} else {
		res = itemStyle.Render(str)
	}

	fmt.Fprint(w, res)
}

type model struct {
	list     list.Model[item]
	choice   string
	quitting bool
}

func (m *model) Init() []tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.list.SetWidth(msg.Width)
		return nil
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return []tea.Cmd{tea.Quit}
		case "enter":
			i, ok := m.list.SelectedItem()
			if ok {
				m.choice = string(i)
			}
			return []tea.Cmd{tea.Quit}
		}
	}

	return m.list.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	if m.choice != "" {
		r.Write(quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice)))
		return
	}

	if m.quitting {
		r.Write(quitTextStyle.Render("Not hungry? That’s cool."))
		return
	}

	r.Write("\n" + m.list.View())
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
	l.Styles.HelpStyle = helpStyle

	_, err := tea.NewProgram(ctx, &model{list: l}).Run()
	return err
}
