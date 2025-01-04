// Package paginator provides a Tea package for calculating pagination
// and rendering pagination info. Note that this package does not render actual
// pages: it's purely for handling keystrokes related to pagination, and
// rendering pagination status.
package pager

import (
	"fmt"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/headless/pager"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
)

// Type specifies the way we render pagination.
type Type func(m *Model, vb tea.Viewbox)

// Pagination rendering options.
var (
	Arabic = (*Model).arabicView
	Dots   = (*Model).dotsView
)

// KeyMap is the key bindings for different actions within the paginator.
type KeyMap struct {
	PrevPage key.Binding
	NextPage key.Binding
}

// DefaultKeyMap is the default set of key bindings for navigating and acting
// upon the paginator.
var DefaultKeyMap = KeyMap{
	PrevPage: key.Binding{Keys: []string{"pgup", "left", "h"}},
	NextPage: key.Binding{Keys: []string{"pgdown", "right", "l"}},
}

// Model is the Tea model for this user interface.
type Model struct {
	Pager *pager.Pager

	Type             Type // Type configures how the pagination is rendered (Arabic, Dots).
	ActiveDot        rune // ActiveDot is used to mark the current page under the Dots display type.
	ActiveDotStyle   styles.Style
	InactiveDot      rune // InactiveDot is used to mark inactive pages under the Dots display type.
	InactiveDotStyle styles.Style
	ArabicFormat     string // ArabicFormat is the printf-style format to use for the Arabic display type.

	// KeyMap encodes the keybindings recognized by the widget.
	KeyMap KeyMap
}

// SetTotalPages is a helper function for calculating the total number of pages
// from a given number of items. Its use is optional since this pager can be
// used for other things beyond navigating sets. Note that it both returns the
// number of total pages and alters the model.
func (m *Model) SetTotalPages(items int) int {
	m.Pager.TotalSet(items)
	return m.Pager.Pages()
}

// ItemsOnPage is a helper function for returning the number of items on the
// current page given the total number of items passed as an argument.
func (m *Model) ItemsOnPage(totalItems int) int {
	if totalItems < 1 {
		return 0
	}
	start, end := m.GetSliceBounds(totalItems)
	return end - start
}

// GetSliceBounds is a helper function for paginating slices. Pass the length
// of the slice you're rendering and you'll receive the start and end bounds
// corresponding to the pagination. For example:
//
//	bunchOfStuff := []stuff{...}
//	start, end := model.GetSliceBounds(len(bunchOfStuff))
//	sliceToRender := bunchOfStuff[start:end]
func (m *Model) GetSliceBounds(length int) (int, int) {
	start := m.Pager.Page() * m.Pager.PerPage()
	end := min(start+m.Pager.PerPage(), length)
	return start, end
}

// PrevPage is a helper function for navigating one page backward. It will not
// page beyond the first page (i.e. page 0).
func (m *Model) PrevPage() {
	m.Pager.PagePrev()
}

// NextPage is a helper function for navigating one page forward. It will not
// page beyond the last page (i.e. totalPages - 1).
func (m *Model) NextPage() {
	m.NextPage()
}

// OnLastPage returns whether or not we're on the last page.
func (m *Model) OnLastPage() bool {
	return m.Pager.Page() == m.Pager.Pages()-1
}

// New creates a new model with defaults.
func New() Model {
	return Model{
		Type:         Arabic,
		Pager:        pager.New(0, 1),
		KeyMap:       DefaultKeyMap,
		ActiveDot:    '•',
		InactiveDot:  '○',
		ArabicFormat: "%d/%d",
	}
}

// Update is the Tea update function which binds keystrokes to pagination.
func (m *Model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.KeyMap.NextPage):
			m.NextPage()
		case key.Matches(msg, m.KeyMap.PrevPage):
			m.PrevPage()
		}
	}

	return nil
}

func (m *Model) dotsView(vb tea.Viewbox) {
	for i := 0; i < m.Pager.Total(); i++ {
		vb.Styled(m.InactiveDotStyle).Set(0, i, m.InactiveDot)
	}
	vb.Styled(m.ActiveDotStyle).Set(0, m.Pager.Page(), m.ActiveDot)
}

func (m *Model) arabicView(vb tea.Viewbox) {
	vb.WriteLine(fmt.Sprintf(m.ArabicFormat, m.Pager.Page()+1, m.Pager.Total()))
}

// View renders the pagination.
func (m *Model) View(vb tea.Viewbox) {
	m.Type(m, vb)
}
