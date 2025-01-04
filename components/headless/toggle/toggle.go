package toggle

type Toggle struct {
	state bool
}

func New() *Toggle {
	return &Toggle{
		state: false,
	}
}

func (t *Toggle) Value() bool { return t.state }

func (t *Toggle) Toggle()        { t.state = !t.state }
func (t *Toggle) On()            { t.state = true }
func (t *Toggle) Off()           { t.state = false }
func (t *Toggle) Set(value bool) { t.state = value }
