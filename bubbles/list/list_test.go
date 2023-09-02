package list

import (
	"fmt"
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
func (d itemDelegate) Render(vb tea.Viewbox, m *Model[item], index int, i item) {
	vb.WriteLine(0, 0, m.Styles.TitleBar.Render(fmt.Sprintf("%d. %s", index+1, i)))
}

func TestStatusBarItemName(t *testing.T) {
	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	assert.Substring(t, list.statusView(), "2 items")

	list.SetItems([]item{item("foo")})
	assert.Substring(t, list.statusView(), "1 item")
}

func TestStatusBarWithoutItems(t *testing.T) {
	list := New[item]([]item{}, itemDelegate{}, 10, 10)
	assert.Substring(t, list.statusView(), "No items")
}

func TestCustomStatusBarItemName(t *testing.T) {
	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
	list.SetStatusBarItemName("connection", "connections")
	assert.Substring(t, list.statusView(), "2 connections")

	list.SetItems([]item{item("foo")})
	assert.Substring(t, list.statusView(), "1 connection")

	list.SetItems([]item{})
	assert.Substring(t, list.statusView(), "No connections")
}
