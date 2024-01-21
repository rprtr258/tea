package ansi

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muesli/termenv"
	"github.com/rprtr258/assert"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const (
	_generateExamples = false
	_generateIssues   = false
	_examplesDir      = "../styles/examples/"
	_issuesDir        = "../testdata/issues/"
)

func TestRenderer(t *testing.T) {
	files, err := filepath.Glob(_examplesDir + "*.md")
	assert.NoError(t, err)

	for _, f := range files {
		bn := strings.TrimSuffix(filepath.Base(f), ".md")
		sn := filepath.Join(_examplesDir, bn+".style")
		tn := filepath.Join("..", "testdata", bn+".test")

		in := assert.UseFileContent(t, f)

		b := assert.UseFileContent(t, sn)

		options := Options{
			WordWrap:     80,
			ColorProfile: termenv.TrueColor,
			Styles:       assert.UseJSON[StyleConfig](t, b),
		}

		md := goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.DefinitionList,
				emoji.Emoji,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)

		ar := NewRenderer(options)
		md.SetRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(util.Prioritized(ar, 1000))))

		var buf bytes.Buffer
		assert.NoError(t, md.Convert(in, &buf))

		// generate
		if _generateExamples {
			err = os.WriteFile(tn, buf.Bytes(), 0o644)
			assert.NoError(t, err)
			continue
		}

		// verify
		td := assert.UseFileContent(t, tn)

		assert.Equal(t, td, buf.Bytes())
	}
}

func TestRendererIssues(t *testing.T) {
	files, err := filepath.Glob(_issuesDir + "*.md")
	assert.NoError(t, err)

	for _, f := range files {
		bn := strings.TrimSuffix(filepath.Base(f), ".md")
		t.Run(bn, func(t *testing.T) {
			tn := filepath.Join(_issuesDir, bn+".test")

			in := assert.UseFileContent(t, f)

			b := assert.UseFileContent(t, "../styles/dark.json")

			options := Options{
				WordWrap:     80,
				ColorProfile: termenv.TrueColor,
				Styles:       assert.UseJSON[StyleConfig](t, b),
			}

			md := goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
					extension.DefinitionList,
					emoji.Emoji,
				),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
			)

			ar := NewRenderer(options)
			md.SetRenderer(
				renderer.NewRenderer(
					renderer.WithNodeRenderers(util.Prioritized(ar, 1000))))

			var buf bytes.Buffer
			assert.NoError(t, md.Convert(in, &buf))

			// generate
			if _generateIssues {
				assert.NoError(t, os.WriteFile(tn, buf.Bytes(), 0o644))
				return
			}

			// verify
			actual := assert.UseFileContent(t, tn)

			assert.Equal(t, string(actual), buf.String())
		})
	}
}
