package history

type History[T any] struct {
	items   []T
	current int
}

func New[T any](init T) *History[T] {
	return &History[T]{
		items:   []T{init},
		current: 0,
	}
}

func (h *History[T]) Items() []T {
	return h.items
}

func (h *History[T]) Push(item T) {
	h.items = append(h.items[:h.current+1], item)
}

func (h *History[T]) BackCan() bool { return h.current > 0 }
func (h *History[T]) Back() {
	if h.BackCan() {
		h.current--
	}
}

func (h *History[T]) ForwardCan() bool { return h.current < len(h.items)-1 }
func (h *History[T]) Forward() {
	if h.ForwardCan() {
		h.current++
	}
}

func (h *History[T]) Current() T { return h.items[h.current] }
