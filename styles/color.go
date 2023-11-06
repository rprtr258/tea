package styles

import (
	"strconv"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
)

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor = scuf.Modifier

func FgRGB(hex string) TerminalColor {
	return scuf.FgRGB(scuf.MustParseHexRGB(hex))
}

func BgRGB(hex string) TerminalColor {
	return scuf.BgRGB(scuf.MustParseHexRGB(hex))
}

// FgColor specifies a color by hex or ANSI value. For example:
//
//	ansiColor := styles.FgColor("21")
//	hexColor := styles.FgColor("#0000ff")
func FgColor(s string) TerminalColor {
	if s == "" {
		return nil
	}

	if s[0] == '#' {
		return scuf.FgRGB(scuf.MustParseHexRGB(s))
	}

	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}

	if x < 16 {
		return scuf.FgANSI(x)
	}

	return scuf.FgANSI256(x)
}

func BgColor(s string) TerminalColor {
	if s[0] == '#' {
		return scuf.BgRGB(scuf.MustParseHexRGB(s))
	}

	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}
	return scuf.BgANSI256(x)
}

// ANSIColor is a color specified by an ANSI color value. It's merely syntactic
// sugar for the more general Color function. Invalid colors will render as
// black.
//
// Example usage:
//
//	// These two statements are equivalent.
//	colorA := styles.ANSIColor(21)
//	colorB := styles.Color("21")
func ANSIColor(x uint) TerminalColor {
	return scuf.FgANSI256(int(x))
}

// adaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := styles.FgAdaptiveColor("#0000ff", "#000099")
func adaptiveColor(Light, Dark TerminalColor) TerminalColor {
	return fun.IF(_renderer.HasDarkBackground(), Dark, Light)
}

func FgAdaptiveColor(light, dark string) TerminalColor {
	return adaptiveColor(FgColor(light), FgColor(dark))
}

func BgAdaptiveColor(light, dark string) TerminalColor {
	return adaptiveColor(BgColor(light), BgColor(dark))
}
