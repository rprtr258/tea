package list

import (
	"strings"

	"github.com/muesli/reflow/truncate"
	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/box"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
)

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle styles.Style
	NormalDesc  styles.Style

	// The selected item state.
	SelectedTitle styles.Style
	SelectedDesc  styles.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle styles.Style
	DimmedDesc  styles.Style

	// Characters matching the current filter, if any.
	FilterMatch styles.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
var NewDefaultItemStyles = func() DefaultItemStyles {
	s := DefaultItemStyles{
		NormalTitle: styles.Style{}.Foreground(styles.FgAdaptiveColor("#1a1a1a", "#dddddd")),
		SelectedTitle: styles.Style{}.
			// BorderForeground(styles.FgAdaptiveColor("#F793FF", "#AD58B4")).
			Foreground(styles.FgAdaptiveColor("#EE6FF8", "#EE6FF8")),
		DimmedTitle: styles.Style{}.Foreground(styles.FgAdaptiveColor("#A49FA5", "#777777")),
		FilterMatch: styles.Style{}.Underline(),
	}
	s.NormalDesc = s.NormalTitle.Copy().Foreground(styles.FgAdaptiveColor("#A49FA5", "#777777"))
	s.SelectedDesc = s.SelectedTitle.Copy().Foreground(styles.FgAdaptiveColor("#F793FF", "#AD58B4"))
	s.DimmedDesc = s.DimmedTitle.Copy().Foreground(styles.FgAdaptiveColor("#C2B8C2", "#4D4D4D"))

	return s
}()

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
		Styles:          NewDefaultItemStyles,
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
	textwidth := uint(m.width /*- s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()*/)
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
		vb = vb.PaddingLeft(2)

		vb.Styled(s.DimmedTitle).WriteLine(title)

		if d.ShowDescription {
			vb = vb.PaddingTop(1)
			vb.Styled(s.DimmedDesc).WriteLine(desc)
		}
	case isSelected && m.FilterState() != Filtering:
		vb.Styled(s.SelectedTitle).Set(0, 0, box.NormalBorder.Left)
		if isFiltered {
			// Highlight matches
			// unmatched := s.SelectedTitle.Inline(true)
			// matched := unmatched.Copy().Inherit(s.FilterMatch)
			// title = styles.StyleRunes(title, matchedRunes, matched, unmatched)
			vb.PaddingLeft(2).WriteLine(title)
		} else {
			vb.PaddingLeft(2).Styled(s.SelectedTitle).WriteLine(title)
		}

		if d.ShowDescription {
			vb = vb.PaddingTop(1)
			vb.Styled(s.SelectedTitle).Set(0, 0, box.NormalBorder.Left)
			vb.PaddingLeft(2).Styled(s.SelectedDesc).WriteLine(desc)
		}
	default:
		vb = vb.PaddingLeft(2)

		if isFiltered {
			// Highlight matches
			// unmatched := s.NormalTitle.Inline(true)
			// matched := unmatched.Copy().Inherit(s.FilterMatch)
			// title = styles.StyleRunes(title, matchedRunes, matched, unmatched)
			vb.WriteLine(title)
		} else {
			vb.Styled(s.NormalTitle).WriteLine(title)
		}

		if d.ShowDescription {
			vb.PaddingTop(1).Styled(s.NormalDesc).WriteLine(desc)
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
