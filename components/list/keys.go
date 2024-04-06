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
	CursorUp: key.Binding{
		Keys: []string{"up", "k"},
		Help: key.Help{"↑/k", "up"},
	},
	CursorDown: key.Binding{
		Keys: []string{"down", "j"},
		Help: key.Help{"↓/j", "down"},
	},
	PrevPage: key.Binding{
		Keys: []string{"left", "h", "pgup", "b", "u"},
		Help: key.Help{"←/h/pgup", "prev page"},
	},
	NextPage: key.Binding{
		Keys: []string{"right", "l", "pgdown", "f", "d"},
		Help: key.Help{"→/l/pgdn", "next page"},
	},
	GoToStart: key.Binding{
		Keys: []string{"home", "g"},
		Help: key.Help{"g/home", "go to start"},
	},
	GoToEnd: key.Binding{
		Keys: []string{"end", "G"},
		Help: key.Help{"G/end", "go to end"},
	},
	Filter: key.Binding{
		Keys: []string{"/"},
		Help: key.Help{"/", "filter"},
	},
	ClearFilter: key.Binding{
		Keys: []string{"esc"},
		Help: key.Help{"esc", "clear filter"},
	},

	// Filtering.
	CancelWhileFiltering: key.Binding{
		Keys: []string{"esc"},
		Help: key.Help{"esc", "cancel"},
	},
	AcceptWhileFiltering: key.Binding{
		Keys: []string{"enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"},
		Help: key.Help{"enter", "apply filter"},
	},

	// Toggle help.
	ShowFullHelp: key.Binding{
		Keys: []string{"?"},
		Help: key.Help{"?", "more"},
	},
	CloseFullHelp: key.Binding{
		Keys: []string{"?"},
		Help: key.Help{"?", "close help"},
	},

	// Quitting.
	Quit: key.Binding{
		Keys: []string{"q", "esc"},
		Help: key.Help{"q", "quit"},
	},
	ForceQuit: key.Binding{Keys: []string{"ctrl+c"}},
}
