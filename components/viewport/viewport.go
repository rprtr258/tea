package viewport

import (
	"cmp"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
	"github.com/samber/lo"
)

const _spacebar = " "

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
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", _spacebar, "f"),
		key.WithHelp("f/pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "b"),
		key.WithHelp("b/pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("u", "ctrl+u"),
		key.WithHelp("u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
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

	// Style applies a styles style to the viewport. Realistically, it's most
	// useful for setting borders, margins and padding.
	Style styles.Style

	lines []string
}

// New create model with given width, height and default key mappings.
func New(width, height int) Model {
	return Model{
		Width:             width,
		Height:            height,
		KeyMap:            DefaultKeyMap,
		MouseWheelEnabled: true,
		MouseWheelDelta:   3,
		lines:             nil,
		YOffset:           0,
		Style:             styles.Style{},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// AtTop returns whether or not the viewport is at the very top position.
func (m *Model) AtTop() bool {
	return m.YOffset <= 0
}

// AtBottom returns whether or not the viewport is at or past the very bottom position.
func (m *Model) AtBottom() bool {
	return m.YOffset >= m.maxYOffset()
}

// PastBottom returns whether or not the viewport is scrolled beyond the last line.
// This can happen when adjusting the viewport height.
func (m *Model) PastBottom() bool {
	return m.YOffset > m.maxYOffset()
}

// ScrollPercent returns the amount scrolled as a float between 0 and 1.
func (m *Model) ScrollPercent() float64 {
	if m.Height >= len(m.lines) {
		return 1.0
	}

	v := float64(m.YOffset) / float64(len(m.lines)-1-m.Height)
	return clamp(v, 0, 1)
}

// SetContent set the pager's text content.
func (m *Model) SetContent(lines []string) {
	m.lines = lines

	if m.YOffset > len(m.lines)-1 {
		m.GotoBottom()
	}
}

// maxYOffset returns the maximum possible value of the y-offset based on the
// viewport's content and set height.
func (m *Model) maxYOffset() int {
	return max(0, len(m.lines)-m.Height)
}

// visibleLines returns the lines that should currently be visible in the
// viewport.
func (m *Model) visibleLines() []string {
	if len(m.lines) == 0 {
		return nil
	}

	return lo.Slice(m.lines, m.YOffset, m.YOffset+m.Height)
}

// SetYOffset sets the Y offset.
func (m *Model) SetYOffset(n int) {
	m.YOffset = clamp(n, 0, m.maxYOffset())
}

// LineDown moves the view down by the given number of lines.
func (m *Model) LineDown(n int) {
	if m.AtBottom() || n == 0 || len(m.lines) == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach the bottom.
	m.SetYOffset(m.YOffset + n)
}

// LineUp moves the view down by the given number of lines. Returns the new
// lines to show.
func (m *Model) LineUp(n int) {
	if m.AtTop() || n == 0 || len(m.lines) == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	m.SetYOffset(m.YOffset - n)
}

// TotalLineCount returns the total number of lines (both hidden and visible) within the viewport.
func (m *Model) TotalLineCount() int {
	return len(m.lines)
}

// VisibleLineCount returns the number of the visible lines within the viewport.
func (m *Model) VisibleLineCount() int {
	return len(m.visibleLines())
}

// GotoTop sets the viewport to the top position.
func (m *Model) GotoTop() {
	m.SetYOffset(0)
	m.visibleLines()
}

// GotoBottom sets the viewport to the bottom position.
func (m *Model) GotoBottom() {
	m.SetYOffset(m.maxYOffset())
	m.visibleLines()
}

// Update handles standard message-based viewport updates.
func (m *Model) Update(msg tea.Msg) []tea.Cmd {
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

	return nil
}

// View renders the viewport into a string.
func (m *Model) View(vb tea.Viewbox) {
	for _, line := range m.visibleLines() {
		vb.WriteLine(line)
		vb = vb.PaddingTop(1)
	}
}

func clamp[T cmp.Ordered](v, low, high T) T {
	if high < low {
		low, high = high, low
	}
	return min(max(v, low), high)
}
