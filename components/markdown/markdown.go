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
type TermRendererOption func(*TermRenderer) error

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

// RenderBytes initializes a new TermRenderer and renders a markdown with a specific style.
func RenderBytes(in []byte, style ansi.StyleConfig) ([]byte, error) {
	r, err := NewTermRenderer(WithStyles(style))
	if err != nil {
		return nil, err
	}

	return r.RenderBytes(in)
}

// Render initializes new TermRenderer and renders markdown with a specific style
func Render(in string, style ansi.StyleConfig) (string, error) {
	b, err := RenderBytes([]byte(in), style)
	return string(b), err
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
	for _, o := range options {
		if err := o(tr); err != nil {
			return nil, err
		}
	}
	ar := ansi.NewRenderer(tr.ansiOptions)
	tr.md.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(ar, 1000),
			),
		),
	)
	return tr, nil
}

// WithBaseURL sets a TermRenderer's base URL.
func WithBaseURL(baseURL string) TermRendererOption {
	return func(tr *TermRenderer) error {
		tr.ansiOptions.BaseURL = baseURL
		return nil
	}
}

// WithColorProfile sets the TermRenderer's color profile
// (TrueColor / ANSI256 / ANSI).
func WithColorProfile(profile termenv.Profile) TermRendererOption {
	return func(tr *TermRenderer) error {
		tr.ansiOptions.ColorProfile = profile
		return nil
	}
}

// WithAutoStyle sets a TermRenderer's styles with either the standard dark
// or light style, depending on the terminal's background color at run-time.
func WithAutoStyle() TermRendererOption {
	return WithStyles(AutoStyle())
}

// WithStyles sets a TermRenderer's styles.
func WithStyles(styles ansi.StyleConfig) TermRendererOption {
	return func(tr *TermRenderer) error {
		tr.ansiOptions.Styles = styles
		return nil
	}
}

// WithStylesFromJSONBytes sets a TermRenderer's styles by parsing styles from
// jsonBytes.
func WithStylesFromJSONBytes(jsonBytes []byte) TermRendererOption {
	return func(tr *TermRenderer) error {
		return json.Unmarshal(jsonBytes, &tr.ansiOptions.Styles)
	}
}

// WithStylesFromJSONFile sets a TermRenderer's styles from a JSON file.
func WithStylesFromJSONFile(filename string) TermRendererOption {
	return func(tr *TermRenderer) error {
		jsonBytes, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonBytes, &tr.ansiOptions.Styles)
	}
}

// WithWordWrap sets a TermRenderer's word wrap.
func WithWordWrap(wordWrap int) TermRendererOption {
	return func(tr *TermRenderer) error {
		tr.ansiOptions.WordWrap = wordWrap
		return nil
	}
}

// WithPreservedNewlines preserves newlines from being replaced.
func WithPreservedNewLines() TermRendererOption {
	return func(tr *TermRenderer) error {
		tr.ansiOptions.PreserveNewLines = true
		return nil
	}
}

// WithEmoji sets a TermRenderer's emoji rendering.
func WithEmoji() TermRendererOption {
	return func(tr *TermRenderer) error {
		emoji.New().Extend(tr.md)
		return nil
	}
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
	b, err := tr.RenderBytes([]byte(in))
	return string(b), err
}

// RenderBytes returns the markdown rendered into a byte slice.
func (tr *TermRenderer) RenderBytes(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := tr.md.Convert(in, &buf)
	return buf.Bytes(), err
}

func AutoStyle() ansi.StyleConfig {
	return fun.IF(termenv.HasDarkBackground(), DarkStyle, LightStyle)
}
