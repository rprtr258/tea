package reactive

type Cell[T any] interface {
	RCell[T]
	Set(T)
}

type fcell[T any] struct {
	get func() T
	set func(T)
}

func (c fcell[T]) Get() T  { return c.get() }
func (c fcell[T]) Set(t T) { c.set(t) }

type stored[T any] struct{ t T }

func Stored[T any](init T) *stored[T] { return &stored[T]{init} }
func (c *stored[T]) Get() T           { return c.t }
func (c *stored[T]) Set(t T)          { c.t = t }

type transient[T, R any] struct {
	inner Cell[T]
	from  func(T) R
	to    func(R) T
}

func (c *transient[T, R]) Get() R  { return c.from(c.inner.Get()) }
func (c *transient[T, R]) Set(t R) { c.inner.Set(c.to(t)) }
func Transient[T, R any](inner Cell[T], from func(T) R, to func(R) T) *transient[T, R] {
	return &transient[T, R]{inner, from, to}
}

type cached[T any] struct {
	inner Cell[T]
	value T
	valid bool
}

func Cached[T any](inner Cell[T]) *cached[T] {
	return &cached[T]{inner, *new(T), false}
}
func (c *cached[T]) Get() T {
	if !c.valid {
		c.value = c.inner.Get()
		c.valid = true
	}
	return c.value
}
func (c *cached[T]) Set(t T) {
	c.inner.Set(t)
	c.valid = false
}

func Sync[T any](a, b *Cell[T]) {
	innerA := *a
	innerB := *b
	set := func(t T) {
		innerA.Set(t)
		innerB.Set(t)
	}
	*a = fcell[T]{
		get: func() T { return innerA.Get() },
		set: set,
	}
	*b = fcell[T]{
		get: func() T { return innerB.Get() },
		set: set,
	}
}
