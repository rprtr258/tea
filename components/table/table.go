package table

import (
	"github.com/mattn/go-runewidth"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/box"
	"github.com/rprtr258/tea/components/headless/table"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/viewport"
	"github.com/rprtr258/tea/styles"
)

// Column defines the table structure
type Column = table.Column[[]string]

// Model defines a state for the table widget
type Model struct {
	Table  *table.Table[[]string]
	widths []int

	styles   Styles
	KeyMap   KeyMap
	viewport viewport.Model
}

// KeyMap defines keybindings.
// It satisfies to the help.KeyMap interface, which is used to render the menu.
type KeyMap struct {
	LineUp, LineDown         key.Binding
	PageUp, PageDown         key.Binding
	HalfPageUp, HalfPageDown key.Binding
	GotoTop, GotoBottom      key.Binding
}

// DefaultKeyMap returns a default set of keybindings
var DefaultKeyMap = KeyMap{
	LineUp: key.Binding{
		Keys: []string{"up", "k"},
		Help: key.Help{"↑/k", "up"},
	},
	LineDown: key.Binding{
		Keys: []string{"down", "j"},
		Help: key.Help{"↓/j", "down"},
	},
	PageUp: key.Binding{
		Keys: []string{"b", "pgup"},
		Help: key.Help{"b/pgup", "page up"},
	},
	PageDown: key.Binding{
		Keys: []string{"f", "pgdown", " "},
		Help: key.Help{"f/pgdn", "page down"},
	},
	HalfPageUp: key.Binding{
		Keys: []string{"u", "ctrl+u"},
		Help: key.Help{"u", "½ page up"},
	},
	HalfPageDown: key.Binding{
		Keys: []string{"d", "ctrl+d"},
		Help: key.Help{"d", "½ page down"},
	},
	GotoTop: key.Binding{
		Keys: []string{"home", "g"},
		Help: key.Help{"g/home", "go to start"},
	},
	GotoBottom: key.Binding{
		Keys: []string{"end", "G"},
		Help: key.Help{"G/end", "go to end"},
	},
}

// Styles contains style definitions for this list component.
// By default, these values are generated by DefaultStyles.
type Styles struct {
	Header   styles.Style
	Cell     styles.Style
	Selected styles.Style
}

// DefaultStyles returns a set of default style definitions for this table
var DefaultStyles = Styles{
	Selected: styles.Style{}.Bold(true).Foreground(styles.FgColor("212")),
	Header:   styles.Style{}.Bold(true), /*.Padding(0, 1)*/
	Cell:     styles.Style{},            /*.Padding(0, 1)*/
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option = func(*Model)

// New creates a new model for the table widget
func New(
	cols []Column,
	widths []int,
	rows [][]string,
	height, width int,
	styles Styles,
	keyMap KeyMap,
) Model {
	m := Model{
		Table:    table.New(cols, rows),
		widths:   widths,
		viewport: viewport.New(0, 20),

		KeyMap: keyMap, // DefaultKeyMap,
		styles: styles, // DefaultStyles,
	}
	m.Table.MoveTo(0)
	m.viewport.Height = height
	m.viewport.Width = width
	return m
}

// SetStyles sets the table styles
func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

// Update is the Tea update loop
func (m *Model) Update(msg tea.Msg) {
	if msg, ok := msg.(tea.MsgKey); ok {
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
	}
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m *Model) SelectedRow() []string {
	if m.Table.Cursor() == -1 {
		return nil
	}

	return m.Table.Selected()
}

func (m *Model) Height() int     { return m.viewport.Height } // Height returns the viewport height of the table
func (m *Model) SetHeight(h int) { m.viewport.Height = h }    // SetHeight sets the height of the viewport of the table

func (m *Model) Width() int     { return m.viewport.Width } // Width returns the viewport width of the table
func (m *Model) SetWidth(w int) { m.viewport.Width = w }    // SetWidth sets the width of the viewport of the table

func (m *Model) Cursor() int     { return m.Table.Cursor() } // Cursor returns the index of the selected row
func (m *Model) SetCursor(n int) { m.Table.MoveTo(n) }       // SetCursor sets the cursor position in the table

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.Table.MoveUp(n)
	m.viewport.SetYOffset(m.Table.Cursor())
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {
	m.Table.MoveDown(n)
	m.viewport.SetYOffset(m.Table.Cursor())
}

// GotoTop moves the selection to the first row
func (m *Model) GotoTop() {
	m.MoveUp(m.Table.Cursor())
}

// GotoBottom moves the selection to the last row
func (m *Model) GotoBottom() {
	m.MoveDown(m.Table.RowsCount())
}

// View renders the component
func (m *Model) View(vb tea.Viewbox) {
	const _gap = 2
	// 2 for borders, 2 for margin
	totalWidth := 2 + 2 + _gap*(len(m.widths)-1)
	for _, width := range m.widths {
		totalWidth += width
	}

	vb = vb.
		MaxWidth(totalWidth).
		// 2 for borders, 1 for header, 1 for split line
		MaxHeight(2 + 2 + m.viewport.Height)
	box.Box(
		vb,
		func(vb tea.Viewbox) {
			// header
			vbh := vb.Styled(m.styles.Header).PaddingLeft(1)
			for i, col := range m.Table.Columns() {
				vbh.WriteLine(runewidth.Truncate(col.Title, m.widths[i], "…"))
				vbh = vbh.PaddingLeft(m.widths[i]).PaddingLeft(_gap)
			}

			// split line
			vb.Row(1).Styled(styles.Style{}.Foreground(styles.FgColor("240"))).Fill(box.NormalBorder.Top)

			// rows
			m.viewport.View(vb.PaddingTop(2), func(vbRow tea.Viewbox, i int) {
				if i == m.Table.Cursor() {
					vbRow = vbRow.Styled(m.styles.Selected)
				}

				vbRow = vbRow.PaddingLeft(1)
				for i, value := range m.Table.Rows()[i] {
					vbRow.WriteLine(m.styles.Cell.Render(runewidth.Truncate(value, m.widths[i], "…")))
					vbRow = vbRow.PaddingLeft(m.widths[i]).PaddingLeft(_gap)
				}
			})
		},
		box.NormalBorder,
		box.BorderMaskAll,
		box.Colors(nil),
		box.Colors(styles.FgColor("240")),
	)
}
