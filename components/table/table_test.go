package table

import (
	"testing"

	"github.com/rprtr258/assert"
)

func TestFromValues(t *testing.T) {
	table := New([]Column{{Title: "Foo"}, {Title: "Bar"}}, []int{0, 0},
		[][]string{
			{"foo1", "bar1"},
			{"foo2", "bar2"},
			{"foo3", "bar3"},
		}, 0, 0, DefaultStyles, DefaultKeyMap)

	assert.Equal(t, [][]string{
		{"foo1", "bar1"},
		{"foo2", "bar2"},
		{"foo3", "bar3"},
	}, table.Table.Rows())
}

func TestFromValuesWithTabSeparator(t *testing.T) {
	table := New([]Column{{Title: "Foo"}, {Title: "Bar"}}, []int{0, 0},
		[][]string{
			{"foo1.", "bar1"},
			{"foo,bar,baz", "bar,2"},
		}, 0, 0, DefaultStyles, DefaultKeyMap)

	assert.Equal(t, [][]string{
		{"foo1.", "bar1"},
		{"foo,bar,baz", "bar,2"},
	}, table.Table.Rows())
}
