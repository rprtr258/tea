package menu

import (
	"context"
	"log"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/viewport"
	"github.com/rprtr258/tea/glamour"
	"github.com/rprtr258/tea/lipgloss"
)

const content = `
# Today’s Menu

## Appetizers

| Name        | Price | Notes                           |
| ---         | ---   | ---                             |
| Tsukemono   | $2    | Just an appetizer               |
| Tomato Soup | $4    | Made with San Marzano tomatoes  |
| Okonomiyaki | $4    | Takes a few minutes to make     |
| Curry       | $3    | We can add squash if you’d like |

## Seasonal Dishes

| Name                 | Price | Notes              |
| ---                  | ---   | ---                |
| Steamed bitter melon | $2    | Not so bitter      |
| Takoyaki             | $3    | Fun to eat         |
| Winter squash        | $3    | Today it's pumpkin |

## Desserts

| Name         | Price | Notes                 |
| ---          | ---   | ---                   |
| Dorayaki     | $4    | Looks good on rabbits |
| Banana Split | $5    | A classic             |
| Cream Puff   | $3    | Pretty creamy!        |

All our dishes are made in-house by Karen, our chef. Most of our ingredients
are from our garden or the fish market down the street.

Some famous people that have eaten here lately:

* [x] René Redzepi
* [x] David Chang
* [ ] Jiro Ono (maybe some day)

Bon appétit!
`

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type model struct {
	viewport viewport.Model
}

func newExample() (*model, error) {
	const width = 78

	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &model{
		viewport: vp,
	}, nil
}

func (m *model) Init() []tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return []tea.Cmd{tea.Quit}
		default:
			return m.viewport.Update(msg)
		}
	default:
		return nil
	}
}

func (m *model) View(r tea.Renderer) {
	r.Write(m.viewport.View() + m.helpView())
}

func (e *model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}

func Main() {
	model, err := newExample()
	if err != nil {
		log.Fatalln("Could not initialize Bubble Tea model:", err.Error())
	}

	if _, err := tea.NewProgram(context.Background(), model).Run(); err != nil {
		log.Fatalln("Bummer, there's been an error:", err.Error())
	}
}
