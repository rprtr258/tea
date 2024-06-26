package ansi

import (
	"io"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

// chromaStyleTheme name used for rendering.
const chromaStyleTheme = "charm"

// mutex for synchronizing access to the chroma style registry.
// Related https://github.com/alecthomas/chroma/pull/650
var mutex sync.Mutex

// A CodeBlockElement is used to render code blocks.
type CodeBlockElement struct {
	Code     string
	Language string
}

func chromaStyle(style StylePrimitive) string {
	var s string

	if style.ForegroundColor != nil {
		s = *style.ForegroundColor
	}
	if style.BackgroundColor != nil {
		if s != "" {
			s += " "
		}
		s += "bg:" + *style.BackgroundColor
	}
	if style.Italic != nil && *style.Italic {
		if s != "" {
			s += " "
		}
		s += "italic"
	}
	if style.Bold != nil && *style.Bold {
		if s != "" {
			s += " "
		}
		s += "bold"
	}
	if style.Underline != nil && *style.Underline {
		if s != "" {
			s += " "
		}
		s += "underline"
	}

	return s
}

func (e *CodeBlockElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack

	rules := ctx.options.Styles.CodeBlock

	var indentation uint
	if rules.Indent != nil {
		indentation = *rules.Indent
	}

	var margin uint
	if rules.Margin != nil {
		margin = *rules.Margin
	}

	theme := rules.Theme

	if rules.Chroma != nil && ctx.options.ColorProfile != termenv.Ascii {
		theme = chromaStyleTheme
		mutex.Lock()
		// Don't register the style if it's already registered.
		_, ok := styles.Registry[theme]
		if !ok {
			styles.Register(chroma.MustNewStyle(theme, chroma.StyleEntries{
				chroma.Text:                chromaStyle(rules.Chroma.Text),
				chroma.Error:               chromaStyle(rules.Chroma.Error),
				chroma.Comment:             chromaStyle(rules.Chroma.Comment),
				chroma.CommentPreproc:      chromaStyle(rules.Chroma.CommentPreproc),
				chroma.Keyword:             chromaStyle(rules.Chroma.Keyword),
				chroma.KeywordReserved:     chromaStyle(rules.Chroma.KeywordReserved),
				chroma.KeywordNamespace:    chromaStyle(rules.Chroma.KeywordNamespace),
				chroma.KeywordType:         chromaStyle(rules.Chroma.KeywordType),
				chroma.Operator:            chromaStyle(rules.Chroma.Operator),
				chroma.Punctuation:         chromaStyle(rules.Chroma.Punctuation),
				chroma.Name:                chromaStyle(rules.Chroma.Name),
				chroma.NameBuiltin:         chromaStyle(rules.Chroma.NameBuiltin),
				chroma.NameTag:             chromaStyle(rules.Chroma.NameTag),
				chroma.NameAttribute:       chromaStyle(rules.Chroma.NameAttribute),
				chroma.NameClass:           chromaStyle(rules.Chroma.NameClass),
				chroma.NameConstant:        chromaStyle(rules.Chroma.NameConstant),
				chroma.NameDecorator:       chromaStyle(rules.Chroma.NameDecorator),
				chroma.NameException:       chromaStyle(rules.Chroma.NameException),
				chroma.NameFunction:        chromaStyle(rules.Chroma.NameFunction),
				chroma.NameOther:           chromaStyle(rules.Chroma.NameOther),
				chroma.Literal:             chromaStyle(rules.Chroma.Literal),
				chroma.LiteralNumber:       chromaStyle(rules.Chroma.LiteralNumber),
				chroma.LiteralDate:         chromaStyle(rules.Chroma.LiteralDate),
				chroma.LiteralString:       chromaStyle(rules.Chroma.LiteralString),
				chroma.LiteralStringEscape: chromaStyle(rules.Chroma.LiteralStringEscape),
				chroma.GenericDeleted:      chromaStyle(rules.Chroma.GenericDeleted),
				chroma.GenericEmph:         chromaStyle(rules.Chroma.GenericEmph),
				chroma.GenericInserted:     chromaStyle(rules.Chroma.GenericInserted),
				chroma.GenericStrong:       chromaStyle(rules.Chroma.GenericStrong),
				chroma.GenericSubheading:   chromaStyle(rules.Chroma.GenericSubheading),
				chroma.Background:          chromaStyle(rules.Chroma.Background),
			}))
		}
		mutex.Unlock()
	}

	iw := indent.NewWriterPipe(w, indentation+margin, func(wr io.Writer) {
		renderText(w, bs.Current().Style.StylePrimitive, " ")
	})

	if theme != "" {
		renderText(iw, bs.Current().Style.StylePrimitive, rules.BlockPrefix)
		err := quick.Highlight(iw, e.Code, e.Language, "terminal256", theme)
		if err != nil {
			return err
		}
		renderText(iw, bs.Current().Style.StylePrimitive, rules.BlockSuffix)
		return nil
	}

	// fallback rendering
	el := &BaseElement{
		Token: e.Code,
		Style: rules.StylePrimitive,
	}

	return el.Render(iw, ctx)
}
