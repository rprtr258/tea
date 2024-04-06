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

func (t *Toggle) Value() bool { return t.State }

func (t *Toggle) Toggle()        { t.State = !t.State }
func (t *Toggle) On()            { t.State = true }
func (t *Toggle) Off()           { t.State = false }
func (t *Toggle) Set(value bool) { t.State = value }
