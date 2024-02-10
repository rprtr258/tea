package viewport

import (
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
)

// KeyMap defines the keybindings for the viewport. Note that you don't
// necessary need to use keybindings at all; the viewport can be controlled
// programmatically with methods like Model.LineDown(1). See the GoDocs for
// details.
type KeyMap struct {
	PageDown     key.Binding
	PageUp       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	Down         key.Binding
	Up           key.Binding
}

// DefaultKeyMap returns a set of pager-like default keybindings.
var DefaultKeyMap = KeyMap{
	PageDown: key.Binding{
		Keys: []string{"pgdown", " ", "f"},
		Help: key.Help{"f/pgdn", "page down"},
	},
	PageUp: key.Binding{
		Keys: []string{"pgup", "b"},
		Help: key.Help{"b/pgup", "page up"},
	},
	HalfPageUp: key.Binding{
		Keys: []string{"u", "ctrl+u"},
		Help: key.Help{"u", "½ page up"},
	},
	HalfPageDown: key.Binding{
		Keys: []string{"d", "ctrl+d"},
		Help: key.Help{"d", "½ page down"},
	},
	Up: key.Binding{
		Keys: []string{"up", "k"},
		Help: key.Help{"↑/k", "up"},
	},
	Down: key.Binding{
		Keys: []string{"down", "j"},
		Help: key.Help{"↓/j", "down"},
	},
}

// Model is Tea model for this viewport element.
type Model struct {
	Width  int
	Height int
	KeyMap KeyMap

	// Whether or not to respond to the mouse.
	// The mouse must be enabled in Tea for this to work.
	MouseWheelEnabled bool

	// The number of lines the mouse wheel will scroll.
	// By default, this is 3.
	MouseWheelDelta int

	// YOffset is the vertical scroll position.
	YOffset int
}

// New create model with given width, height and default key mappings.
func New(width, height int) Model {
	return Model{
		Width:             width,
		Height:            height,
		KeyMap:            DefaultKeyMap,
		MouseWheelEnabled: true,
		MouseWheelDelta:   3,
		YOffset:           0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// AtTop returns whether or not the viewport is at the very top position.
func (m *Model) AtTop() bool {
	return m.YOffset <= 0
}

// SetYOffset sets the Y offset.
func (m *Model) SetYOffset(n int) {
	m.YOffset = max(n, 0)
}

// LineDown moves the view down by the given number of lines.
func (m *Model) LineDown(n int) {
	if n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach the bottom.
	m.SetYOffset(m.YOffset + n)
}

// LineUp moves the view down by the given number of lines. Returns the new
// lines to show.
func (m *Model) LineUp(n int) {
	if m.AtTop() || n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	m.SetYOffset(m.YOffset - n)
}

// GotoTop sets the viewport to the top position.
func (m *Model) GotoTop() {
	m.SetYOffset(0)
}

// GotoBottom sets the viewport to the bottom position.
func (m *Model) GotoBottom(contentHeight int) {
	m.SetYOffset(max(0, contentHeight-m.Height))
}

// Update handles standard message-based viewport updates.
func (m *Model) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.KeyMap.PageDown):
			// move view down by the number of lines in the viewport
			m.LineDown(m.Height)
		case key.Matches(msg, m.KeyMap.PageUp):
			// move view up by one height of the viewport
			m.LineUp(m.Height)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			// move view down by half the height of the viewport
			m.LineDown(m.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			// move view up by half the height of the viewport
			m.LineUp(m.Height / 2)
		case key.Matches(msg, m.KeyMap.Down):
			m.LineDown(1)
		case key.Matches(msg, m.KeyMap.Up):
			m.LineUp(1)
		}
	case tea.MsgMouse:
		if !m.MouseWheelEnabled {
			break
		}

		switch msg.Type {
		case tea.MouseWheelUp:
			m.LineUp(m.MouseWheelDelta)
		case tea.MouseWheelDown:
			m.LineDown(m.MouseWheelDelta)
		}
	}
}

// View renders the viewport into a string.
func (m *Model) View(vb tea.Viewbox, lines func(tea.Viewbox, int)) {
	for i := 0; i < min(m.Height, vb.Height); i++ {
		lines(vb.Row(i), m.YOffset+i)
	}
}
