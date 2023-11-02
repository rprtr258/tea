package lipgloss

import (
	"fmt"
	"image/color"
	"log"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
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
			scuf.FgANSI256(62),
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

func TestxToColor(t *testing.T) {
	t.Parallel()

	for input, expected := range map[string]uint{
		"#FF0000":       0xFF0000,
		"#00F":          0x0000FF,
		"#6B50FF":       0x6B50FF,
		"invalid color": 0x0,
	} {
		input := input
		expected := expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			h := hexToColor(input)
			o := uint(h.R)<<16 + uint(h.G)<<8 + uint(h.B)
			assert.Equal(t, expected, o)
		})
	}
}

func TestRGBA(t *testing.T) {
	for i, test := range []struct {
		profile  termenv.Profile
		darkBg   bool
		input    TerminalColor
		expected string
	}{
		// lipgloss.Color
		{
			termenv.TrueColor,
			true,
			FgColor("#FF0000"),
			"FF0000",
		},
		{
			termenv.TrueColor,
			true,
			FgColor("9"),
			"FF0000",
		},
		{
			termenv.TrueColor,
			true,
			FgColor("21"),
			"0000FF",
		},
		// lipgloss.AdaptiveColor
		{
			termenv.TrueColor,
			true,
			AdaptiveColor{Light: FgColor("#0000FF"), Dark: FgColor("#FF0000")},
			"FF0000",
		},
		{
			termenv.TrueColor,
			false,
			AdaptiveColor{Light: FgColor("#0000FF"), Dark: FgColor("#FF0000")},
			"0000FF",
		},
		{
			termenv.TrueColor,
			true,
			AdaptiveColor{Light: FgColor("21"), Dark: FgColor("9")},
			"FF0000",
		},
		{
			termenv.TrueColor,
			false,
			AdaptiveColor{Light: FgColor("21"), Dark: FgColor("9")},
			"0000FF",
		},
		// CompleteColor
		{
			termenv.TrueColor,
			true,
			CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			"FF0000",
		},
		{
			termenv.ANSI256,
			true,
			CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			"FFFFFF",
		},
		{
			termenv.ANSI,
			true,
			CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			"0000FF",
		},
		{
			termenv.TrueColor,
			true,
			CompleteColor{TrueColor: FgColor("#000000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			"000000",
		},
		// lipgloss.CompleteAdaptiveColor
		// dark
		{
			termenv.TrueColor,
			true,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#0000FF"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"FF0000",
		},
		{
			termenv.ANSI256,
			true,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("21"), ANSI: FgColor("12")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"FFFFFF",
		},
		{
			termenv.ANSI,
			true,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("9")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"0000FF",
		},
		// light
		{
			termenv.TrueColor,
			false,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#0000FF"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"0000FF",
		},
		{
			termenv.ANSI256,
			false,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("21"), ANSI: FgColor("12")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"0000FF",
		},
		{
			termenv.ANSI,
			false,
			CompleteAdaptiveColor{
				Light: CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("9")},
				Dark:  CompleteColor{TrueColor: FgColor("#FF0000"), ANSI256: FgColor("231"), ANSI: FgColor("12")},
			},
			"FF0000",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_renderer.SetColorProfile(test.profile)
			_renderer.SetHasDarkBackground(test.darkBg)

			c, err := colorful.Hex(scuf.ToHex(test.input.color()))
			log.Println(scuf.ToHex(test.input.color()))
			assert.NoError(t, err)
			r, g, b, _ := c.RGBA()
			assert.Equal(t, test.expected, fmt.Sprintf("%02X%02X%02X", r/256, g/256, b/256))
		})
	}
}

// hexToColor translates a hex color string (#RRGGBB or #RGB) into a color.RGB,
// which satisfies the color.Color interface. If an invalid string is passed
// black with 100% opacity will be returned: or, in hex format, 0x000000FF.
func hexToColor(hex string) color.RGBA {
	if hex == "" || hex[0] != '#' {
		return color.RGBA{0, 0, 0, 0xFF}
	}

	switch len(hex) {
	case 7: // full format: #RRGGBB
		const offset = 4
		return color.RGBA{
			R: hexToByte(hex[1])<<offset + hexToByte(hex[2]),
			G: hexToByte(hex[3])<<offset + hexToByte(hex[4]),
			B: hexToByte(hex[5])<<offset + hexToByte(hex[6]),
			A: 0xFF,
		}
	case 4: // short format: #RGB
		const offset = 0x11
		return color.RGBA{
			R: hexToByte(hex[1]) * offset,
			G: hexToByte(hex[2]) * offset,
			B: hexToByte(hex[3]) * offset,
			A: 0xFF,
		}
	default:
		return color.RGBA{0, 0, 0, 0xFF}
	}
}

func hexToByte(b byte) byte {
	const offset = 10
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + offset
	case b >= 'A' && b <= 'F':
		return b - 'A' + offset
	default: // Invalid, but just return 0.
		return 0
	}
}
