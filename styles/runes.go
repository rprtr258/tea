package styles

import (
	"strings"
	"unicode/utf8"
)

// StyleRunes apply a given style to runes at the given indices in the string.
// Note that you must provide styling options for both matched and unmatched
// runes. Indices out of bounds will be ignored.
func StyleRunes(str string, indices []int, matched, unmatched Style) string {
	// Convert slice of indices to a map for easier lookups
	m := make(map[int]struct{})
	for _, i := range indices {
		m[i] = struct{}{}
	}

	runeCount := utf8.RuneCountInString(str)

	i := 0
	var out, group strings.Builder
	for _, r := range str {
		group.WriteRune(r)

		_, matches := m[i]
		_, nextMatches := m[i+1]

		if matches != nextMatches || i == runeCount-1 {
			// Flush
			style := unmatched
			if matches {
				style = matched
			}

			out.WriteString(style.Render(group.String()))
			group.Reset()
		}

		i++
	}

	return out.String()
}
