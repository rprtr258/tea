package markdown

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/muesli/termenv"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/tea/components/markdown/ansi"
)

// A TermRendererOption sets an option on a TermRenderer.
type TermRendererOption func(*TermRenderer)

// TermRenderer can be used to render markdown content, posing a depth of
// customization and styles to fit your needs.
type TermRenderer struct {
	md          goldmark.Markdown
	ansiOptions ansi.Options
	buf         bytes.Buffer
	renderBuf   bytes.Buffer
}

func parseStyleFile(stylePath string) (ansi.StyleConfig, error) {
	jsonBytes, err := os.ReadFile(stylePath)
	if err != nil {
		return ansi.StyleConfig{}, fmt.Errorf("read style file: %w", err)
	}

	var res ansi.StyleConfig
	if err := json.Unmarshal(jsonBytes, &res); err != nil {
		return ansi.StyleConfig{}, fmt.Errorf("parse style json spec: %w", err)
	}

	return res, nil
}

// WithBaseURL sets a TermRenderer's base URL.
func WithBaseURL(baseURL string) TermRendererOption {
	return func(tr *TermRenderer) {
		tr.ansiOptions.BaseURL = baseURL
	}
}

// WithColorProfile sets the TermRenderer's color profile (TrueColor / ANSI256 / ANSI)
func WithColorProfile(profile termenv.Profile) TermRendererOption {
	return func(tr *TermRenderer) {
		tr.ansiOptions.ColorProfile = profile
	}
}

// WithAutoStyle sets a TermRenderer's styles with either the standard dark
// or light style, depending on the terminal's background color at run-time.
func WithAutoStyle() TermRendererOption {
	return WithStyles(fun.IF(termenv.HasDarkBackground(), DarkStyle, LightStyle))
}

// WithStyles sets a TermRenderer's styles
func WithStyles(styles ansi.StyleConfig) TermRendererOption {
	return func(tr *TermRenderer) {
		tr.ansiOptions.Styles = styles
	}
}

// WithWordWrap sets a TermRenderer's word wrap
func WithWordWrap(wordWrap int) TermRendererOption {
	return func(tr *TermRenderer) {
		tr.ansiOptions.WordWrap = wordWrap
	}
}

// WithPreservedNewlines preserves newlines from being replaced.
func WithPreservedNewLines() TermRendererOption {
	return func(tr *TermRenderer) {
		tr.ansiOptions.PreserveNewLines = true
	}
}

// WithEmoji sets a TermRenderer's emoji rendering.
func WithEmoji() TermRendererOption {
	return func(tr *TermRenderer) {
		emoji.New().Extend(tr.md)
	}
}

// NewTermRenderer returns a new TermRenderer the given options.
func NewTermRenderer(options ...TermRendererOption) (*TermRenderer, error) {
	tr := &TermRenderer{
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.DefinitionList,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		),
		ansiOptions: ansi.Options{
			WordWrap:     80,
			ColorProfile: termenv.TrueColor,
		},
	}
	for _, opt := range options {
		opt(tr)
	}
	tr.md.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(ansi.NewRenderer(tr.ansiOptions), 1000),
			),
		),
	)
	return tr, nil
}

func (tr *TermRenderer) Read(b []byte) (int, error) {
	return tr.renderBuf.Read(b)
}

func (tr *TermRenderer) Write(b []byte) (int, error) {
	return tr.buf.Write(b)
}

// Close must be called after writing to TermRenderer. You can then retrieve
// the rendered markdown by calling Read.
func (tr *TermRenderer) Close() error {
	err := tr.md.Convert(tr.buf.Bytes(), &tr.renderBuf)
	if err != nil {
		return err
	}

	tr.buf.Reset()
	return nil
}

// Render returns the markdown rendered into a string.
func (tr *TermRenderer) Render(in string) (string, error) {
	var buf bytes.Buffer
	err := tr.md.Convert([]byte(in), &buf)
	return string(buf.Bytes()), err
}

// Render initializes new TermRenderer and renders markdown with a specific style
func Render(in string, style ansi.StyleConfig) (string, error) {
	r, err := NewTermRenderer(WithStyles(style))
	if err != nil {
		return "", err
	}

	return r.Render(in)
}
