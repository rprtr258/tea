package list_default //nolint:revive,stylecheck

import (
	"context"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/list"
)

var appPadding = tea.PaddingOptions{
	Top:    1,
	Left:   2,
	Bottom: 1,
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model[item]
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if msg.String() == "ctrl+c" {
			f(tea.Quit)
			return
		}
	case tea.MsgWindowSize:
		m.list.SetSize(
			msg.Width-appPadding.Left-appPadding.Right,
			msg.Height-appPadding.Top-appPadding.Bottom,
		)
	}

	m.list.Update(msg, f)
}

func (m *model) View(vb tea.Viewbox) {
	m.list.View(vb.Padding(appPadding))
}

func Main(ctx context.Context) error {
	items := fun.MapToSlice(map[string]string{
		"Raspberry Pi’s":       "I have ’em all over my house",
		"Nutella":              "It's good on toast",
		"Bitter melon":         "It cools you down",
		"Nice socks":           "And by that I mean socks without holes",
		"Eight hours of sleep": "I had this once",
		"Cats":                 "Usually",
		"Plantasia, the album": "My plants love it too",
		"Pour over coffee":     "It takes forever to make though",
		"VR":                   "Virtual reality...what is there to say?",
		"Noguchi Lamps":        "Such pleasing organic forms",
		"Linux":                "Pretty much the best OS",
		"Business school":      "Just kidding",
		"Pottery":              "Wet clay is a great feeling",
		"Shampoo":              "Nothing like clean hair",
		"Table tennis":         "It’s surprisingly exhausting",
		"Milk crates":          "Great for packing in your extra stuff",
		"Afternoon tea":        "Especially the tea sandwich part",
		"Stickers":             "The thicker the vinyl the better",
		"20° Weather":          "Celsius, not Fahrenheit",
		"Warm light":           "Like around 2700 Kelvin",
		"The vernal equinox":   "The autumnal equinox is pretty good too",
		"Gaffer’s tape":        "Basically sticky fabric",
		"Terrycloth":           "In other words, towel fabric",
	}, func(title, desc string) item {
		return item{
			title: title,
			desc:  desc,
		}
	})

	m := &model{
		list: list.New[item](items, list.NewDefaultDelegate[item](), 0, 0),
	}
	m.list.Title = "My Fave Things"

	_, err := tea.
		NewProgram(ctx, m).
		WithAltScreen().
		Run()
	return err
}
