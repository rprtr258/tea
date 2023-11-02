package lipgloss

import (
	"strconv"

	"github.com/muesli/termenv"
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
//	ansiColor := lipgloss.FgColor("21")
//	hexColor := lipgloss.FgColor("#0000ff")
func FgColor(s string) Raw {
	if s[0] == '#' {
		return Raw(scuf.FgRGB(scuf.MustParseHexRGB(s)))
	}

	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
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
//	colorA := lipgloss.ANSIColor(21)
//	colorB := lipgloss.Color("21")
func ANSIColor(x uint) TerminalColor {
	return Raw(scuf.FgANSI256(int(x)))
}

// AdaptiveColor provides color options for light and dark backgrounds. The
// appropriate color will be returned at runtime based on the darkness of the
// terminal background color.
//
// Example usage:
//
//	color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}
type AdaptiveColor struct {
	Light TerminalColor
	Dark  TerminalColor
}

func (ac AdaptiveColor) color() scuf.Modifier {
	return fun.IF(_renderer.HasDarkBackground(), ac.Dark, ac.Light).color()
}

// CompleteColor specifies exact values for truecolor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed.
type CompleteColor struct {
	TrueColor TerminalColor
	ANSI256   TerminalColor
	ANSI      TerminalColor
}

func (c CompleteColor) color() scuf.Modifier {
	switch _renderer.ColorProfile() {
	case termenv.TrueColor:
		return c.TrueColor.color()
	case termenv.ANSI256:
		return c.ANSI256.color()
	case termenv.ANSI:
		return c.ANSI.color()
	default:
		return nil
	}
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
