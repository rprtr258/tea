package list

import (
	"strings"

	"github.com/muesli/reflow/truncate"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/lipgloss"
)

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() DefaultItemStyles {
	s := DefaultItemStyles{
		NormalTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}),
		SelectedTitle: lipgloss.NewStyle().
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}),
		DimmedTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}),
		FilterMatch: lipgloss.NewStyle().
			Underline(true),
	}
	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})
	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	return s
}

// DefaultItem describes an items designed to work with DefaultDelegate.
type DefaultItem interface {
	Item
	Title() string
	Description() string
}

// DefaultDelegate is a standard delegate designed to work in lists. It's
// styled by DefaultItemStyles, which can be customized as you like.
//
// The description line can be hidden by setting Description to false, which
// renders the list as single-line-items. The spacing between items can be set
// with the SetSpacing method.
//
// Setting UpdateFunc is optional. If it's set it will be called when the
// ItemDelegate called, which is called when the list's Update function is
// invoked.
//
// Settings ShortHelpFunc and FullHelpFunc is optional. They can be set to
// include items in the list's default short and full help menus.
type DefaultDelegate[I DefaultItem] struct {
	ShowDescription bool
	Styles          DefaultItemStyles
	UpdateFunc      func(tea.Msg, *Model[I]) []tea.Cmd
	ShortHelpFunc   func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
	height          int
	spacing         int
}

// NewDefaultDelegate creates a new delegate with default styles.
func NewDefaultDelegate[I DefaultItem]() DefaultDelegate[I] {
	return DefaultDelegate[I]{
		ShowDescription: true,
		Styles:          NewDefaultItemStyles(),
		height:          2,
		spacing:         1,
	}
}

// SetHeight sets delegate's preferred height.
func (d *DefaultDelegate[I]) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d DefaultDelegate[I]) Height() int {
	return fun.IF(d.ShowDescription, d.height, 1)
}

// SetSpacing sets the delegate's spacing.
func (d *DefaultDelegate[I]) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d DefaultDelegate[I]) Spacing() int {
	return d.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d DefaultDelegate[I]) Update(msg tea.Msg, m *Model[I]) []tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

// Render prints an item.
func (d DefaultDelegate[I]) Render(vb tea.Viewbox, m *Model[I], index int, item I) {
	if m.width <= 0 {
		// short-circuit
		return
	}

	title := item.Title()
	desc := item.Description()

	s := &d.Styles
	// Prevent text from exceeding list width
	textwidth := uint(m.width - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.height-1 {
				break
			}
			lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	// Conditions
	isSelected := index == m.Index()
	emptyFilter := m.FilterState() == Filtering && m.FilterValue() == ""
	isFiltered := m.FilterState() == Filtering || m.FilterState() == FilterApplied

	// var matchedRunes []int
	// if isFiltered && index < len(m.filteredItems) {
	// 	// Get indices of matched characters
	// 	matchedRunes = m.MatchesForItem(index)
	// }

	switch {
	case emptyFilter:
		vb = vb.Padding(tea.PaddingOptions{Left: 2})

		vb.Styled(s.DimmedTitle).WriteLine(0, 0, title)

		if d.ShowDescription {
			vb.Styled(s.DimmedDesc).WriteLine(1, 0, desc)
		}
	case isSelected && m.FilterState() != Filtering:
		vb = vb.Padding(tea.PaddingOptions{Left: 1})

		vb.Styled(s.SelectedTitle).WriteLine(0, 0, lipgloss.NormalBorder.Left)

		if isFiltered {
			// Highlight matches
			// unmatched := s.SelectedTitle.Inline(true)
			// matched := unmatched.Copy().Inherit(s.FilterMatch)
			// title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			vb.WriteLine(0, 1, title)
		} else {
			vb.Styled(s.SelectedTitle).WriteLine(0, 1, title)
		}

		if d.ShowDescription {
			vb.Styled(s.SelectedTitle).WriteLine(1, 0, lipgloss.NormalBorder.Left)
			vb.Styled(s.SelectedDesc).WriteLine(1, 1, desc)
		}
	default:
		vb = vb.Padding(tea.PaddingOptions{Left: 2})

		if isFiltered {
			// Highlight matches
			// unmatched := s.NormalTitle.Inline(true)
			// matched := unmatched.Copy().Inherit(s.FilterMatch)
			// title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			vb.Styled(s.NormalTitle).WriteLine(0, 0, title)
		} else {
			vb.Styled(s.NormalTitle).WriteLine(0, 0, title)
		}

		if d.ShowDescription {
			vb.Styled(s.NormalDesc).WriteLine(1, 0, desc)
		}
	}
}

// ShortHelp returns the delegate's short help.
func (d DefaultDelegate[I]) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d DefaultDelegate[I]) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
