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

type item struct{}

type model struct {
	list list.Model[list.DefaultItem[item]]
}

func (m *model) Init(tea.Context[*model]) {}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if msg.String() == "ctrl+c" {
			c.Dispatch(tea.Quit)
			return
		}
	case tea.MsgWindowSize:
		m.list.SetWidth(
			msg.Width - appPadding.Left - appPadding.Right,
		)
	}

	ctxList := tea.Of(c, func(m *model) *list.Model[list.DefaultItem[item]] { return &m.list })
	m.list.Update(ctxList, msg)
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
	}, func(title, description string) list.DefaultItem[item] {
		return list.DefaultItem[item]{
			Title:       title,
			Description: description,
		}
	})

	m := &model{
		list: list.New(items, list.NewDefaultDelegate[item](nil, nil, nil), func(di list.DefaultItem[item]) string {
			return di.Title
		}),
	}
	m.list.Title = "My Fave Things"

	_, err := tea.
		NewProgram2(ctx, m).
		WithAltScreen().
		Run()
	return err
}
