package ansi

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/rprtr258/fun"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
)

func textFromChildren(node ast.Node, source []byte) string {
	var sb strings.Builder
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindText {
			cn := c.(*ast.Text)
			sb.Write(cn.Segment.Value(source))

			if cn.HardLineBreak() || cn.SoftLineBreak() {
				sb.WriteRune('\n')
			}
		} else {
			sb.Write(c.Text(source))
		}
	}

	return sb.String()
}

// ElementRenderer is called when entering a markdown node.
type ElementRenderer interface {
	Render(w io.Writer, ctx RenderContext) error
}

// ElementFinisher is called when leaving a markdown node.
type ElementFinisher interface {
	Finish(w io.Writer, ctx RenderContext) error
}

// An Element is used to instruct the renderer how to handle individual markdown
// nodes.
type Element struct {
	Entering string
	Exiting  string
	Renderer ElementRenderer
	Finisher ElementFinisher
}

// NewElement returns the appropriate render Element for a given node.
func (r *Renderer) NewElement(node ast.Node, source []byte) Element {
	ctx := r.context

	switch node.Kind() {
	// Document
	case ast.KindDocument:
		e := &BlockElement{
			Block:  &bytes.Buffer{},
			Style:  ctx.options.Styles.Document,
			Margin: true,
		}
		return Element{
			Renderer: e,
			Finisher: e,
		}
	// Heading
	case ast.KindHeading:
		n := node.(*ast.Heading)
		he := &HeadingElement{
			Level: n.Level,
			First: node.PreviousSibling() == nil,
		}
		return Element{
			Exiting:  "",
			Renderer: he,
			Finisher: he,
		}
	// Paragraph
	case ast.KindParagraph:
		if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem {
			return Element{}
		}

		return Element{
			Renderer: &ParagraphElement{
				First: node.PreviousSibling() == nil,
			},
			Finisher: &ParagraphElement{},
		}
	// Blockquote
	case ast.KindBlockquote:
		e := &BlockElement{
			Block:   &bytes.Buffer{},
			Style:   cascadeStyle(ctx.blockStack.Current().Style, ctx.options.Styles.BlockQuote, false),
			Margin:  true,
			Newline: true,
		}
		return Element{
			Entering: "\n",
			Renderer: e,
			Finisher: e,
		}
	// Lists
	case ast.KindList:
		s := ctx.options.Styles.List.StyleBlock
		if s.Indent == nil {
			s.Indent = new(uint)
		}

		n := node.Parent()
		for n != nil {
			if n.Kind() == ast.KindList {
				i := ctx.options.Styles.List.LevelIndent
				s.Indent = &i
				break
			}
			n = n.Parent()
		}

		e := &BlockElement{
			Block:   &bytes.Buffer{},
			Style:   cascadeStyle(ctx.blockStack.Current().Style, s, false),
			Margin:  true,
			Newline: true,
		}
		return Element{
			Entering: "\n",
			Renderer: e,
			Finisher: e,
		}
	case ast.KindListItem:
		l := uint(1)
		for n := node; n.PreviousSibling() != nil && n.PreviousSibling().Kind() == ast.KindListItem; n = n.PreviousSibling() {
			l++
		}

		var e uint
		if p := node.Parent().(*ast.List); p.IsOrdered() {
			e = l
			if p.Start != 1 {
				e += uint(p.Start) - 1
			}
		}

		post := fun.IF(
			node.LastChild() != nil && node.LastChild().Kind() == ast.KindList ||
				node.NextSibling() == nil,
			"", "\n")

		if node.FirstChild() != nil &&
			node.FirstChild().FirstChild() != nil &&
			node.FirstChild().FirstChild().Kind() == astext.KindTaskCheckBox {
			nc := node.FirstChild().FirstChild().(*astext.TaskCheckBox)

			return Element{
				Exiting: post,
				Renderer: &TaskElement{
					Checked: nc.IsChecked,
				},
			}
		}

		return Element{
			Exiting: post,
			Renderer: &ItemElement{
				IsOrdered:   node.Parent().(*ast.List).IsOrdered(),
				Enumeration: e,
			},
		}
	// Text Elements
	case ast.KindText:
		n := node.(*ast.Text)

		s := string(n.Segment.Value(source))
		if n.HardLineBreak() || n.SoftLineBreak() {
			s += "\n"
		}

		return Element{
			Renderer: &BaseElement{
				Token: html.UnescapeString(s),
				Style: ctx.options.Styles.Text,
			},
		}

	case ast.KindEmphasis:
		n := node.(*ast.Emphasis)

		return Element{
			Renderer: &BaseElement{
				Token: html.UnescapeString(string(n.Text(source))),
				Style: fun.IF(n.Level > 1,
					ctx.options.Styles.Strong,
					ctx.options.Styles.Emph),
			},
		}
	case astext.KindStrikethrough:
		n := node.(*astext.Strikethrough)
		return Element{
			Renderer: &BaseElement{
				Token: html.UnescapeString(string(n.Text(source))),
				Style: ctx.options.Styles.Strikethrough,
			},
		}
	case ast.KindThematicBreak:
		return Element{
			Entering: "",
			Exiting:  "",
			Renderer: &BaseElement{
				Style: ctx.options.Styles.HorizontalRule,
			},
		}
	// Links
	case ast.KindLink:
		n := node.(*ast.Link)
		return Element{
			Renderer: &LinkElement{
				Text:    textFromChildren(node, source),
				BaseURL: ctx.options.BaseURL,
				URL:     string(n.Destination),
			},
		}
	case ast.KindAutoLink:
		n := node.(*ast.AutoLink)
		u := string(n.URL(source))
		if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(u), "mailto:") {
			u = "mailto:" + u
		}

		return Element{
			Renderer: &LinkElement{
				Text:    string(n.Label(source)),
				BaseURL: ctx.options.BaseURL,
				URL:     u,
			},
		}
	// Images
	case ast.KindImage:
		n := node.(*ast.Image)
		return Element{
			Renderer: &ImageElement{
				Text:    string(n.Text(source)),
				BaseURL: ctx.options.BaseURL,
				URL:     string(n.Destination),
			},
		}
	// Code
	case ast.KindFencedCodeBlock:
		n := node.(*ast.FencedCodeBlock)

		s := ""
		for i := 0; i < n.Lines().Len(); i++ {
			s += string(fun.Ptr(n.Lines().At(i)).Value(source))
		}

		return Element{
			Entering: "\n",
			Renderer: &CodeBlockElement{
				Code:     s,
				Language: string(n.Language(source)),
			},
		}

	case ast.KindCodeBlock:
		n := node.(*ast.CodeBlock)
		l := n.Lines().Len()
		s := ""
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			s += string(line.Value(source))
		}
		return Element{
			Entering: "\n",
			Renderer: &CodeBlockElement{
				Code: s,
			},
		}

	case ast.KindCodeSpan:
		// n := node.(*ast.CodeSpan)
		e := &BlockElement{
			Block: &bytes.Buffer{},
			Style: cascadeStyle(ctx.blockStack.Current().Style, ctx.options.Styles.Code, false),
		}
		return Element{
			Renderer: e,
			Finisher: e,
		}

	// Tables
	case astext.KindTable:
		te := &TableElement{}
		return Element{
			Entering: "\n",
			Renderer: te,
			Finisher: te,
		}

	case astext.KindTableCell:
		s := ""
		n := node.FirstChild()
		for n != nil {
			s += string(n.Text(source))
			// s += string(n.LinkData.Destination)
			n = n.NextSibling()
		}

		return Element{
			Renderer: &TableCellElement{
				Text: s,
				Head: node.Parent().Kind() == astext.KindTableHeader,
			},
		}

	case astext.KindTableHeader:
		return Element{
			Finisher: &TableHeadElement{},
		}
	case astext.KindTableRow:
		return Element{
			Finisher: &TableRowElement{},
		}

	// HTML Elements
	case ast.KindHTMLBlock:
		n := node.(*ast.HTMLBlock)
		return Element{
			Renderer: &BaseElement{
				Token: ctx.SanitizeHTML(string(n.Text(source)), true),
				Style: ctx.options.Styles.HTMLBlock.StylePrimitive,
			},
		}
	case ast.KindRawHTML:
		n := node.(*ast.RawHTML)
		return Element{
			Renderer: &BaseElement{
				Token: ctx.SanitizeHTML(string(n.Text(source)), true),
				Style: ctx.options.Styles.HTMLSpan.StylePrimitive,
			},
		}
	// Definition Lists
	case astext.KindDefinitionList:
		e := &BlockElement{
			Block:   &bytes.Buffer{},
			Style:   cascadeStyle(ctx.blockStack.Current().Style, ctx.options.Styles.DefinitionList, false),
			Margin:  true,
			Newline: true,
		}
		return Element{
			Entering: "\n",
			Renderer: e,
			Finisher: e,
		}
	case astext.KindDefinitionTerm:
		return Element{
			Renderer: &BaseElement{
				Style: ctx.options.Styles.DefinitionTerm,
			},
		}
	case astext.KindDefinitionDescription:
		return Element{
			Renderer: &BaseElement{
				Style: ctx.options.Styles.DefinitionDescription,
			},
		}
	// Handled by parents
	case astext.KindTaskCheckBox:
		// handled by KindListItem
		return Element{}
	case ast.KindTextBlock:
		return Element{}
	case east.KindEmoji:
		n := node.(*east.Emoji)
		return Element{
			Renderer: &BaseElement{
				Token: string(n.Value.Unicode),
			},
		}
	// Unknown case
	default:
		fmt.Println("Warning: unhandled element", node.Kind().String())
		return Element{}
	}
}
