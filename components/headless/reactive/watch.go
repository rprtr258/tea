package reactive

import "time"

func Watch[T comparable](
	cell *Cell[T],
	f func(old, new T),
) {
	old := (*cell).Get()
	inner := *cell
	*cell = fcell[T]{
		get: func() T {
			return old
		},
		set: func(t T) {
			inner.Set(t)
			if new := inner.Get(); old != new {
				f(old, new)
				old = new
			}
		},
	}
}

func Whenever(
	cell *Cell[bool],
	f func(),
) {
	Watch(cell, func(old, new bool) {
		if !old && new {
			f()
		}
	})
	if (*cell).Get() {
		f()
	}
}

func WatchDebounced[T comparable](
	cell *Cell[T],
	f func(old, new T),
	debounce time.Duration,
) {
	old := (*cell).Get()
	cancel := func() {}
	Watch(cell, func(_, new T) {
		cancel()
		timer := time.AfterFunc(debounce, func() {
			f(old, new)
			cancel = func() {}
			old = new
		})
		cancel = func() { timer.Stop() }
	})
}
