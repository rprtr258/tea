package selector

import (
	"testing"

	"github.com/rprtr258/assert"
)

func exampleSelector() *Selector[struct{}] {
	return New([]Item[struct{}]{
		{Label: "Apple", Value: struct{}{}},
		{Label: "Orange", Value: struct{}{}},
		{Label: "Banana", Value: struct{}{}},
	}, -1, -1, nil)
}

func TestShouldHandleDropdownOpenCloseState(t *testing.T) {
	s := exampleSelector()
	assert.False(t, s.IsOpen())
	s.Toggle()
	assert.True(t, s.IsOpen())
	s.Toggle()
	assert.False(t, s.IsOpen())
}

func TestSelectItem(t *testing.T) {
	s := exampleSelector()
	s.Toggle()
	s.HighlightNext()
	s.Select()
	item, ok := s.SelectedItem()
	assert.True(t, ok)
	assert.Equal(t, "Apple", item.Label)
}
