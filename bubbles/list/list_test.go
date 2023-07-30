package list

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/tea"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

var _ ItemDelegate[item] = itemDelegate{}

func (d itemDelegate) Height() int                                { return 1 }
func (d itemDelegate) Spacing() int                               { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *Model[item]) []tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m *Model[item], index int, i item) {
	fmt.Fprint(w, m.Styles.TitleBar.Render(fmt.Sprintf("%d. %s", index+1, i)))
}

func TestStatusBarItemName(t *testing.T) {
	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	assert.Contains(t, list.statusView(), "2 items")

	list.SetItems([]item{item("foo")})
	assert.Contains(t, list.statusView(), "1 item")
}

func TestStatusBarWithoutItems(t *testing.T) {
	list := New[item]([]item{}, itemDelegate{}, 10, 10)
	assert.Contains(t, list.statusView(), "No items")
}

func TestCustomStatusBarItemName(t *testing.T) {
	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	list.SetStatusBarItemName("connection", "connections")
	assert.Contains(t, list.statusView(), "2 connections")

	list.SetItems([]item{item("foo")})
	assert.Contains(t, list.statusView(), "1 connection")

	list.SetItems([]item{})
	assert.Contains(t, list.statusView(), "No connections")
}
