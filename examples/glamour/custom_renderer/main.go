package custom_renderer

import (
	"fmt"

	"github.com/rprtr258/tea/glamour"
)

func Main() {
	in := `# Custom Renderer

Word-wrapping will occur when lines exceed the limit of 40 characters.
`

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(40),
	)

	out, _ := r.Render(in)
	fmt.Print(out)
}
