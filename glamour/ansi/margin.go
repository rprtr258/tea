package ansi

import (
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/padding"
)

// MarginWriter is a Writer that applies indentation and padding around
// whatever you write to it.
type MarginWriter struct {
	w  io.Writer
	pw *padding.Writer
	iw *indent.Writer
}

// NewMarginWriter returns a new MarginWriter.
func NewMarginWriter(ctx RenderContext, w io.Writer, rules StyleBlock) *MarginWriter {
	bs := ctx.blockStack
	pw := padding.NewWriterPipe(w, bs.Width(ctx), func(wr io.Writer) {
		renderText(w, ctx.options.ColorProfile, rules.StylePrimitive, " ")
	})

	ic := " "
	if rules.IndentToken != nil {
		ic = *rules.IndentToken
	}

	var indentation uint
	if rules.Indent != nil {
		indentation = *rules.Indent
	}

	var margin uint
	if rules.Margin != nil {
		margin = *rules.Margin
	}

	return &MarginWriter{
		w:  w,
		pw: pw,
		iw: indent.NewWriterPipe(pw, indentation+margin, func(wr io.Writer) {
			renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, ic)
		}),
	}
}

func (w *MarginWriter) Write(b []byte) (int, error) {
	return w.iw.Write(b)
}
