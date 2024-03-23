package toggle

type Toggle struct {
	State bool
}

func NewInitial(initialValue bool) *Toggle {
	return &Toggle{
		State: initialValue,
	}
}

func New() *Toggle {
	return NewInitial(false)
}

func (t *Toggle) Toggle() {
	t.State = !t.State
}
