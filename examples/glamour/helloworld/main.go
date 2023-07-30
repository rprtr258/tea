package helloworld

import (
	"context"
	"fmt"

	"github.com/rprtr258/tea/glamour"
)

func Main(context.Context) error {
	in := `# Hello World

This is a simple example of Markdown rendering with Glamour!
Check out the [other examples](https://github.com/rprtr258/tea/glamour/tree/master/examples) too.

` + "```" + `go
package main

import main

func main() {
	fmt.Println("Hello World!")
}
` + "```" + `

Bye!
`

	out, _ := glamour.Render(in, "dark")
	fmt.Print(out)
	return nil
}
