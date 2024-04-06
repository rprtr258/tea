package main

import (
	"encoding/json"
	"log"
	"os"

	// "github.com/rprtr258/tea/components/markdown"
	"github.com/rprtr258/tea/components/markdown/ansi"
)

func writeStyleJSON(filename string, styleConfig ansi.StyleConfig) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(styleConfig)
}

func run() error {
	// TODO: what for?
	// for style, styleConfig := range markdown.DefaultStyles {
	// 	if err := writeStyleJSON(filepath.Join("styles", style+".json"), styleConfig); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
