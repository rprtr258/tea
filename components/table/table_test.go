package table

import (
	"testing"

	"github.com/rprtr258/assert"
)

func TestFromValues(t *testing.T) {
	table := New(WithColumns([]Column{{Title: "Foo"}, {Title: "Bar"}}))
	table.SetRows([][]string{
		{"foo1", "bar1"},
		{"foo2", "bar2"},
		{"foo3", "bar3"},
	}...)

	assert.Equal(t, [][]string{
		{"foo1", "bar1"},
		{"foo2", "bar2"},
		{"foo3", "bar3"},
	}, table.rows)
}

func TestFromValuesWithTabSeparator(t *testing.T) {
	table := New(WithColumns([]Column{{Title: "Foo"}, {Title: "Bar"}}))
	table.SetRows([][]string{
		{"foo1.", "bar1"},
		{"foo,bar,baz", "bar,2"},
	}...)

	assert.Equal(t, [][]string{
		{"foo1.", "bar1"},
		{"foo,bar,baz", "bar,2"},
	}, table.rows)
}
