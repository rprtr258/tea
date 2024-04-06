package styles

import (
	"fmt"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/rprtr258/assert"
	"github.com/rprtr258/scuf"
)

func TestSetColorProfile(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		col      scuf.Modifier
		expected string
	}{
		"ascii": {
			nil,
			"hello",
		},
		"ansi": {
			scuf.FgANSI(12),
			"\x1b[94mhello\x1b[0m",
		},
		"ansi256": {
			scuf.FgANSI(62),
			"\x1b[38;5;62mhello\x1b[0m",
		},
		"truecolor": {
			scuf.FgRGB(scuf.MustParseHexRGB("#5956E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, scuf.String("hello", test.col))
		})
	}
}

func TestRGBA(t *testing.T) {
	type testcase struct {
		darkBg   bool
		input    scuf.Modifier
		expected string
	}
	assert.TableSlice(t, []testcase{
		// styles.Color
		{
			true,
			FgColor("#FF0000"),
			"FF0000",
		},
		{
			true,
			FgColor("9"),
			"FF0000",
		},
		{
			true,
			FgColor("21"),
			"0000FF",
		},
		{
			true,
			FgColor("#FFFFFF"),
			"FFFFFF",
		},
		{
			true,
			FgColor("#0000FF"),
			"0000FF",
		},
		{
			true,
			FgColor("#000000"),
			"000000",
		},
		// styles.AdaptiveColor
		{
			true,
			FgAdaptiveColor("#0000FF", "#FF0000"),
			"FF0000",
		},
		{
			false,
			FgAdaptiveColor("#0000FF", "#FF0000"),
			"0000FF",
		},
		{
			true,
			FgAdaptiveColor("21", "9"),
			"FF0000",
		},
		{
			false,
			FgAdaptiveColor("21", "9"),
			"0000FF",
		},
		// styles.CompleteAdaptiveColor
		// dark
		{
			true,
			adaptiveColor(FgColor("#0000FF"), FgColor("#FF0000")),
			"FF0000",
		},
		{
			true,
			adaptiveColor(FgColor("#FF0000"), FgColor("#FFFFFF")),
			"FFFFFF",
		},
		{
			true,
			adaptiveColor(FgColor("#FF0000"), FgColor("#0000FF")),
			"0000FF",
		},
		{
			false,
			adaptiveColor(FgColor("#0000FF"), FgColor("#FF0000")),
			"0000FF",
		},
	}, func(t *testing.T, test testcase) {
		_renderer.SetHasDarkBackground(test.darkBg)

		c, err := colorful.Hex(scuf.ToHex(test.input))
		assert.NoError(t, err)
		r, g, b, _ := c.RGBA()
		assert.Equal(t, test.expected, fmt.Sprintf("%02X%02X%02X", r/256, g/256, b/256))
	})
}
