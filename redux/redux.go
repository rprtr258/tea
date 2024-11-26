package redux

import (
	"reflect"
)

// ReadCell read only cell
type ReadCell[T any] func() T

func (c ReadCell[T]) Get() T {
	return c()
}

func (c ReadCell[T]) Map(fn func(T) T) ReadCell[T] {
	return func() T {
		return fn(c.Get())
	}
}

func MapRO[T, U any](c ReadCell[T], fn func(T) U) ReadCell[U] {
	return func() U {
		return fn(c.Get())
	}
}

// Cell - if no args provided, returns stored cell
// if one arg provided, value is set to the arg, and returned
type Cell[T any] func(...T) T

func (c Cell[T]) Set(t T) T {
	return c(t)
}

func (c Cell[T]) Get() T {
	return c()
}

func (c Cell[T]) Update(fn func(T) T) T {
	return c.Set(fn(c.Get()))
}

func (c Cell[T]) Subscribe(callback func(old, new T)) Cell[T] {
	return func(t ...T) T {
		old := c.Get()
		if len(t) == 0 {
			return old
		}

		new := t[0]
		if reflect.DeepEqual(old, new) {
			return old
		}
		callback(old, new)
		return c.Set(new)
	}
}

func (c Cell[T]) Trace(name string) Cell[T] {
	return c.Subscribe(func(old, new T) {
		println(name, ":", old, "->", new)
	})
}

func (c Cell[T]) RO() ReadCell[T] {
	return func() T {
		return c.Get()
	}
}

func Stored[T any](init T) Cell[T] {
	return func(t ...T) T {
		if len(t) > 0 {
			init = t[0]
		}
		return init
	}
}

func Map[T, R any](
	inner Cell[T],
	from func(T) R,
	to func(R) T,
) Cell[R] {
	return func(t ...R) R {
		if len(t) == 0 {
			return from(inner.Get())
		}

		return from(inner.Set(to(t[0])))
	}
}

type storeInt struct {
	Value Cell[int]
}

func (s storeInt) Increment(delta int) {
	s.Value.Update(func(x int) int { return x + delta })
}

func (s storeInt) Decrement(delta int) {
	s.Value.Update(func(x int) int { return x - delta })
}

func (s storeInt) ViewPositive() ReadCell[int] {
	return s.Value.RO().Map(func(x int) int {
		return max(x, 0)
	})
}
