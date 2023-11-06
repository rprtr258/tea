package list

import "github.com/rprtr258/tea/components/key"

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu.
type KeyMap struct {
	// Keybindings used when browsing the list.
	CursorUp    key.Binding
	CursorDown  key.Binding
	NextPage    key.Binding
	PrevPage    key.Binding
	GoToStart   key.Binding
	GoToEnd     key.Binding
	Filter      key.Binding
	ClearFilter key.Binding

	// Keybindings used when setting a filter.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
var DefaultKeyMap = KeyMap{
	// Browsing.
	CursorUp: key.NewBinding(
		[]string{"up", "k"},
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		[]string{"down", "j"},
		key.WithHelp("↓/j", "down"),
	),
	PrevPage: key.NewBinding(
		[]string{"left", "h", "pgup", "b", "u"},
		key.WithHelp("←/h/pgup", "prev page"),
	),
	NextPage: key.NewBinding(
		[]string{"right", "l", "pgdown", "f", "d"},
		key.WithHelp("→/l/pgdn", "next page"),
	),
	GoToStart: key.NewBinding(
		[]string{"home", "g"},
		key.WithHelp("g/home", "go to start"),
	),
	GoToEnd: key.NewBinding(
		[]string{"end", "G"},
		key.WithHelp("G/end", "go to end"),
	),
	Filter: key.NewBinding(
		[]string{"/"},
		key.WithHelp("/", "filter"),
	),
	ClearFilter: key.NewBinding(
		[]string{"esc"},
		key.WithHelp("esc", "clear filter"),
	),

	// Filtering.
	CancelWhileFiltering: key.NewBinding(
		[]string{"esc"},
		key.WithHelp("esc", "cancel"),
	),
	AcceptWhileFiltering: key.NewBinding(
		[]string{"enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"},
		key.WithHelp("enter", "apply filter"),
	),

	// Toggle help.
	ShowFullHelp: key.NewBinding(
		[]string{"?"},
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		[]string{"?"},
		key.WithHelp("?", "close help"),
	),

	// Quitting.
	Quit: key.NewBinding(
		[]string{"q", "esc"},
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding([]string{"ctrl+c"}),
}
