package list_fancy //nolint:revive,stylecheck

import (
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/components/list"
)

func newItemDelegate[I list.DefaultItem](keys *delegateKeyMap) list.DefaultDelegate[I] {
	d := list.NewDefaultDelegate[I]()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model[I]) []tea.Cmd {
		var title string

		if i, ok := m.SelectedItem(); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) { //nolint:gocritic
		case tea.MsgKey:
			switch {
			case key.Matches(msg, keys.choose):
				return []tea.Cmd{m.CmdNewStatusMessage(statusMessageStyle("You chose " + title))}

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return []tea.Cmd{m.CmdNewStatusMessage(statusMessageStyle("Deleted " + title))}
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

var delegateKeys = &delegateKeyMap{
	choose: key.Binding{
		Keys: []string{"enter"},
		Help: key.Help{"enter", "choose"},
	},
	remove: key.Binding{
		Keys: []string{"x", "backspace"},
		Help: key.Help{"x", "delete"},
	},
}
