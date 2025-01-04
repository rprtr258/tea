package result

type Result[T any] struct {
	value     T
	err       error
	isLoading bool
}

func New[T any](
	fetch func() (T, error),
) (_ *Result[T], refresh func()) {
	res := &Result[T]{
		isLoading: true,
	}
	load := func() {
		if !res.isLoading {
			return
		}
		go func() {
			value, err := fetch()
			res.value = value
			res.err = err
			res.isLoading = false
		}()
	}
	load()
	return res, func() {
		res.isLoading = true
		load()
	}
}

func (r *Result[T]) Value() (T, bool) {
	return r.value, !r.isLoading && r.err == nil
}

func (r *Result[T]) FetchErr() error {
	return r.err
}

func (r *Result[T]) IsLoading() bool {
	return r.isLoading
}
