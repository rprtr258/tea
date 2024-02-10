// Package list provides a feature-rich Tea component for browsing
// a general purpose list of items. It features optional filtering, pagination,
// help, status messages, and a spinner to indicate activity.
package list

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/muesli/reflow/truncate"
	"github.com/rprtr258/fun"
	"github.com/sahilm/fuzzy"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/help"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/paginator"
	"github.com/rprtr258/tea/components/spinner"
	"github.com/rprtr258/tea/components/textinput"
	"github.com/rprtr258/tea/styles"
)

// Item is an item that appears in the list.
type Item interface {
	// FilterValue is the value we use when filtering against this item when
	// we're filtering the list.
	FilterValue() string
}

// ItemDelegate encapsulates the general functionality for all list items. The
// benefit to separating this logic from the item itself is that you can change
// the functionality of items without changing the actual items themselves.
//
// Note that if the delegate also implements help.KeyMap delegate-related
// help items will be added to the help view.
type ItemDelegate[I Item] interface {
	// Render renders the item's view.
	Render(vb tea.Viewbox, m *Model[I], index int, item I)

	// Height is the height of the list item.
	Height() int

	// Spacing is the size of the horizontal gap between list items in cells.
	Spacing() int

	// Update is the update loop for items. All messages in the list's update
	// loop will pass through here except when the user is setting a filter.
	// Use this method to perform item-level updates appropriate to this
	// delegate.
	Update(msg tea.Msg, m *Model[I]) []tea.Cmd
}

type filteredItem[I Item] struct {
	item    I     // item matched
	matches []int // rune indices of matched items
}

type filteredItems[I Item] []filteredItem[I]

func (f filteredItems[I]) items() []I {
	return fun.Map[I](
		func(v filteredItem[I]) I {
			return v.item
		},
		f...)
}

// MsgFilterMatches contains data about items matched during filtering. The
// message should be routed to Update for processing.
type MsgFilterMatches[I Item] []filteredItem[I]

// FilterFunc takes a term and a list of strings to search through
// (defined by Item#FilterValue).
// It should return a sorted list of ranks.
type FilterFunc func(string, []string) []Rank

// Rank defines a rank for a given item.
type Rank struct {
	// The index of the item in the original input.
	Index int
	// Indices of the actual word that were matched against the filter term.
	MatchedIndexes []int
}

// DefaultFilter uses the sahilm/fuzzy to filter through the list.
// This is set by default.
func DefaultFilter(term string, targets []string) []Rank {
	ranks := fuzzy.Find(term, targets)
	slices.SortStableFunc(ranks, func(i, j fuzzy.Match) int {
		return j.Score - i.Score
	})
	return fun.Map[Rank](
		func(r fuzzy.Match) Rank {
			return Rank{
				Index:          r.Index,
				MatchedIndexes: r.MatchedIndexes,
			}
		},
		ranks...)
}

type msgStatusMessageTimeout struct{}

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

// String returns a human-readable string of the current filter state.
func (f FilterState) String() string {
	return [...]string{
		"unfiltered",
		"filtering",
		"filter applied",
	}[f]
}

// Model contains the state of this component.
type Model[I Item] struct {
	showTitle        bool
	showFilter       bool
	showStatusBar    bool
	showPagination   bool
	showHelp         bool
	filteringEnabled bool

	itemNameSingular string
	itemNamePlural   string

	Title             string
	Styles            Styles
	InfiniteScrolling bool

	// Key mappings for navigating the list.
	KeyMap KeyMap

	// Filter is used to filter the list.
	Filter FilterFunc

	disableQuitKeybindings bool

	// Additional key mappings for the short and full help views. This allows
	// you to add additional key mappings to the help menu without
	// re-implementing the help component. Of course, you can also disable the
	// list's help component and implement a new one if you need more
	// flexibility.
	AdditionalShortHelpKeys func() []key.Binding
	AdditionalFullHelpKeys  func() []key.Binding

	spinner     spinner.Model
	showSpinner bool
	width       int
	height      int
	Paginator   paginator.Model
	cursor      int
	Help        help.Model
	FilterInput textinput.Model
	filterState FilterState

	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration

	statusMessage      string
	statusMessageTimer *time.Timer

	// The master set of items we're working with.
	items []I

	// Filtered items we're currently displaying. Filtering, toggles and so on
	// will alter this slice so we can show what is relevant. For that reason,
	// this field should be considered ephemeral.
	filteredItems filteredItems[I]

	delegate ItemDelegate[I]
}

// New returns a new model with sensible defaults.
func New[I Item](items []I, delegate ItemDelegate[I], width, height int) Model[I] {
	filterInput := textinput.New()
	filterInput.Prompt = "Filter: "
	filterInput.PromptStyle = DefaultStyle.FilterPrompt
	filterInput.Cursor.Style = DefaultStyle.FilterCursor
	filterInput.CharLimit = 64
	filterInput.Focus()

	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDotStyle = DefaultStyle.ActivePaginationDot
	p.InactiveDotStyle = DefaultStyle.InactivePaginationDot
	p.InactiveDot = bullet

	m := Model[I]{
		showTitle:             true,
		showFilter:            true,
		showStatusBar:         true,
		showPagination:        true,
		showHelp:              true,
		itemNameSingular:      "item",
		itemNamePlural:        "items",
		filteringEnabled:      true,
		KeyMap:                DefaultKeyMap,
		Filter:                DefaultFilter,
		Styles:                DefaultStyle,
		Title:                 "List",
		FilterInput:           filterInput,
		StatusMessageLifetime: time.Second,

		width:     width,
		height:    height,
		delegate:  delegate,
		items:     items,
		Paginator: p,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(DefaultStyle.Spinner),
		),
		Help: help.New(),
	}

	m.updatePagination()
	m.updateKeybindings()
	return m
}

// SetFilteringEnabled enables or disables filtering. Note that this is different
// from ShowFilter, which merely hides or shows the input view.
func (m *Model[I]) SetFilteringEnabled(v bool) {
	m.filteringEnabled = v
	if !v {
		m.resetFiltering()
	}
	m.updateKeybindings()
}

// FilteringEnabled returns whether or not filtering is enabled.
func (m *Model[I]) FilteringEnabled() bool {
	return m.filteringEnabled
}

// SetShowTitle shows or hides the title bar.
func (m *Model[I]) SetShowTitle(v bool) {
	m.showTitle = v
	m.updatePagination()
}

// ShowTitle returns whether or not the title bar is set to be rendered.
func (m *Model[I]) ShowTitle() bool {
	return m.showTitle
}

// SetShowFilter shows or hides the filter bar. Note that this does not disable
// filtering, it simply hides the built-in filter view. This allows you to
// use the FilterInput to render the filtering UI differently without having to
// re-implement filtering from scratch.
//
// To disable filtering entirely use EnableFiltering.
func (m *Model[I]) SetShowFilter(v bool) {
	m.showFilter = v
	m.updatePagination()
}

// ShowFilter returns whether or not the filter is set to be rendered. Note
// that this is separate from FilteringEnabled, so filtering can be hidden yet
// still invoked. This allows you to render filtering differently without
// having to re-implement it from scratch.
func (m *Model[I]) ShowFilter() bool {
	return m.showFilter
}

// SetShowStatusBar shows or hides the view that displays metadata about the
// list, such as item counts.
func (m *Model[I]) SetShowStatusBar(v bool) {
	m.showStatusBar = v
	m.updatePagination()
}

// ShowStatusBar returns whether or not the status bar is set to be rendered.
func (m *Model[I]) ShowStatusBar() bool {
	return m.showStatusBar
}

// SetStatusBarItemName defines a replacement for the item's identifier.
// Defaults to item/items.
func (m *Model[I]) SetStatusBarItemName(singular, plural string) {
	m.itemNameSingular = singular
	m.itemNamePlural = plural
}

// StatusBarItemName returns singular and plural status bar item names.
func (m *Model[I]) StatusBarItemName() (string, string) {
	return m.itemNameSingular, m.itemNamePlural
}

// SetShowPagination hides or shows the paginator. Note that pagination will
// still be active, it simply won't be displayed.
func (m *Model[I]) SetShowPagination(v bool) {
	m.showPagination = v
	m.updatePagination()
}

// ShowPagination returns whether the pagination is visible.
func (m *Model[I]) ShowPagination() bool {
	return m.showPagination
}

// SetShowHelp shows or hides the help view.
func (m *Model[I]) SetShowHelp(v bool) {
	m.showHelp = v
	m.updatePagination()
}

// ShowHelp returns whether or not the help is set to be rendered.
func (m *Model[I]) ShowHelp() bool {
	return m.showHelp
}

// Items returns the items in the list.
func (m *Model[I]) Items() []I {
	return m.items
}

// SetItems sets the items available in the list. This returns a command.
func (m *Model[I]) SetItems(items []I) tea.Cmd {
	m.items = items

	var cmd tea.Cmd
	if m.filterState != Unfiltered {
		m.filteredItems = nil
		cmd = cmdFilterItems(*m)
	}

	m.updatePagination()
	m.updateKeybindings()
	return cmd
}

// Select selects the given index of the list and goes to its respective page.
func (m *Model[I]) Select(index int) {
	m.Paginator.Page = index / m.Paginator.PerPage
	m.cursor = index % m.Paginator.PerPage
}

// ResetSelected resets the selected item to the first item in the first page of the list.
func (m *Model[I]) ResetSelected() {
	m.Select(0)
}

// ResetFilter resets the current filtering state.
func (m *Model[I]) ResetFilter() {
	m.resetFiltering()
}

// SetItem replaces an item at the given index. This returns a command.
func (m *Model[I]) SetItem(index int, item I) []tea.Cmd {
	m.items[index] = item

	var cmd []tea.Cmd
	if m.filterState != Unfiltered {
		cmd = append(cmd, cmdFilterItems(*m))
	}

	m.updatePagination()
	return cmd
}

// InsertItem inserts an item at the given index. If the index is out of the upper bound,
// the item will be appended. This returns a command.
func (m *Model[I]) InsertItem(index int, item I) []tea.Cmd {
	m.items = insertItemIntoSlice(m.items, item, index)

	var cmd []tea.Cmd
	if m.filterState != Unfiltered {
		cmd = append(cmd, cmdFilterItems(*m))
	}

	m.updatePagination()
	m.updateKeybindings()
	return cmd
}

// RemoveItem removes an item at the given index. If the index is out of bounds
// this will be a no-op. O(n) complexity, which probably won't matter in the
// case of a TUI.
func (m *Model[I]) RemoveItem(index int) {
	m.items = removeItemFromSlice(m.items, index)
	if m.filterState != Unfiltered {
		m.filteredItems = removeItemFromSlice(m.filteredItems, index)
		if len(m.filteredItems) == 0 {
			m.resetFiltering()
		}
	}
	m.updatePagination()
}

// SetDelegate sets the item delegate.
func (m *Model[I]) SetDelegate(d ItemDelegate[I]) {
	m.delegate = d
	m.updatePagination()
}

// VisibleItems returns the total items available to be shown.
func (m *Model[I]) VisibleItems() []I {
	if m.filterState != Unfiltered {
		return m.filteredItems.items()
	}
	return m.items
}

// SelectedItem returns the current selected item in the list and true if it
// exists. False otherwise.
func (m *Model[I]) SelectedItem() (I, bool) {
	i := m.Index()

	items := m.VisibleItems()
	if i < 0 || i >= len(items) {
		var i I
		return i, false
	}

	return items[i], true
}

// MatchesForItem returns rune positions matched by the current filter, if any.
// Use this to style runes matched by the active filter.
//
// See DefaultItemView for a usage example.
func (m *Model[I]) MatchesForItem(index int) []int {
	if m.filteredItems == nil || index >= len(m.filteredItems) {
		return nil
	}
	return m.filteredItems[index].matches
}

// Index returns the index of the currently selected item as it appears in the
// entire slice of items.
func (m *Model[I]) Index() int {
	return m.Paginator.Page*m.Paginator.PerPage + m.cursor
}

// Cursor returns the index of the cursor on the current page.
func (m *Model[I]) Cursor() int {
	return m.cursor
}

// CursorUp moves the cursor up. This can also move the state to the previous
// page.
func (m *Model[I]) CursorUp() {
	m.cursor--

	// If we're at the start, stop
	if m.cursor < 0 && m.Paginator.Page == 0 {
		// if infinite scrolling is enabled, go to the last item
		if m.InfiniteScrolling {
			m.Paginator.Page = m.Paginator.TotalPages - 1
			m.cursor = m.Paginator.ItemsOnPage(len(m.VisibleItems())) - 1
			return
		}

		m.cursor = 0
		return
	}

	// Move the cursor as normal
	if m.cursor >= 0 {
		return
	}

	// Go to the previous page
	m.Paginator.PrevPage()
	m.cursor = m.Paginator.ItemsOnPage(len(m.VisibleItems())) - 1
}

// CursorDown moves the cursor down. This can also advance the state to the
// next page.
func (m *Model[I]) CursorDown() {
	itemsOnPage := m.Paginator.ItemsOnPage(len(m.VisibleItems()))

	m.cursor++

	// If we're at the end, stop
	if m.cursor < itemsOnPage {
		return
	}

	// Go to the next page
	if !m.Paginator.OnLastPage() {
		m.Paginator.NextPage()
		m.cursor = 0
		return
	}

	// During filtering the cursor position can exceed the number of
	// itemsOnPage. It's more intuitive to start the cursor at the
	// topmost position when moving it down in this scenario.
	if m.cursor > itemsOnPage {
		m.cursor = 0
		return
	}

	m.cursor = itemsOnPage - 1

	// if infinite scrolling is enabled, go to the first item
	if m.InfiniteScrolling {
		m.Paginator.Page = 0
		m.cursor = 0
	}
}

// PrevPage moves to the previous page, if available.
func (m *Model[I]) PrevPage() {
	m.Paginator.PrevPage()
}

// NextPage moves to the next page, if available.
func (m *Model[I]) NextPage() {
	m.Paginator.NextPage()
}

// FilterState returns the current filter state.
func (m *Model[I]) FilterState() FilterState {
	return m.filterState
}

// FilterValue returns the current value of the filter.
func (m *Model[I]) FilterValue() string {
	return m.FilterInput.Value()
}

// SettingFilter returns whether or not the user is currently editing the
// filter value. It's purely a convenience method for the following:
//
//	m.FilterState() == Filtering
//
// It's included here because it's a common thing to check for when
// implementing this component.
func (m *Model[I]) SettingFilter() bool {
	return m.filterState == Filtering
}

// IsFiltered returns whether or not the list is currently filtered.
// It's purely a convenience method for the following:
//
//	m.FilterState() == FilterApplied
func (m *Model[I]) IsFiltered() bool {
	return m.filterState == FilterApplied
}

// Width returns the current width setting.
func (m *Model[I]) Width() int {
	return m.width
}

// Height returns the current height setting.
func (m *Model[I]) Height() int {
	return m.height
}

// SetSpinner allows to set the spinner style.
func (m *Model[I]) SetSpinner(spinner spinner.Spinner) {
	m.spinner.Spinner = spinner
}

// ToggleSpinner toggles the spinner. Note that this also returns a command.
func (m *Model[I]) ToggleSpinner() []tea.Cmd {
	if !m.showSpinner {
		return []tea.Cmd{m.StartSpinner()}
	}

	m.StopSpinner()
	return nil
}

// StartSpinner starts the spinner. Note that this returns a command.
func (m *Model[I]) StartSpinner() tea.Cmd {
	m.showSpinner = true
	return m.spinner.CmdTick
}

// StopSpinner stops the spinner.
func (m *Model[I]) StopSpinner() {
	m.showSpinner = false
}

// DisableQuitKeybindings is a helper for disabling the keybindings used for quitting,
// in case you want to handle this elsewhere in your application.
func (m *Model[I]) DisableQuitKeybindings() {
	m.disableQuitKeybindings = true
	m.KeyMap.Quit.SetEnabled(false)
	m.KeyMap.ForceQuit.SetEnabled(false)
}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *Model[I]) CmdNewStatusMessage(s string) tea.Cmd {
	m.statusMessage = s
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	// Wait for timeout
	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return msgStatusMessageTimeout{}
	}
}

// SetSize sets the width and height of this component.
func (m *Model[I]) SetSize(width, height int) {
	m.setSize(width, height)
}

// SetWidth sets the width of this component.
func (m *Model[I]) SetWidth(v int) {
	m.setSize(v, m.height)
}

// SetHeight sets the height of this component.
func (m *Model[I]) SetHeight(v int) {
	m.setSize(m.width, v)
}

func (m *Model[I]) setSize(width, height int) {
	promptWidth := styles.Width(m.Styles.Title.Render(m.FilterInput.Prompt))

	m.width = width
	m.height = height
	m.Help.Width = width
	m.FilterInput.Width = width - promptWidth - styles.Width(m.spinnerView())
	m.updatePagination()
}

func (m *Model[I]) resetFiltering() {
	if m.filterState == Unfiltered {
		return
	}

	m.filterState = Unfiltered
	m.FilterInput.Reset()
	m.filteredItems = nil
	m.updatePagination()
	m.updateKeybindings()
}

func (m *Model[I]) itemsAsFilterItems() filteredItems[I] {
	fi := make([]filteredItem[I], len(m.items))
	for i, item := range m.items {
		fi[i] = filteredItem[I]{
			item: item,
		}
	}
	return fi
}

// Set keybindings according to the filter state.
func (m *Model[I]) updateKeybindings() {
	switch m.filterState {
	case Filtering:
		m.KeyMap.CursorUp.SetEnabled(false)
		m.KeyMap.CursorDown.SetEnabled(false)
		m.KeyMap.NextPage.SetEnabled(false)
		m.KeyMap.PrevPage.SetEnabled(false)
		m.KeyMap.GoToStart.SetEnabled(false)
		m.KeyMap.GoToEnd.SetEnabled(false)
		m.KeyMap.Filter.SetEnabled(false)
		m.KeyMap.ClearFilter.SetEnabled(false)
		m.KeyMap.CancelWhileFiltering.SetEnabled(true)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(m.FilterInput.Value() != "")
		m.KeyMap.Quit.SetEnabled(false)
		m.KeyMap.ShowFullHelp.SetEnabled(false)
		m.KeyMap.CloseFullHelp.SetEnabled(false)

	default:
		hasItems := len(m.items) != 0
		m.KeyMap.CursorUp.SetEnabled(hasItems)
		m.KeyMap.CursorDown.SetEnabled(hasItems)

		hasPages := m.Paginator.TotalPages > 1
		m.KeyMap.NextPage.SetEnabled(hasPages)
		m.KeyMap.PrevPage.SetEnabled(hasPages)

		m.KeyMap.GoToStart.SetEnabled(hasItems)
		m.KeyMap.GoToEnd.SetEnabled(hasItems)

		m.KeyMap.Filter.SetEnabled(m.filteringEnabled && hasItems)
		m.KeyMap.ClearFilter.SetEnabled(m.filterState == FilterApplied)
		m.KeyMap.CancelWhileFiltering.SetEnabled(false)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(false)
		m.KeyMap.Quit.SetEnabled(!m.disableQuitKeybindings)

		if m.Help.ShowAll {
			m.KeyMap.ShowFullHelp.SetEnabled(true)
			m.KeyMap.CloseFullHelp.SetEnabled(true)
		} else {
			minHelp := countEnabledBindings(m.FullHelp()) > 1
			m.KeyMap.ShowFullHelp.SetEnabled(minHelp)
			m.KeyMap.CloseFullHelp.SetEnabled(minHelp)
		}
	}
}

// Update pagination according to the amount of items for the current state.
func (m *Model[I]) updatePagination() {
	index := m.Index()
	availHeight := m.height

	if m.showTitle || m.showFilter && m.filteringEnabled {
		availHeight -= 2
	}
	if m.showStatusBar {
		availHeight -= 2
	}
	if m.showPagination {
		availHeight -= 2
	}
	if m.showHelp {
		availHeight -= 2
	}

	m.Paginator.PerPage = max(1, availHeight/(m.delegate.Height()+m.delegate.Spacing()))

	if pages := len(m.VisibleItems()); pages < 1 {
		m.Paginator.SetTotalPages(1)
	} else {
		m.Paginator.SetTotalPages(pages)
	}

	// Restore index
	m.Paginator.Page = index / m.Paginator.PerPage
	m.cursor = index % m.Paginator.PerPage

	// Make sure the page stays in bounds
	m.Paginator.Page = min(m.Paginator.Page, m.Paginator.TotalPages-1)
}

func (m *Model[I]) hideStatusMessage() {
	m.statusMessage = ""
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}

// Update is the Tea update loop.
func (m *Model[I]) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if key.Matches(msg, m.KeyMap.ForceQuit) {
			f(tea.Quit)
			return
		}

	case MsgFilterMatches[I]:
		m.filteredItems = filteredItems[I](msg)
		return

	case spinner.MsgTick:
		m.spinner.Update(msg, nil)
		// if m.showSpinner {
		// 	cmds = append(cmds, cmd...)
		// }

	case msgStatusMessageTimeout:
		m.hideStatusMessage()
	}

	if m.filterState == Filtering {
		m.handleFiltering(msg, f)
		return
	}

	f(m.handleBrowsing(msg)...)
}

// Updates for when a user is browsing the list.
func (m *Model[I]) handleBrowsing(msg tea.Msg) []tea.Cmd {
	numItems := len(m.VisibleItems())

	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch {
		// Note: we match clear filter before quit because, by default, they're
		// both mapped to escape.
		case key.Matches(msg, m.KeyMap.ClearFilter):
			m.resetFiltering()

		case key.Matches(msg, m.KeyMap.Quit):
			return []tea.Cmd{tea.Quit}

		case key.Matches(msg, m.KeyMap.CursorUp):
			m.CursorUp()

		case key.Matches(msg, m.KeyMap.CursorDown):
			m.CursorDown()

		case key.Matches(msg, m.KeyMap.PrevPage):
			m.Paginator.PrevPage()

		case key.Matches(msg, m.KeyMap.NextPage):
			m.Paginator.NextPage()

		case key.Matches(msg, m.KeyMap.GoToStart):
			m.Paginator.Page = 0
			m.cursor = 0

		case key.Matches(msg, m.KeyMap.GoToEnd):
			m.Paginator.Page = m.Paginator.TotalPages - 1
			m.cursor = m.Paginator.ItemsOnPage(numItems) - 1

		case key.Matches(msg, m.KeyMap.Filter):
			m.hideStatusMessage()
			if m.FilterInput.Value() == "" {
				// Populate filter with all items only if the filter is empty.
				m.filteredItems = m.itemsAsFilterItems()
			}
			m.Paginator.Page = 0
			m.cursor = 0
			m.filterState = Filtering
			m.FilterInput.CursorEnd()
			m.FilterInput.Focus()
			m.updateKeybindings()
			return []tea.Cmd{textinput.Blink}

		case key.Matches(msg, m.KeyMap.ShowFullHelp) || key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
			m.updatePagination()
		}
	}

	cmds := m.delegate.Update(msg, m)

	// Keep the index in bounds when paginating
	itemsOnPage := m.Paginator.ItemsOnPage(len(m.VisibleItems()))
	if m.cursor > itemsOnPage-1 {
		m.cursor = max(0, itemsOnPage-1)
	}

	return cmds
}

// Updates for when a user is in the filter editing interface.
func (m *Model[I]) handleFiltering(msg tea.Msg, f func(...tea.Cmd)) {
	// Handle keys
	if msg, ok := msg.(tea.MsgKey); ok {
		switch {
		case key.Matches(msg, m.KeyMap.CancelWhileFiltering):
			m.resetFiltering()
			m.KeyMap.Filter.SetEnabled(true)
			m.KeyMap.ClearFilter.SetEnabled(false)

		case key.Matches(msg, m.KeyMap.AcceptWhileFiltering):
			m.hideStatusMessage()

			if len(m.items) == 0 {
				break
			}

			h := m.VisibleItems()

			// If we've filtered down to nothing, clear the filter
			if len(h) == 0 {
				m.resetFiltering()
				break
			}

			m.FilterInput.Blur()
			m.filterState = FilterApplied
			m.updateKeybindings()

			if m.FilterInput.Value() == "" {
				m.resetFiltering()
			}
		}
	}

	// Update the filter text input component
	oldValue := m.FilterInput.Value()
	m.FilterInput.Update(msg, f)
	filterChanged := oldValue != m.FilterInput.Value()

	// If the filtering input has changed, request updated filtering
	if filterChanged {
		f(cmdFilterItems(*m))
		m.KeyMap.AcceptWhileFiltering.SetEnabled(m.FilterInput.Value() != "")
	}

	// Update pagination
	m.updatePagination()

	return
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m *Model[I]) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
	}

	filtering := m.filterState == Filtering

	// If the delegate implements the help.KeyMap interface add the short help
	// items to the short help after the cursor movement keys.
	if !filtering {
		if b, ok := m.delegate.(help.KeyMap); ok {
			kb = append(kb, b.ShortHelp()...)
		}
	}

	kb = append(kb,
		m.KeyMap.Filter,
		m.KeyMap.ClearFilter,
		m.KeyMap.AcceptWhileFiltering,
		m.KeyMap.CancelWhileFiltering,
	)

	if !filtering && m.AdditionalShortHelpKeys != nil {
		kb = append(kb, m.AdditionalShortHelpKeys()...)
	}

	return append(kb,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	)
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m *Model[I]) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.NextPage,
		m.KeyMap.PrevPage,
		m.KeyMap.GoToStart,
		m.KeyMap.GoToEnd,
	}}

	filtering := m.filterState == Filtering

	// If the delegate implements the help.KeyMap interface add full help
	// keybindings to a special section of the full help.
	if !filtering {
		if b, ok := m.delegate.(help.KeyMap); ok {
			kb = append(kb, b.FullHelp()...)
		}
	}

	listLevelBindings := []key.Binding{
		m.KeyMap.Filter,
		m.KeyMap.ClearFilter,
		m.KeyMap.AcceptWhileFiltering,
		m.KeyMap.CancelWhileFiltering,
	}

	if !filtering && m.AdditionalFullHelpKeys != nil {
		listLevelBindings = append(listLevelBindings, m.AdditionalFullHelpKeys()...)
	}

	return append(kb,
		listLevelBindings,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}

// View renders the component.
func (m *Model[I]) View(vb tea.Viewbox) {
	if m.showHelp {
		m.helpView(vb.PaddingTop(vb.Height - 1).PaddingLeft(2))
		vb = vb.Padding(tea.PaddingOptions{Bottom: 2})
	}

	if m.showPagination {
		m.paginationView(vb.PaddingTop(vb.Height - 1).PaddingLeft(2))
		vb = vb.Padding(tea.PaddingOptions{Bottom: 2})
	}

	if m.showTitle || m.showFilter && m.filteringEnabled {
		m.titleView(vb.PaddingLeft(2))
		vb = vb.PaddingTop(2)
	}

	if m.showStatusBar {
		m.statusView(vb.PaddingLeft(2))
		vb = vb.PaddingTop(2)
	}

	m.populatedView(vb)
}

func (m *Model[I]) titleView(vb tea.Viewbox) {
	vb = vb.Styled(m.Styles.TitleBar)

	// We need to account for the size of the spinner, even if we don't
	// render it, to reserve some space for it should we turn it on later.
	// spinnerView := m.spinnerView()
	// spinnerWidth := styles.Width(spinnerView)
	// spinnerLeftGap := " "
	// spinnerOnLeft := /*titleBarStyle.GetPaddingLeft() >= spinnerWidth+styles.Width(spinnerLeftGap) &&*/ m.showSpinner

	// If the filter's showing, draw that. Otherwise draw the title.
	if m.showFilter && m.filterState == Filtering {
		m.FilterInput.View(vb)
	} else if m.showTitle {
		// if m.showSpinner && spinnerOnLeft {
		// 	vb = vb.WriteLine(spinnerView + spinnerLeftGap)
		// 	// titleBarGap := titleBarStyle.GetPaddingLeft()
		// 	// titleBarStyle = titleBarStyle.PaddingLeft(titleBarGap - spinnerWidth - styles.Width(spinnerLeftGap))
		// }

		vb.Styled(m.Styles.Title).WriteLineX(" ").WriteLineX(m.Title).WriteLineX(" ")

		// // Status message
		// if m.filterState != Filtering {
		// 	vb.WriteLine(truncate.StringWithTail("  "+m.statusMessage, uint(m.width-spinnerWidth), ellipsis))
		// }
	}

	// // Spinner
	// if m.showSpinner && !spinnerOnLeft {
	// 	// Place spinner on the right
	// 	// availSpace := m.width - styles.Width(m.Styles.TitleBar.Render(view))
	// 	// if availSpace > spinnerWidth {
	// 	// 	x = vb.WriteLine(0, x, strings.Repeat(" ", availSpace-spinnerWidth)+spinnerView)
	// 	// }
	// }
}

func (m *Model[I]) statusView(vb tea.Viewbox) {
	vb = vb.Styled(m.Styles.StatusBar)

	visibleItems := len(m.VisibleItems())

	itemName := m.itemNameSingular
	if visibleItems != 1 {
		itemName = m.itemNamePlural
	}

	itemsDisplay := fmt.Sprintf("%d %s", visibleItems, itemName)

	switch {
	case m.filterState == Filtering:
		// Filter results
		if visibleItems == 0 {
			vb = vb.Styled(m.Styles.StatusEmpty).WriteLineX("Nothing matched")
		} else {
			vb = vb.WriteLineX(itemsDisplay)
		}
	case len(m.items) == 0:
		// Not filtering: no items.
		vb = vb.Styled(m.Styles.StatusEmpty).WriteLineX("No " + m.itemNamePlural)
	default:
		// Normal
		filtered := m.FilterState() == FilterApplied

		status := ""
		if filtered {
			status = fmt.Sprintf("“%s” ", truncate.StringWithTail(strings.TrimSpace(m.FilterInput.Value()), 10, "…"))
		}
		status += itemsDisplay

		vb = vb.WriteLineX(status)
	}

	totalItems := len(m.items)
	if numFiltered := totalItems - visibleItems; numFiltered > 0 {
		vb = vb.Styled(m.Styles.DividerDot).WriteLineX(string([]rune{' ', bullet, ' '}))
		vb.Styled(m.Styles.StatusBarFilterCount).WriteLine(fmt.Sprintf("%d filtered", numFiltered))
	}
}

func (m *Model[I]) paginationView(vb tea.Viewbox) {
	if m.Paginator.TotalPages < 2 { //nolint:gomnd
		return
	}

	style := m.Styles.PaginationStyle
	// If the dot pagination is wider than the width of the window
	// use the arabic paginator.
	if vb.Width > 0 && m.Paginator.TotalPages > vb.Width {
		m.Paginator.Type = paginator.Arabic
		style = m.Styles.ArabicPagination
	}

	m.Paginator.View(vb.Styled(style))
}

func (m *Model[I]) populatedView(vb tea.Viewbox) {
	items := m.VisibleItems()

	// Empty states
	if len(items) == 0 {
		if m.filterState != Filtering {
			vb.Styled(m.Styles.NoItems).WriteLine("No " + m.itemNamePlural + ".")
		}

		return
	}

	start, end := m.Paginator.GetSliceBounds(len(items))
	docs := items[start:end]

	for i, item := range docs {
		m.delegate.Render(vb, m, i+start, item)
		if i < len(docs)-1 {
			vb = vb.PaddingTop(m.delegate.Spacing() + m.delegate.Height())
		}
	}
}

func (m *Model[I]) helpView(vb tea.Viewbox) {
	m.Help.View(vb, m)
}

func (m *Model[I]) spinnerView() string {
	m.spinner.View(tea.Viewbox{})
	return ""
}

func cmdFilterItems[I Item](m Model[I]) tea.Cmd {
	return func() tea.Msg {
		if m.FilterInput.Value() == "" || m.filterState == Unfiltered {
			return MsgFilterMatches[I](m.itemsAsFilterItems()) // return nothing
		}

		targets := []string{}
		items := m.items

		for _, t := range items {
			targets = append(targets, t.FilterValue())
		}

		filterMatches := []filteredItem[I]{}
		for _, r := range m.Filter(m.FilterInput.Value(), targets) {
			filterMatches = append(filterMatches, filteredItem[I]{
				item:    items[r.Index],
				matches: r.MatchedIndexes,
			})
		}

		return MsgFilterMatches[I](filterMatches)
	}
}

func insertItemIntoSlice[I Item](items []I, item I, index int) []I {
	if items == nil {
		return []I{item}
	}

	if index >= len(items) {
		return append(items, item)
	}

	index = max(0, index)

	items = append(items, item)
	for i := len(items) - 1; i > index; i-- {
		items[i], items[i-1] = items[i-1], items[i]
	}
	return items
}

// Remove an item from a slice of items at the given index. This runs in O(n).
func removeItemFromSlice[I any](i []I, index int) []I {
	if index >= len(i) {
		return i // noop
	}

	copy(i[index:], i[index+1:])
	return i[:len(i)-1]
}

func countEnabledBindings(groups [][]key.Binding) int {
	agg := 0
	for _, group := range groups {
		for _, kb := range group {
			if kb.Enabled() {
				agg++
			}
		}
	}
	return agg
}
