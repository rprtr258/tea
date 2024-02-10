package help

import (
	"cmp"

	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
)

// KeyMap is a map of keybindings used to generate help. Since it's an
// interface it can be any type, though struct or a map[string][]key.Binding
// are likely candidates.
//
// Note that if a key is disabled (via key.Binding.SetEnabled) it will not be
// rendered in the help view, so in theory generated help should self-manage.
type KeyMap interface {
	// ShortHelp returns a slice of bindings to be displayed in the short
	// version of the help. The help bubble will render help in the order in
	// which the help items are returned here.
	ShortHelp() []key.Binding

	// FullHelp returns an extended group of help items, grouped by columns.
	// The help bubble will render the help in the order in which the help
	// items are returned here.
	FullHelp() [][]key.Binding
}

// Styles is a set of available style definitions for the Help bubble.
type Styles struct {
	Ellipsis styles.Style

	// Styling for the short help
	ShortKey       styles.Style
	ShortDesc      styles.Style
	ShortSeparator styles.Style

	// Styling for the full help
	FullKey  styles.Style
	FullDesc styles.Style
}

// Model contains the state of the help view.
type Model struct {
	Width   int
	ShowAll bool // if true, render the "full" help menu

	ShortSeparator string

	// The symbol we use in the short help when help items have been truncated
	// due to width. Periods of ellipsis by default.
	Ellipsis string

	Styles Styles
}

// New creates a new help view with some useful defaults.
func New() Model {
	keyStyle := styles.Style{}.Foreground(styles.FgAdaptiveColor("#909090", "#626262"))
	descStyle := styles.Style{}.Foreground(styles.FgAdaptiveColor("#B2B2B2", "#4A4A4A"))
	sepStyle := styles.Style{}.Foreground(styles.FgAdaptiveColor("#DDDADA", "#3C3C3C"))
	return Model{
		ShortSeparator: " • ",
		Ellipsis:       "…",
		Styles: Styles{
			ShortKey:       keyStyle,
			ShortDesc:      descStyle,
			ShortSeparator: sepStyle,
			Ellipsis:       sepStyle.Copy(),
			FullKey:        keyStyle.Copy(),
			FullDesc:       descStyle.Copy(),
		},
	}
}

// Update helps satisfy the Tea Model interface. It's a no-op.
func (m *Model) Update(_ tea.Msg) []tea.Cmd {
	return nil
}

// View renders the help view's current state.
func (m *Model) View(vb tea.Viewbox, k KeyMap) {
	if m.ShowAll {
		m.FullHelpView(vb, k.FullHelp())
	} else {
		m.ShortHelpView(vb, k.ShortHelp())
	}
}

// ShortHelpView renders a single line help view from a slice of keybindings.
// If the line is longer than the maximum width it will be gracefully
// truncated, showing only as many help items as possible.
func (m *Model) ShortHelpView(vb tea.Viewbox, bindings []key.Binding) {
	if len(bindings) == 0 {
		return
	}

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		if i > 0 {
			vb = vb.Styled(m.Styles.ShortSeparator).WriteLineX(m.ShortSeparator)
		}

		vb = vb.Styled(m.Styles.ShortKey).WriteLineX(kb.Help.Key())
		vb = vb.PaddingLeft(1)
		vb = vb.Styled(m.Styles.ShortDesc).WriteLineX(kb.Help.Desc())
	}
}

func maxFunc[T any, R cmp.Ordered](slice []T, f func(T) R) R {
	res := f(slice[0])
	for _, v := range slice[1:] {
		res = max(res, f(v))
	}
	return res
}

// FullHelpView renders help columns from a slice of key binding slices.
// Each top level slice entry renders into a column.
func (m *Model) FullHelpView(vb tea.Viewbox, groups [][]key.Binding) {
	// Iterate over groups to build columns
	for _, group := range groups {
		if !fun.Any(key.Binding.Enabled, group...) {
			continue
		}

		var keys, descriptions []string
		// Separate keys and descriptions into different slices
		for _, kb := range group {
			if !kb.Enabled() {
				continue
			}

			keys = append(keys, kb.Help.Key())
			descriptions = append(descriptions, kb.Help.Desc())
		}

		maxKeyLength := maxFunc(keys, ansi.PrintableRuneWidth)
		maxDescLength := maxFunc(descriptions, ansi.PrintableRuneWidth)

		vbKeys := vb.MaxWidth(maxKeyLength).Styled(m.Styles.FullKey)
		vbDescs := vb.PaddingLeft(maxKeyLength + 1).MaxWidth(maxDescLength).Styled(m.Styles.FullDesc)
		for i, key := range keys {
			vbKeys.PaddingTop(i).WriteLine(key)
			vbDescs.PaddingTop(i).WriteLine(descriptions[i])
		}

		vb = vb.PaddingLeft(maxKeyLength + 1 + maxDescLength + 3)
	}
}
