package list_fancy

import (
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/key"
	"github.com/rprtr258/tea/bubbles/list"
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

		switch msg := msg.(type) {
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

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
