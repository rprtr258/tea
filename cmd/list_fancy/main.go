package list_fancy //nolint:revive,stylecheck

import (
	"context"
	"math/rand"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/list"
	"github.com/rprtr258/tea/styles"
)

var (
	appPadding = tea.PaddingOptions{
		Top:    1,
		Left:   2,
		Bottom: 1,
	}

	titleStyle = styles.Style{}.
			Foreground(styles.FgColor("#FFFDF5")).
			Background(styles.BgColor("#25A065"))
		// Padding(0, 1)

	statusMessageStyle = styles.Style{}.
				Foreground(styles.FgAdaptiveColor("#04B575", "#04B575")).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.Binding{
			Keys: []string{"a"},
			Help: key.Help{"a", "add item"},
		},
		toggleSpinner: key.Binding{
			Keys: []string{"s"},
			Help: key.Help{"s", "toggle spinner"},
		},
		toggleTitleBar: key.Binding{
			Keys: []string{"T"},
			Help: key.Help{"T", "toggle title"},
		},
		toggleStatusBar: key.Binding{
			Keys: []string{"S"},
			Help: key.Help{"S", "toggle status"},
		},
		togglePagination: key.Binding{
			Keys: []string{"P"},
			Help: key.Help{"P", "toggle pagination"},
		},
		toggleHelpMenu: key.Binding{
			Keys: []string{"H"},
			Help: key.Help{"H", "toggle help"},
		},
	}
}

type model struct {
	list          list.Model[item]
	itemGenerator *randomItemGenerator
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func newModel() *model {
	var itemGenerator randomItemGenerator
	// Make initial list of items
	const numItems = 24
	items := make([]item, numItems)
	for i := range items {
		items[i] = itemGenerator.next()
	}

	listKeys := newListKeyMap()
	// Setup list
	delegate := newItemDelegate[item](delegateKeys)
	groceryList := list.New[item](items, delegate, 0, 0)
	groceryList.Title = "Groceries"
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return &model{
		list:          groceryList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &itemGenerator,
	}
}

func (m *model) Init(f func(...tea.Cmd)) {
	f(tea.EnterAltScreen)
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.list.SetSize(msg.Width-appPadding.Left-appPadding.Right, msg.Height-appPadding.Top-appPadding.Bottom)
	case tea.MsgKey:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			yield(m.list.ToggleSpinner()...)
		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			newItem := m.itemGenerator.next()
			yield(m.list.InsertItem(0, newItem)...)
			yield(m.list.CmdNewStatusMessage(statusMessageStyle("Added " + newItem.Title())))
		}
	}

	// This will also call our delegate's update function.
	m.list.Update(msg, yield)
}

func (m *model) View(vb tea.Viewbox) {
	m.list.View(vb.Padding(appPadding))
}

func Main(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())

	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
