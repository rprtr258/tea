package list

import "iter"

type List[T any] struct {
	items    []T
	filter   func(T) bool
	selected int
}

func New[T any](items []T) *List[T] {
	return &List[T]{items, func(T) bool { return true }, -1}
}

func (l *List[T]) Items() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, item := range l.items {
			if l.filter(item) && !yield(i, item) {
				return
			}
		}
	}
}
func (l *List[T]) total() int {
	total := 0
	for range l.Items() {
		total++
	}
	return total
}

func (l *List[T]) Selected() (T, bool) {
	if l.selected < 0 || l.selected >= l.total() {
		return *new(T), false
	}
	return l.items[l.selected], true
}
func (l *List[T]) SelectedIndex() int { return l.selected }
func (l *List[T]) Select(index int)   { l.selected = min(max(0, index), l.total()-1) }
func (l *List[T]) SelectNext()        { l.Select(l.SelectedIndex() + 1) }
func (l *List[T]) SelectPrev()        { l.Select(l.SelectedIndex() - 1) }
func (l *List[T]) SelectFirst()       { l.Select(0) }
func (l *List[T]) SelectLast()        { l.Select(l.total() - 1) }

func (l *List[T]) FilterSet(filter func(T) bool) {
	l.filter = filter
	for l.selected >= 0 && !filter(l.items[l.selected]) {
		l.selected--
	}
}
func (l *List[T]) FilterReset() { l.FilterSet(func(T) bool { return true }) }
