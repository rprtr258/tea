package selector

import (
	"slices"
	"strings"

	"github.com/rprtr258/fun"
)

type Item[T any] struct {
	Label string
	Value T
}

type Selector[T any] struct {
	items []Item[T]
	// indexSelected is the index of the currently selected item. Must always
	// satisfy 0 <= indexSelected < len(items)
	// TODO: allow unselected
	indexSelected int
	// indexHighlighted is the index of the currently highlighted item plus one,
	// such that zero value is no highlight
	// TODO: allow unhighlighted
	indexHighlighted int
	isOpen           bool
}

func New[T any](
	items []Item[T],
	defaultSelectedIndex int,
	defaultHighlightedIndex int,
	less func(a, b Item[T]) int,
) *Selector[T] {
	if len(items) == 0 {
		panic("no items in selector")
	}

	if defaultSelectedIndex < -1 || defaultSelectedIndex >= len(items) {
		panic("invalid default selected index")
	}

	if defaultHighlightedIndex < -1 || defaultHighlightedIndex >= len(items) {
		panic("invalid default highlighted index")
	}

	if less != nil {
		slices.SortFunc(items, less)
	}

	return &Selector[T]{
		items:            items,
		indexSelected:    defaultSelectedIndex,
		indexHighlighted: defaultHighlightedIndex,
		isOpen:           false,
	}
}

func (s *Selector[T]) IsOpen() bool { return s.isOpen }
func (s *Selector[T]) Show()        { s.isOpen = true }
func (s *Selector[T]) Hide() {
	s.isOpen = false
	s.indexHighlighted = -1
}
func (s *Selector[T]) Toggle() {
	if s.IsOpen() {
		s.Hide()
	} else {
		s.Show()
	}
}

func (s *Selector[T]) ResetHighlight() {
	s.indexHighlighted = 0
}

func (s *Selector[T]) Highlighted() (Item[T], bool) {
	if s.indexHighlighted == 0 {
		return Item[T]{}, false
	}

	return s.items[s.indexHighlighted-1], true
}

// Select highlighted item. No change if no highlight.
func (s *Selector[T]) Select() {
	if s.indexHighlighted == -1 {
		return
	}

	s.indexSelected = s.indexHighlighted
}

func (s *Selector[T]) SelectedItem() (Item[T], bool) {
	if s.indexSelected == -1 {
		return Item[T]{}, false
	}

	return s.items[s.indexSelected], true
}

func (s *Selector[T]) Selected() (T, bool) {
	item, ok := s.SelectedItem()
	return item.Value, ok
}

func (s *Selector[T]) Items() []Item[T] {
	return s.items
}

func (s *Selector[T]) HighlightPrev() {
	s.Show()

	switch s.indexHighlighted {
	case -1: // no highlight yet
		s.indexHighlighted = 0 // select first
	case 0: // first item highlighted
		s.indexHighlighted = len(s.items) - 1
	default: // not first item highlighted
		s.indexHighlighted-- // highlight previous
	}
}

func (s *Selector[T]) HighlightNext() {
	s.Show()

	switch s.indexHighlighted {
	case -1: // no highlight yet
		s.indexHighlighted = 0 // select first
	case len(s.items) - 1: // last item highlighted
		s.indexHighlighted = 0
	default: // not last item highlighted
		s.indexHighlighted++ // highlight next
	}
}

// Prefix filter, returns list of setters to according items.
func Prefix[T any](prefix string) func(item Item[T]) bool {
	return func(item Item[T]) bool {
		return strings.HasPrefix(item.Label, prefix)
	}
}

// Filter by predicate, returns list of setters to according items.
// TODO: change semantics
func (s *Selector[T]) FilterBy(predicate func(Item[T]) bool) []Item[func()] {
	return fun.FilterMap[Item[func()]](
		func(item Item[T], i int) (Item[func()], bool) {
			return Item[func()]{
				Label: item.Label,
				Value: func() {
					s.indexSelected = i
					s.indexHighlighted = 0
				},
			}, predicate(item)
		},
		s.items...,
	)
}
