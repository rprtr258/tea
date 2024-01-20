package styles

import (
	"fmt"
	"io"
	"testing"

	"github.com/rprtr258/assert"
	"github.com/rprtr258/scuf"
)

func TestStyleRender(t *testing.T) {
	type testcase struct {
		style    Style
		expected string
	}
	assert.TableSlice(t, []testcase{
		{
			Style{}.Foreground(FgColor("#5956E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			Style{}.Foreground(FgAdaptiveColor("#fffe12", "#5956E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			Style{}.Bold(true),
			"\x1b[1mhello\x1b[0m",
		},
		{
			Style{}.Italic(),
			"\x1b[3mhello\x1b[0m",
		},
		{
			Style{}.Underline(),
			"\x1b[4;4mh\x1b[0m\x1b[4;4me\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4mo\x1b[0m",
		},
		{
			Style{}.Blink(),
			"\x1b[5mhello\x1b[0m",
		},
		{
			Style{}.Faint(),
			"\x1b[2mhello\x1b[0m",
		},
	}, func(t *testing.T, test testcase) {
		_renderer.SetHasDarkBackground(true)

		s := test.style.Copy().SetString("hello")
		res := s.Render()
		assert.Equal(t, test.expected, res)
	})
}

func TestStyleCustomRender(t *testing.T) {
	for i, tc := range []struct {
		style    Style
		expected string
	}{
		{
			Style{}.Foreground(FgColor("#5956E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			Style{}.Foreground(FgAdaptiveColor("#fffe12", "#5A56E0")),
			"\x1b[38;2;255;254;18mhello\x1b[0m",
		},
		{
			Style{}.Bold(true),
			"\x1b[1mhello\x1b[0m",
		},
		{
			Style{}.Italic(),
			"\x1b[3mhello\x1b[0m",
		},
		{
			Style{}.Underline(),
			"\x1b[4;4mh\x1b[0m\x1b[4;4me\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4mo\x1b[0m",
		},
		{
			Style{}.Blink(),
			"\x1b[5mhello\x1b[0m",
		},
		{
			Style{}.Faint(),
			"\x1b[2mhello\x1b[0m",
		},
		{
			Style{}.Faint(),
			"\x1b[2mhello\x1b[0m",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := NewRenderer(io.Discard)
			r.SetHasDarkBackground(false)
			_renderer = r

			s := tc.style.Copy().SetString("hello")
			res := s.Render()
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestValueCopy(t *testing.T) {
	t.Parallel()

	s := Style{}.
		Bold(true)

	i := s
	i.Bold(false)

	assert.Equal(t, s.GetBold(), i.GetBold())
}

func TestStyleInherit(t *testing.T) {
	t.Parallel()

	s := Style{}.
		Bold(true).
		Italic().
		Underline().
		Strikethrough(true).
		Blink().
		Faint().
		Foreground(FgColor("#ffffff")).
		Background(BgColor("#111111"))

	i := Style{}.Inherit(s)

	assert.Equal(t, s.GetBold(), i.GetBold())
	assert.Equal(t, s.GetItalic(), i.GetItalic())
	assert.Equal(t, s.GetUnderline(), i.GetUnderline())
	assert.Equal(t, s.GetStrikethrough(), i.GetStrikethrough())
	assert.Equal(t, s.GetBlink(), i.GetBlink())
	assert.Equal(t, s.GetFaint(), i.GetFaint())
	assert.Equal(t, s.GetForeground(), i.GetForeground())
	assert.Equal(t, s.GetBackground(), i.GetBackground())
}

func TestStyleCopy(t *testing.T) {
	t.Parallel()

	s := Style{}.
		Bold(true).
		Italic().
		Underline().
		Strikethrough(true).
		Blink().
		Faint().
		Foreground(FgColor("#ffffff")).
		Background(BgColor("#111111"))

	i := s.Copy()

	assert.Equal(t, s.GetBold(), i.GetBold())
	assert.Equal(t, s.GetItalic(), i.GetItalic())
	assert.Equal(t, s.GetUnderline(), i.GetUnderline())
	assert.Equal(t, s.GetStrikethrough(), i.GetStrikethrough())
	assert.Equal(t, s.GetBlink(), i.GetBlink())
	assert.Equal(t, s.GetFaint(), i.GetFaint())
	assert.Equal(t, s.GetForeground(), i.GetForeground())
	assert.Equal(t, s.GetBackground(), i.GetBackground())
}

func TestStyleUnset(t *testing.T) {
	t.Parallel()

	s := Style{}.Bold(true)
	assert.True(t, s.GetBold())
	s.UnsetBold()
	assert.False(t, s.GetBold())

	s = Style{}.Italic()
	assert.True(t, s.GetItalic())
	s.UnsetItalic()
	assert.False(t, s.GetItalic())

	s = Style{}.Underline()
	assert.True(t, s.GetUnderline())
	s.UnsetUnderline()
	assert.False(t, s.GetUnderline())

	s = Style{}.Strikethrough(true)
	assert.True(t, s.GetStrikethrough())
	s.UnsetStrikethrough()
	assert.False(t, s.GetStrikethrough())

	s = Style{}.Reverse(true)
	assert.True(t, s.GetReverse())
	s.UnsetReverse()
	assert.False(t, s.GetReverse())

	s = Style{}.Blink()
	assert.True(t, s.GetBlink())
	s.UnsetBlink()
	assert.False(t, s.GetBlink())

	s = Style{}.Faint()
	assert.True(t, s.GetFaint())
	s.UnsetFaint()
	assert.False(t, s.GetFaint())

	// colors
	colcol := scuf.Modifier(scuf.FgRGB(scuf.MustParseHexRGB("#ffffff")))
	s = Style{}.Foreground(colcol)
	assert.Equal(t, colcol, s.GetForeground())
	s.UnsetForeground()
	assert.NotEqual(t, colcol, s.GetForeground())

	s = Style{}.Background(colcol)
	assert.Equal(t, colcol, s.GetBackground())
	s.UnsetBackground()
	assert.NotEqual(t, colcol, s.GetBackground())
}

func TestStyleValue(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		style    Style
		expected string
	}{
		"empty": {
			style:    Style{},
			expected: "foo",
		},
		"set string": {
			style:    Style{}.SetString("bar"),
			expected: "bar foo",
		},
		"set string with bold": {
			style:    Style{}.SetString("bar").Bold(true),
			expected: "\x1b[1mbar foo\x1b[0m",
		},
		"new style with string": {
			style:    Style{}.SetString("bar", "foobar"),
			expected: "bar foobar foo",
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := test.style.Render("foo")
			assert.Equal(t, test.expected, res)
		})
	}
}

func BenchmarkStyleRender(b *testing.B) {
	s := Style{}.
		Bold(true).
		Foreground(FgColor("#ffffff"))

	for i := 0; i < b.N; i++ {
		s.Render("Hello world")
	}
}
