package list_default

import (
	"context"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/list"
	"github.com/rprtr258/tea/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model[item]
}

func (m *model) Init() []tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if msg.String() == "ctrl+c" {
			return []tea.Cmd{tea.Quit}
		}
	case tea.MsgWindowSize:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	return m.list.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	r.Write(docStyle.Render(m.list.View()))
}

func Main(ctx context.Context) error {
	items := []item{
		{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		{title: "Nutella", desc: "It's good on toast"},
		{title: "Bitter melon", desc: "It cools you down"},
		{title: "Nice socks", desc: "And by that I mean socks without holes"},
		{title: "Eight hours of sleep", desc: "I had this once"},
		{title: "Cats", desc: "Usually"},
		{title: "Plantasia, the album", desc: "My plants love it too"},
		{title: "Pour over coffee", desc: "It takes forever to make though"},
		{title: "VR", desc: "Virtual reality...what is there to say?"},
		{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		{title: "Linux", desc: "Pretty much the best OS"},
		{title: "Business school", desc: "Just kidding"},
		{title: "Pottery", desc: "Wet clay is a great feeling"},
		{title: "Shampoo", desc: "Nothing like clean hair"},
		{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		{title: "Stickers", desc: "The thicker the vinyl the better"},
		{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		{title: "Warm light", desc: "Like around 2700 Kelvin"},
		{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := &model{
		list: list.New[item](items, list.NewDefaultDelegate[item](), 0, 0),
	}
	m.list.Title = "My Fave Things"

	_, err := tea.
		NewProgram(context.Background(), m).
		WithAltScreen().
		Run()
	return err
}
