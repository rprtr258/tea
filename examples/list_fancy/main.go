package list_fancy

import (
	"context"
	"math/rand"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/bubbles/list"
	"github.com/rprtr258/tea/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
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
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list          list.Model[item]
	itemGenerator *randomItemGenerator
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func newModel() *model {
	var (
		itemGenerator randomItemGenerator
		delegateKeys  = newDelegateKeyMap()
		listKeys      = newListKeyMap()
	)

	// Make initial list of items
	const numItems = 24
	items := make([]item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = itemGenerator.next()
	}

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

func (m *model) Init() []tea.Cmd {
	return []tea.Cmd{tea.EnterAltScreen}
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.MsgKey:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			return m.list.ToggleSpinner()
		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			newItem := m.itemGenerator.next()
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			return append(insCmd, statusCmd)
		}
	}

	// This will also call our delegate's update function.
	return append(cmds, m.list.Update(msg)...)
}

func (m *model) View(r tea.Renderer) {
	r.Write(appStyle.Render(m.list.View()))
}

func Main(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())

	_, err := tea.NewProgram(ctx, newModel()).Run()
	return err
}
