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
	content  string
	ready    bool
	viewport viewport.Model
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			f(tea.Quit)
			return
		}

	case tea.MsgWindowSize:
		const headerHeight = 3
		const footerHeight = 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport.
			// The initial dimensions come in quickly, though asynchronously,
			// which is why we wait for them here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetContent(strings.Split(m.content, "\n"))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	f(m.viewport.Update(msg)...)
}

func (m *model) headerView(vb tea.Viewbox) {
	title := "Mr. Pager"
	box.Box(
		vb.Sub(tea.Rectangle{
			Width: len(title) + 2 + 2,
		}),
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
	vb = vb.PaddingLeft(1)
	for i := 0; i < vb.Width; i++ {
		vb.Set(0, i, '─')
	}
}

func (m *model) footerView(vb tea.Viewbox) {
	info := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
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
	if !m.ready {
		vb.WriteLine("  Initializing...")
		return
	}

	m.headerView(vb.Sub(tea.Rectangle{
		Height: 3,
		Width:  vb.Width,
	}))
	vb = vb.PaddingTop(3)
	m.viewport.View(vb.Sub(tea.Rectangle{
		Height: m.viewport.Height,
		Width:  vb.Width,
	}))
	m.footerView(vb.PaddingTop(m.viewport.Height))
}

func Main(ctx context.Context) error {
	// Load some text for our viewport
	content, err := os.ReadFile("./cmd/pager/artichoke.md")
	if err != nil {
		return fmt.Errorf("load file: %w", err)
	}

	_, err = tea.NewProgram(ctx, &model{content: string(content)}).
		WithAltScreen().       // use the full size of the terminal in its "alternate screen buffer"
		WithMouseCellMotion(). // turn on mouse support so we can track the mouse wheel
		Run()
	return err
}
