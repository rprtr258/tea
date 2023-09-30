package custom_renderer

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea/components/markdown"
)

func Main(context.Context) error {
	in := `# Custom Renderer

Word-wrapping will occur when lines exceed the limit of 40 characters.
`

	r, _ := markdown.NewTermRenderer(
		markdown.WithStandardStyle("dark"),
		markdown.WithWordWrap(40),
	)

	out, err := r.Render(in)
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}
