package helloworld

import (
	"fmt"

	"github.com/rprtr258/tea/glamour"
)

func Main() {
	in := `# Hello World

This is a simple example of Markdown rendering with Glamour!
Check out the [other examples](https://github.com/rprtr258/tea/glamour/tree/master/examples) too.

Bye!
`

	out, _ := glamour.Render(in, "dark")
	fmt.Print(out)
}
