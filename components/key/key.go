// Package key provides some types and functions for generating user-definable
// keymappings useful in Tea components. There are a few different ways
// you can define a keymapping with this package. Here's one example:
//
//	type KeyMap struct {
//	    Up key.Binding
//	    Down key.Binding
//	}
//
//	var DefaultKeyMap = KeyMap{
//	    Up: key.NewBinding(
//	        key.WithKeys("k", "up"),        // actual keybindings
//	        key.WithHelp("↑/k", "move up"), // corresponding help text
//	    ),
//	    Down: key.NewBinding(
//	        key.WithKeys("j", "down"),
//	        key.WithHelp("↓/j", "move down"),
//	    ),
//	}
//
//	func (m *model) Update(msg tea.Msg) tea.Cmd {
//	    switch msg := msg.(type) {
//	    case tea.MsgKey:
//	        switch {
//	        case key.Matches(msg, DefaultKeyMap.Up):
//	            // The user pressed up
//	        case key.Matches(msg, DefaultKeyMap.Down):
//	            // The user pressed down
//	        }
//	    }
//
//	    // ...
//	}
//
// The help information, which is not used in the example above, can be used
// to render help text for keystrokes in your views.
package key

import (
	"github.com/rprtr258/tea"
)

// Help is help information for a given keybinding.
type Help [2]string // Key, Desc

func (h Help) Key() string  { return h[0] }
func (h Help) Desc() string { return h[1] }

// Binding describes a set of keybindings and, optionally, their associated
// help text.
type Binding struct {
	Keys     []string
	Help     Help
	Disabled bool
}

// Enabled returns whether or not the keybinding is enabled.
// Disabled keybindings won't be activated and won't show up in help.
// Keybindings are enabled by default.
func (b Binding) Enabled() bool {
	return !b.Disabled && len(b.Keys) != 0
}

// SetEnabled enables or disables the keybinding.
func (b *Binding) SetEnabled(v bool) {
	b.Disabled = !v
}

// Unbind removes the keys and help from this binding, effectively nullifying
// it. This is a step beyond disabling it, since applications can enable
// or disable key bindings based on application state.
func (b *Binding) Unbind() {
	b.Keys = nil
	b.Help = Help{}
}

// Matches checks if the given MsgKey matches the given bindings.
func Matches(key tea.MsgKey, bindings ...Binding) bool {
	keys := key.String()
	for _, binding := range bindings {
		for _, bindingKey := range binding.Keys {
			if keys == bindingKey && binding.Enabled() {
				return true
			}
		}
	}
	return false
}
