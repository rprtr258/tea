package selector_multiple

import (
	"slices"

	"github.com/rprtr258/fun"
)

type Item[T any] struct {
	Label string
	Value T
}

type Selector[T any] struct {
	items []Item[T]
	// selectedIndexes are the indexes of the currently selected items.
	// All of them must satisfy 0 <= selectedIndex < len(items)
	selectedIndexes []int
	// selected - map of selectedIndexes
	selected map[int]struct{}
	// highlightedIndex is the index of the currently highlighted item plus one,
	// such that zero value is no highlight
	highlightedIndex int
}

type Options[T any] struct {
	items                   []Item[T]
	defaultSelectedIndex    int
	defaultHighlightedIndex int
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

func (o *Options[T]) WithDefaultHighlightedIndex(index int) *Options[T] {
	o.defaultHighlightedIndex = index
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

	return &Selector[T]{
		items:            o.items,
		selectedIndexes:  nil,
		selected:         map[int]struct{}{},
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
func (s *Selector[T]) ToggleHighlighted() {
	if s.highlightedIndex == 0 {
		return
	}

	newIndex := s.highlightedIndex - 1

	if _, ok := s.selected[newIndex]; ok {
		delete(s.selected, newIndex)
		_, i, _ := fun.Index(func(index int) bool {
			return index == newIndex
		}, s.selectedIndexes...)
		s.selectedIndexes = slices.Delete(s.selectedIndexes, i, i+1)
	} else {
		s.selected[newIndex] = struct{}{}
		s.selectedIndexes = append(s.selectedIndexes, newIndex)
	}
}

func (s *Selector[T]) AddHighlighted() {
	if s.highlightedIndex == 0 {
		return
	}

	newIndex := s.highlightedIndex - 1

	if _, ok := s.selected[newIndex]; ok {
		return
	}

	s.selected[newIndex] = struct{}{}
	s.selectedIndexes = append(s.selectedIndexes, newIndex)
}

func (s *Selector[T]) SelectedCount() int {
	return len(s.selectedIndexes)
}

// RemoveSelected removes selected item.
// i must satisfy 0 <= i < SelectedCount()
func (s *Selector[T]) RemoveSelected(i int) {
	if i < 0 || i >= s.SelectedCount() {
		panic("remove selected: invalid index")
	}

	if _, ok := s.selected[i]; ok {
		return
	}

	delete(s.selected, i)
	_, j, _ := fun.Index(func(index int) bool {
		return index == i
	}, s.selectedIndexes...)
	s.selectedIndexes = slices.Delete(s.selectedIndexes, j, j+1)
}

func (s *Selector[T]) SelectedItems() []Item[T] {
	return fun.Map[Item[T]](
		func(index int) Item[T] {
			return s.items[index]
		},
		s.selectedIndexes...,
	)
}

func (s *Selector[T]) Selected() []T {
	return fun.Map[T](
		func(index int) T {
			return s.items[index].Value
		},
		s.selectedIndexes...,
	)
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
