package ansi

import (
	"io"
	"net/url"
)

// A LinkElement is used to render hyperlinks.
type LinkElement struct {
	Text    string
	BaseURL string
	URL     string
	Child   ElementRenderer
}

func (e *LinkElement) Render(w io.Writer, ctx RenderContext) error {
	textRendered := e.Text != "" && e.Text != e.URL
	if textRendered {
		el := &BaseElement{
			Token: e.Text,
			Style: ctx.options.Styles.LinkText,
		}

		if err := el.Render(w, ctx); err != nil {
			return err
		}
	}

	/*
		if node.LastChild != nil {
			if node.LastChild.Type == bf.Image {
				el := tr.NewElement(node.LastChild)
				err := el.Renderer.Render(w, node.LastChild, tr)
				if err != nil {
					return err
				}
			}
			if len(node.LastChild.Literal) > 0 &&
				string(node.LastChild.Literal) != string(node.LinkData.Destination) {
				textRendered = true
				el := &BaseElement{
					Token: string(node.LastChild.Literal),
					Style: ctx.style[LinkText],
				}
				err := el.Render(w, node.LastChild, tr)
				if err != nil {
					return err
				}
			}
		}
	*/

	if u, err := url.Parse(e.URL); err == nil &&
		"#"+u.Fragment != e.URL { // if the URL only consists of an anchor, ignore it
		pre := " "
		style := ctx.options.Styles.Link
		if !textRendered {
			pre = ""
			style.BlockPrefix = ""
			style.BlockSuffix = ""
		}

		el := &BaseElement{
			Token:  resolveRelativeURL(e.BaseURL, e.URL),
			Prefix: pre,
			Style:  style,
		}

		if err := el.Render(w, ctx); err != nil {
			return err
		}
	}

	return nil
}
