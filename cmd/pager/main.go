package pager

// An example program demonstrating pager component.

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/box"
	"github.com/rprtr258/tea/components/viewport"
)

type model struct {
	viewport viewport.Model
	lines    []string
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			f(tea.Quit)
			return
		}
	case tea.MsgWindowSize:
		const headerHeight = 3
		const footerHeight = 3
		m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
	}

	m.viewport.Update(msg)
	if m.viewport.YOffset > len(m.lines)-m.viewport.Height {
		m.viewport.YOffset = max(0, len(m.lines)-m.viewport.Height)
	}
}

func (m *model) headerView(vb tea.Viewbox) {
	const title = "Mr. Pager"
	box.Box(
		vb.MaxWidth(len(title)+2+2),
		func(vb tea.Viewbox) {
			vb.PaddingLeft(1).WriteLine(title)
		},
		box.RoundedBorder,
		box.BorderMaskAll,
		box.Colors(nil),
		box.Colors(nil),
	)
	vb = vb.PaddingLeft(1 + 1 + len(title) + 1).Row(1)
	vb.Set(0, 0, '├')
	vb.PaddingLeft(1).Fill('─')
}

func (m *model) footerView(vb tea.Viewbox) {
	info := fmt.Sprintf("%3.f%%", float64(m.viewport.YOffset)/float64(len(m.lines)-m.viewport.Height)*100)
	box.Box(
		vb.Sub(tea.Rectangle{
			Left:   vb.Width - 2 - len(info) - 2,
			Height: 3,
			Width:  2 + len(info) + 2,
		}),
		func(vb tea.Viewbox) {
			vb.PaddingLeft(1).WriteLine(info)
		},
		box.RoundedBorder,
		box.BorderMaskAll,
		box.Colors(nil),
		box.Colors(nil),
	)
	vb = vb.Padding(tea.PaddingOptions{Right: 1 + 1 + len(info) + 1}).Row(1)
	vb.Set(0, vb.Width-1, '┤')
	vb = vb.Padding(tea.PaddingOptions{Right: 1})
	for i := 0; i < vb.Width; i++ {
		vb.Set(0, i, '─')
	}
}

func (m *model) View(vb tea.Viewbox) {
	vbHeader, vbViewport, vbFooter := vb.SplitY3(tea.Fixed(3), tea.Flex(1), tea.Fixed(3)) // TODO: second is auto

	m.headerView(vbHeader)
	m.viewport.View(
		vbViewport,
		func(vb tea.Viewbox, i int) {
			if i >= len(m.lines) {
				return
			}
			vb.WriteLine(m.lines[i])
		})
	m.footerView(vbFooter)
}

func Main(ctx context.Context) error {
	// Load some text for our viewport
	content, err := os.ReadFile("./cmd/pager/artichoke.md")
	if err != nil {
		return fmt.Errorf("load file: %w", err)
	}

	_, err = tea.NewProgram(ctx, &model{
		lines: strings.Split(string(content), "\n"),
	}).
		WithAltScreen().       // use the full size of the terminal in its "alternate screen buffer"
		WithMouseCellMotion(). // turn on mouse support so we can track the mouse wheel
		Run()
	return err
}
