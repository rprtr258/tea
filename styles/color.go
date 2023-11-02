package styles

import (
	"strconv"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
)

// TerminalColor is a color intended to be rendered in the terminal.
type TerminalColor interface {
	color() scuf.Modifier
}

func FgRGB(hex string) TerminalColor {
	return Raw(scuf.FgRGB(scuf.MustParseHexRGB(hex)))
}

func BgRGB(hex string) TerminalColor {
	return Raw(scuf.BgRGB(scuf.MustParseHexRGB(hex)))
}

// TODO: remove, just to deal with cringe
type Raw scuf.Modifier

func (r Raw) color() scuf.Modifier {
	return scuf.Modifier(r)
}

// FgColor specifies a color by hex or ANSI value. For example:
//
//	ansiColor := styles.FgColor("21")
//	hexColor := styles.FgColor("#0000ff")
func FgColor(s string) Raw {
	if s == "" {
		return Raw(nil)
	}

	if s[0] == '#' {
		return Raw(scuf.FgRGB(scuf.MustParseHexRGB(s)))
	}

	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}

	if x < 16 {
		return Raw(scuf.FgANSI(x))
	}

	return Raw(scuf.FgANSI256(x))
}

func BgColor(s string) Raw {
	if s[0] == '#' {
		return Raw(scuf.BgRGB(scuf.MustParseHexRGB(s)))
	}

	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}
	return Raw(scuf.BgANSI256(x))
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
	return Raw(scuf.FgANSI256(int(x)))
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := styles.FgAdaptiveColor("#0000ff", "#000099")
type AdaptiveColor struct {
	Light TerminalColor
	Dark  TerminalColor
}

func FgAdaptiveColor(light, dark string) AdaptiveColor {
	return AdaptiveColor{
		Light: FgColor(light),
		Dark:  FgColor(dark),
	}
}

func BgAdaptiveColor(light, dark string) AdaptiveColor {
	return AdaptiveColor{
		Light: BgColor(light),
		Dark:  BgColor(dark),
	}
}

func (ac AdaptiveColor) color() scuf.Modifier {
	return fun.IF(_renderer.HasDarkBackground(), ac.Dark, ac.Light).color()
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor TerminalColor
}

func (c CompleteColor) color() scuf.Modifier {
	return c.TrueColor.color()
}

// CompleteAdaptiveColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles, with separate options for light and dark backgrounds. Automatic
// color degradation will not be performed.
type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

func (cac CompleteAdaptiveColor) color() scuf.Modifier {
	return fun.IF(_renderer.HasDarkBackground(), cac.Dark, cac.Light).color()
}
