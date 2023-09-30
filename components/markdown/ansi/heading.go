package ansi

import (
	"bytes"
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
)

// A HeadingElement is used to render headings.
type HeadingElement struct {
	Level int
	First bool
}

func (e *HeadingElement) Render(w io.Writer, ctx RenderContext) error {
	rules := ctx.options.Styles.Heading

	if e.Level >= 1 && e.Level <= 6 {
		styleBlock := map[int]StyleBlock{
			1: ctx.options.Styles.H1,
			2: ctx.options.Styles.H2,
			3: ctx.options.Styles.H3,
			4: ctx.options.Styles.H4,
			5: ctx.options.Styles.H5,
			6: ctx.options.Styles.H6,
		}[e.Level]
		rules = cascadeStyles(true, rules, styleBlock)
	}

	bs := ctx.blockStack
	if !e.First {
		renderText(w, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, "\n")
	}

	bs.Push(BlockElement{
		Block: &bytes.Buffer{},
		Style: cascadeStyle(bs.Current().Style, rules, false),
	})

	renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, rules.BlockPrefix)
	renderText(bs.Current().Block, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, rules.Prefix)
	return nil
}

func (e *HeadingElement) Finish(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := bs.Current().Style

	var indentation uint
	if rules.Indent != nil {
		indentation = *rules.Indent
	}

	var margin uint
	if rules.Margin != nil {
		margin = *rules.Margin
	}

	iw := indent.NewWriterPipe(w, indentation+margin, func(wr io.Writer) {
		renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, " ")
	})

	flow := wordwrap.NewWriter(int(bs.Width(ctx) - indentation - margin*2))
	if _, err := flow.Write(bs.Current().Block.Bytes()); err != nil {
		return err
	}
	flow.Close()

	if _, err := iw.Write(flow.Bytes()); err != nil {
		return err
	}

	renderText(w, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, rules.Suffix)
	renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, rules.BlockSuffix)

	bs.Current().Block.Reset()
	bs.Pop()
	return nil
}
