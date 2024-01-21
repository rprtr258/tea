package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/markdown"
	"github.com/rprtr258/tea/components/viewport"
	"github.com/rprtr258/tea/styles"
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

var helpStyle = styles.Style{}.Foreground(styles.FgColor("241")).Render

type model struct {
	viewport viewport.Model
}

func newExample() (*model, error) {
	const width = 78

	vp := viewport.New(width, 20)
	// BorderStyle(styles.RoundedBorder).
	// BorderForeground(styles.FgColor("62"))
	// PaddingRight(2)

	return &model{
		viewport: vp,
	}, nil
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			f(tea.Quit)
		default:
			m.viewport.Update(msg)
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	renderer, err := markdown.NewTermRenderer(
		markdown.WithAutoStyle(),
		markdown.WithWordWrap(vb.Width),
	)
	if err != nil {
		return
		// return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return
		// return nil, err
	}

	lines := strings.Split(str, "\n")
	m.viewport.View(vb, func(v tea.Viewbox, i int) {
		v.WriteLine(lines[i])
	})
	// TODO: right after viewport
	vb.PaddingTop(vb.Height - 1).WriteLine(helpStyle("  ↑/↓: Navigate • q: Quit"))
}

func Main(ctx context.Context) error {
	model, err := newExample()
	if err != nil {
		return fmt.Errorf("initialize model: %w", err)
	}

	_, err = tea.NewProgram(ctx, model).Run()
	return err
}
