package ansi

import (
	"io"
	"net/url"
	"strings"

	"github.com/muesli/termenv"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Options is used to configure an ANSIRenderer.
type Options struct {
	BaseURL          string
	WordWrap         int
	PreserveNewLines bool
	ColorProfile     termenv.Profile
	Styles           StyleConfig
}

// Renderer renders markdown content as ANSI escaped sequences.
type Renderer struct {
	context RenderContext
}

// NewRenderer returns a new ANSIRenderer with style and options set.
func NewRenderer(options Options) *Renderer {
	return &Renderer{
		context: NewRenderContext(options),
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	for _, nodeKind := range []ast.NodeKind{
		astext.KindTaskCheckBox,  // checkboxes
		astext.KindStrikethrough, // strikethrough
		east.KindEmoji,           // emoji
		// blocks
		ast.KindDocument,
		ast.KindHeading,
		ast.KindBlockquote,
		ast.KindCodeBlock,
		ast.KindFencedCodeBlock,
		ast.KindHTMLBlock,
		ast.KindList,
		ast.KindListItem,
		ast.KindParagraph,
		ast.KindTextBlock,
		ast.KindThematicBreak,
		// inlines
		ast.KindAutoLink,
		ast.KindCodeSpan,
		ast.KindEmphasis,
		ast.KindImage,
		ast.KindLink,
		ast.KindRawHTML,
		ast.KindText,
		ast.KindString,
		// tables
		astext.KindTable,
		astext.KindTableHeader,
		astext.KindTableRow,
		astext.KindTableCell,
		// definitions
		astext.KindDefinitionList,
		astext.KindDefinitionTerm,
		astext.KindDefinitionDescription,
		// footnotes
		astext.KindFootnote,
		astext.KindFootnoteList,
		astext.KindFootnoteLink,
		astext.KindFootnoteBacklink,
	} {
		reg.Register(nodeKind, r.renderNode)
	}
}

func (r *Renderer) renderNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// _, _ = w.Write([]byte(node.Type.String()))
	writeTo := io.Writer(w)
	bs := r.context.blockStack

	// children get rendered by their parent
	if isChild(node) {
		return ast.WalkContinue, nil
	}

	e := r.NewElement(node, source)
	if entering {
		// everything below the Document element gets rendered into a block buffer
		if bs.Len() > 0 {
			writeTo = io.Writer(bs.Current().Block)
		}

		_, _ = writeTo.Write([]byte(e.Entering))
		if e.Renderer != nil {
			err := e.Renderer.Render(writeTo, r.context)
			if err != nil {
				return ast.WalkStop, err
			}
		}
	} else {
		// everything below the Document element gets rendered into a block buffer
		if bs.Len() > 0 {
			writeTo = io.Writer(bs.Parent().Block)
		}

		// if we're finished rendering the entire document,
		// flush to the real writer
		if node.Type() == ast.TypeDocument {
			writeTo = w
		}

		if e.Finisher != nil {
			err := e.Finisher.Finish(writeTo, r.context)
			if err != nil {
				return ast.WalkStop, err
			}
		}
		_, _ = bs.Current().Block.WriteString(e.Exiting)
	}

	return ast.WalkContinue, nil
}

func isChild(node ast.Node) bool {
	if node.Parent() != nil && node.Parent().Kind() == ast.KindBlockquote {
		// skip paragraph within blockquote to avoid reflowing text
		return true
	}
	for n := node.Parent(); n != nil; n = n.Parent() {
		// These types are already rendered by their parent
		switch n.Kind() {
		case ast.KindLink, ast.KindImage, ast.KindEmphasis, astext.KindStrikethrough, astext.KindTableCell:
			return true
		}
	}

	return false
}

func resolveRelativeURL(baseURL, rel string) string {
	u, err := url.Parse(rel)
	if err != nil || u.IsAbs() {
		return rel
	}

	u.Path = strings.TrimPrefix(u.Path, "/")

	base, err := url.Parse(baseURL)
	if err != nil {
		return rel
	}

	return base.ResolveReference(u).String()
}
