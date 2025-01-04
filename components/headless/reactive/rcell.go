package reactive

type RCell[T any] interface {
	Get() T
}

type constCell[T any] struct{ t T }

func Const[T any](t T) constCell[T] { return constCell[T]{t} }
func (c constCell[T]) Get() T       { return c.t }

type computed[T, R any] struct {
	inner RCell[T]
	from  func(T) R
}

func Computed[T, R any](inner RCell[T], f func(T) R) computed[T, R] { return computed[T, R]{inner, f} }
func (c computed[T, R]) Get() R                                     { return c.from(c.inner.Get()) }

type binaryOp[A, B, C any] struct {
	a Cell[A]
	b Cell[B]
	f func(A, B) C
}

func BinaryOp[A, B, C any](a Cell[A], b Cell[B], f func(A, B) C) binaryOp[A, B, C] {
	return binaryOp[A, B, C]{a, b, f}
}
func (c binaryOp[A, B, C]) Get() C {
	return c.f(c.a.Get(), c.b.Get())
}

type combine[T, R any] struct {
	cells Cell[[]T]
	f     func(...T) R
}

func Combine[T, R any](
	cells Cell[[]T],
	f func(...T) R,
) combine[T, R] {
	return combine[T, R]{cells, f}
}
func (c combine[T, R]) Get() R {
	return c.f(c.cells.Get()...)
}
