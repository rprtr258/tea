package autocomplete

import (
	"iter"
	"strings"
)

type Autocomplete struct {
	items  []string
	prefix string
}

func New(items ...string) *Autocomplete {
	return &Autocomplete{
		items:  items,
		prefix: "",
	}
}

func (a *Autocomplete) SetPrefix(prefix string) {
	a.prefix = prefix
}

func (a *Autocomplete) Items() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, item := range a.items {
			if strings.HasPrefix(item, a.prefix) && !yield(item) {
				return
			}
		}
	}
}
