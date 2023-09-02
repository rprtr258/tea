package ansi

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const (
	generateExamples = false
	generateIssues   = false
	examplesDir      = "../styles/examples/"
	issuesDir        = "../testdata/issues/"
)

func TestRenderer(t *testing.T) {
	files, err := filepath.Glob(examplesDir + "*.md")
	assert.NoError(t, err)

	for _, f := range files {
		bn := strings.TrimSuffix(filepath.Base(f), ".md")
		sn := filepath.Join(examplesDir, bn+".style")
		tn := filepath.Join("../testdata", bn+".test")

		in, err := os.ReadFile(f)
		assert.NoError(t, err)

		b, err := os.ReadFile(sn)
		assert.NoError(t, err)

		options := Options{
			WordWrap:     80,
			ColorProfile: termenv.TrueColor,
		}
		assert.NoError(t, json.Unmarshal(b, &options.Styles))

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
		if generateExamples {
			err = os.WriteFile(tn, buf.Bytes(), 0o644)
			assert.NoError(t, err)
			continue
		}

		// verify
		td, err := os.ReadFile(tn)
		assert.NoError(t, err)

		assert.Equal(t, td, buf.Bytes())
	}
}

func TestRendererIssues(t *testing.T) {
	files, err := filepath.Glob(issuesDir + "*.md")
	assert.NoError(t, err)

	for _, f := range files {
		bn := strings.TrimSuffix(filepath.Base(f), ".md")
		t.Run(bn, func(t *testing.T) {
			tn := filepath.Join(issuesDir, bn+".test")

			in, err := os.ReadFile(f)
			assert.NoError(t, err)

			b, err := os.ReadFile("../styles/dark.json")
			assert.NoError(t, err)

			options := Options{
				WordWrap:     80,
				ColorProfile: termenv.TrueColor,
			}
			assert.NoError(t, json.Unmarshal(b, &options.Styles))

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
			if generateIssues {
				assert.NoError(t, os.WriteFile(tn, buf.Bytes(), 0o644))
				return
			}

			// verify
			td, err := os.ReadFile(tn)
			assert.NoError(t, err)

			assert.Equal(t, td, buf.Bytes())
		})
	}
}
