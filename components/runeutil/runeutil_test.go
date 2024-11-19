package runeutil

import (
	"testing"
	"unicode/utf8"

	"github.com/rprtr258/assert"
)

func TestSanitize(t *testing.T) {
	for _, test := range []struct {
		input, output string
	}{
		{"", ""},
		{"x", "x"},
		{"\n", "XX"},
		{"\na\n", "XXaXX"},
		{"\n\n", "XXXX"},
		{"\t", ""},
		{"hello", "hello"},
		{"hel\nlo", "helXXlo"},
		{"hel\rlo", "helXXlo"},
		{"hel\tlo", "hello"},
		{"he\n\nl\tlo", "heXXXXllo"},
		{"he\tl\n\nlo", "helXXXXlo"},
		{"hel\x1blo", "hello"},
		{"hello\xc2", "hello"}, // invalid utf8
	} {
		runes := make([]rune, 0, len(test.input))
		b := []byte(test.input)
		for i, w := 0, 0; i < len(b); i += w {
			var r rune
			r, w = utf8.DecodeRune(b[i:])
			runes = append(runes, r)
		}
		t.Logf("input runes: %+v", runes)
		s := NewSanitizer(ReplaceNewlines("XX"), ReplaceTabs(""))
		result := s.Sanitize(runes)
		rs := string(result)
		assert.Equal(assert.Wrap(t).Msgf("input: %q, result: %v", test.input, result), test.output, rs)
	}
}
