package ansi

import (
	"bytes"
	"io"
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

// A ParagraphElement is used to render individual paragraphs.
type ParagraphElement struct {
	First bool
}

func (e *ParagraphElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := ctx.options.Styles.Paragraph

	if !e.First {
		_, _ = w.Write([]byte("\n"))
	}
	be := BlockElement{
		Block: &bytes.Buffer{},
		Style: cascadeStyle(bs.Current().Style, rules, false),
	}
	bs.Push(be)

	renderText(w, bs.Parent().Style.StylePrimitive, rules.BlockPrefix)
	renderText(bs.Current().Block, bs.Current().Style.StylePrimitive, rules.Prefix)
	return nil
}

func (e *ParagraphElement) Finish(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := bs.Current().Style

	mw := NewMarginWriter(ctx, w, rules)
	if strings.TrimSpace(bs.Current().Block.String()) != "" {
		flow := wordwrap.NewWriter(int(bs.Width(ctx)))
		flow.KeepNewlines = ctx.options.PreserveNewLines
		_, _ = flow.Write(bs.Current().Block.Bytes())
		flow.Close()

		_, err := mw.Write(flow.Bytes())
		if err != nil {
			return err
		}
		_, _ = mw.Write([]byte("\n"))
	}

	renderText(w, bs.Current().Style.StylePrimitive, rules.Suffix)
	renderText(w, bs.Parent().Style.StylePrimitive, rules.BlockSuffix)

	bs.Current().Block.Reset()
	bs.Pop()
	return nil
}
