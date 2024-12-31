package list

import (
	"fmt"

	"github.com/rprtr258/tea"
)

type item string

func (i item) FilterValue() string { return "" }

var itemDelegate = ItemDelegate[item]{
	func(vb tea.Viewbox, m *Model[item], index int, i item) {
		vb.WriteLine(m.Styles.TitleBar.Render(fmt.Sprintf("%d. %s", index+1, i)))
	},
	1,
	0,
	func(_ tea.Msg, _ *Model[item]) []tea.Cmd { return nil },
}

// func TestStatusBarItemName(t *testing.T) {
// 	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
// 	assert.True(t, strings.Contains(list.statusView(), "2 items"))

// 	list.SetItems([]item{item("foo")})
// 	assert.True(t, strings.Contains(list.statusView(), "1 item"))
// }

// func TestStatusBarWithoutItems(t *testing.T) {
// 	list := New[item]([]item{}, itemDelegate{}, 10, 10)
// 	assert.True(t, strings.Contains(list.statusView(), "No items"))
// }

// func TestCustomStatusBarItemName(t *testing.T) {
// 	list := New[item]([]item{item("foo"), item("bar")}, itemDelegate{}, 10, 10)
// 	list.SetStatusBarItemName("connection", "connections")
// 	assert.True(t, strings.Contains(list.statusView(), "2 connections"))

// 	list.SetItems([]item{item("foo")})
// 	assert.True(t, strings.Contains(list.statusView(), "1 connection"))

// 	list.SetItems([]item{})
// 	assert.True(t, strings.Contains(list.statusView(), "No connections"))
// }
