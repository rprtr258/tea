package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/rprtr258/tea/components/markdown"
	"github.com/rprtr258/tea/components/markdown/ansi"
)

func writeStyleJSON(filename string, styleConfig *ansi.StyleConfig) error {
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
	for style, styleConfig := range markdown.DefaultStyles {
		filename := filepath.Join("styles", style+".json")
		if err := writeStyleJSON(filename, styleConfig); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
