package ansi

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/rprtr258/scuf"
)

// BaseElement renders a styled primitive element.
type BaseElement struct {
	Token  string
	Prefix string
	Suffix string
	Style  StylePrimitive
}

func formatToken(format, token string) (string, error) {
	tmpl, err := template.New(format).Funcs(TemplateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, map[string]any{
		"text": token,
	})
	return b.String(), err
}

func renderText(w io.Writer, rules StylePrimitive, s string) {
	if s == "" {
		return
	}

	if rules.Upper != nil && *rules.Upper {
		s = strings.ToUpper(s)
	}
	if rules.Lower != nil && *rules.Lower {
		s = strings.ToLower(s)
	}
	if rules.Title != nil && *rules.Title {
		s = strings.Title(s) //nolint:staticcheck
	}

	out := []scuf.Modifier{}
	if rules.ForegroundColor != nil {
		out = append(out, scuf.FgRGB(scuf.MustParseHexRGB(*rules.ForegroundColor)))
	}
	if rules.BackgroundColor != nil {
		out = append(out, scuf.BgRGB(scuf.MustParseHexRGB(*rules.BackgroundColor)))
	}
	if rules.Underline != nil && *rules.Underline {
		out = append(out, scuf.ModUnderline)
	}
	if rules.Bold != nil && *rules.Bold {
		out = append(out, scuf.ModBold)
	}
	if rules.Italic != nil && *rules.Italic {
		out = append(out, scuf.ModItalic)
	}
	if rules.CrossedOut != nil && *rules.CrossedOut {
		out = append(out, scuf.ModCrossout)
	}
	if rules.Overlined != nil && *rules.Overlined {
		out = append(out, scuf.ModOverline)
	}
	if rules.Inverse != nil && *rules.Inverse {
		out = append(out, scuf.ModReverse)
	}
	if rules.Blink != nil && *rules.Blink {
		out = append(out, scuf.ModBlink)
	}

	_, _ = w.Write([]byte(scuf.String(s, out...)))
}

func (e *BaseElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack

	renderText(w, bs.Current().Style.StylePrimitive, e.Prefix)
	defer func() {
		renderText(w, bs.Current().Style.StylePrimitive, e.Suffix)
	}()

	rules := bs.With(e.Style)
	// render unstyled prefix/suffix
	renderText(w, bs.Current().Style.StylePrimitive, rules.BlockPrefix)
	defer func() {
		renderText(w, bs.Current().Style.StylePrimitive, rules.BlockSuffix)
	}()

	// render styled prefix/suffix
	renderText(w, rules, rules.Prefix)
	defer func() {
		renderText(w, rules, rules.Suffix)
	}()

	s := e.Token
	if rules.Format != "" {
		var err error
		s, err = formatToken(rules.Format, s)
		if err != nil {
			return err
		}
	}

	renderText(w, rules, s)
	return nil
}
