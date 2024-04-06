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
	// selectedIndex is the index of the currently selected item. Must always
	// satisfy 0 <= selectedIndex < len(items)
	selectedIndex int
	// highlightedIndex is the index of the currently highlighted item plus one,
	// such that zero value is no highlight
	highlightedIndex int
}

type Options[T any] struct {
	items                   []Item[T]
	defaultSelectedIndex    int
	defaultHighlightedIndex int
	less                    func(a, b Item[T]) int
}

func (o *Options[T]) WithItem(label string, value T) *Options[T] {
	o.items = append(o.items, Item[T]{
		Label: label,
		Value: value,
	})
	return o
}

func (o *Options[T]) WithItems(items ...Item[T]) *Options[T] {
	o.items = items
	return o
}

// WithMapItems, no order guaranteed
func (o *Options[T]) WithMapItems(items map[string]T) *Options[T] {
	for k, v := range items {
		o.items = append(o.items, Item[T]{
			Label: k,
			Value: v,
		})
	}
	return o
}

func (o *Options[T]) WithDefaultSelectedIndex(index int) *Options[T] {
	o.defaultSelectedIndex = index
	return o
}

func (o *Options[T]) WithDefaultHighlightedIndex(index int) *Options[T] {
	o.defaultHighlightedIndex = index
	return o
}

func (o *Options[T]) WithSort(less func(a, b Item[T]) int) *Options[T] {
	o.less = less
	return o
}

func (o *Options[T]) Build() *Selector[T] {
	if len(o.items) == 0 {
		panic("no items in selector")
	}

	if o.defaultSelectedIndex < 0 || o.defaultSelectedIndex >= len(o.items) {
		panic("invalid default selected index")
	}

	if o.defaultHighlightedIndex < 0 || o.defaultHighlightedIndex > len(o.items) {
		panic("invalid default highlighted index")
	}

	if o.less == nil {
		slices.SortFunc(o.items, o.less)
	}

	return &Selector[T]{
		items:            o.items,
		selectedIndex:    o.defaultSelectedIndex,
		highlightedIndex: o.defaultHighlightedIndex,
	}
}

func (s *Selector[T]) ResetHighlight() {
	s.highlightedIndex = 0
}

func (s *Selector[T]) Highlighted() (Item[T], bool) {
	if s.highlightedIndex == 0 {
		return Item[T]{}, false
	}

	return s.items[s.highlightedIndex-1], true
}

// Select highlighted item. No change if no highlight.
func (s *Selector[T]) Select() {
	if s.highlightedIndex == 0 {
		return
	}

	s.selectedIndex = s.highlightedIndex - 1
}

func (s *Selector[T]) SelectedItem() Item[T] {
	return s.items[s.selectedIndex]
}

func (s *Selector[T]) Selected() T {
	return s.SelectedItem().Value
}

func (s *Selector[T]) Items() []Item[T] {
	return s.items
}

func (s *Selector[T]) HighlightPrev() {
	switch {
	case s.highlightedIndex == 0: // no highlight yet
		s.highlightedIndex = 1 // select first
	case s.highlightedIndex > 1: // not first item highlighted
		s.highlightedIndex-- // highlight previous
	}
}

func (s *Selector[T]) HighlightPrevWrapped() {
	switch {
	case s.highlightedIndex == 0: // no highlight yet
		s.highlightedIndex = 1 // select first
	case s.highlightedIndex > 1: // not first item highlighted
		s.highlightedIndex-- // highlight previous
	default: // first item highlighted
		s.highlightedIndex = len(s.items)
	}
}

func (s *Selector[T]) HighlightNext() {
	switch {
	case s.highlightedIndex == 0: // no highlight yet
		s.highlightedIndex = 1 // select first
	case s.highlightedIndex < len(s.items): // not last item highlighted
		s.highlightedIndex++ // highlight next
	}
}

func (s *Selector[T]) HighlightNextWrapped() {
	switch {
	case s.highlightedIndex == 0: // no highlight yet
		s.highlightedIndex = 1 // select first
	case s.highlightedIndex < len(s.items): // not last item highlighted
		s.highlightedIndex++ // highlight next
	default: // last item highlighted
		s.highlightedIndex = 1
	}
}

func (s *Selector[T]) HighlightFirst() {
	s.highlightedIndex = 1
}

func (s *Selector[T]) HighlightLast() {
	s.highlightedIndex = len(s.items)
}

// Filter by predicate, returns list of setters to according items.
// No filter highlight is supported for now.
func (s *Selector[T]) FilterBy(predicate func(Item[T]) bool) []Item[func()] {
	return fun.FilterMap[Item[func()]](
		func(item Item[T], i int) (Item[func()], bool) {
			return Item[func()]{
				Label: item.Label,
				Value: func() {
					s.selectedIndex = i
					s.highlightedIndex = 0
				},
			}, predicate(item)
		},
		s.items...,
	)
}

// Filter by prefix, returns list of setters to according items.
// No filter highlight is supported for now.
func (s *Selector[T]) Filter(prefix string) []Item[func()] {
	return s.FilterBy(func(item Item[T]) bool {
		return strings.HasPrefix(item.Label, prefix)
	})
}
