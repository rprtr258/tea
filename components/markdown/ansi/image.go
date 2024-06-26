package ansi

import "io"

// An ImageElement is used to render images elements.
type ImageElement struct {
	Text    string
	BaseURL string
	URL     string
	Child   ElementRenderer
}

func (e *ImageElement) Render(w io.Writer, ctx RenderContext) error {
	if e.Text != "" {
		el := &BaseElement{
			Token: e.Text,
			Style: ctx.options.Styles.ImageText,
		}

		if err := el.Render(w, ctx); err != nil {
			return err
		}
	}

	if e.URL != "" {
		el := &BaseElement{
			Token:  resolveRelativeURL(e.BaseURL, e.URL),
			Prefix: " ",
			Style:  ctx.options.Styles.Image,
		}

		if err := el.Render(w, ctx); err != nil {
			return err
		}
	}

	return nil
}
