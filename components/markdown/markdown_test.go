package markdown

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rprtr258/assert"
)

const (
	generate = false
	markdown = "testdata/readme.markdown.in"
	testFile = "testdata/readme.test"
)

func TestTermRendererWriter(t *testing.T) {
	r, err := NewTermRenderer(
		WithStyles(DarkStyle),
	)
	assert.NoError(t, err)

	in := assert.UseFileContent(t, markdown)

	_, err = r.Write(in)
	assert.NoError(t, err)

	assert.NoError(t, r.Close())

	b, err := io.ReadAll(r)
	assert.NoError(t, err)

	// generate
	if generate {
		assert.NoError(t, os.WriteFile(testFile, b, 0o644))
		return
	}

	// verify
	td := assert.UseFileContent(t, testFile)

	assert.Equal(t, td, b)
}

func TestTermRenderer(t *testing.T) {
	r, err := NewTermRenderer(WithStyles(DarkStyle))
	assert.NoError(t, err)

	in := assert.UseFileContent(t, markdown)

	b, err := r.Render(string(in))
	assert.NoError(t, err)

	// verify
	td := assert.UseFileContent(t, testFile)

	assert.Equal(t, string(td), b)
}

func TestWithEmoji(t *testing.T) {
	r, err := NewTermRenderer(WithEmoji())
	assert.NoError(t, err)

	b, err := r.Render(":+1:")
	assert.NoError(t, err)

	b = strings.TrimSpace(b)

	// Thumbs up unicode character
	td := "\U0001f44d"

	assert.Equal(t, td, b)
}

func TestWithPreservedNewLines(t *testing.T) {
	r, err := NewTermRenderer(
		WithPreservedNewLines(),
	)
	assert.NoError(t, err)

	in := assert.UseFileContent(t, "testdata/preserved_newline.in")

	b, err := r.Render(string(in))
	assert.NoError(t, err)

	// verify
	td := assert.UseFileContent(t, "testdata/preserved_newline.test")

	assert.Equal(t, string(td), b)
}

func TestStyles(t *testing.T) {
	_, err := NewTermRenderer(WithAutoStyle())
	assert.NoError(t, err)

	_, err = NewTermRenderer(WithAutoStyle())
	assert.NoError(t, err)
}

func TestRenderHelpers(t *testing.T) {
	in := assert.UseFileContent(t, markdown)

	b, err := Render(string(in), DarkStyle)
	assert.NoError(t, err)

	// verify
	td := assert.UseFileContent(t, testFile)

	assert.Equal(t, string(td), b)
}

func TestCapitalization(t *testing.T) {
	p := true
	style := DarkStyle
	style.H1.Upper = &p
	style.H2.Title = &p
	style.H3.Lower = &p

	r, err := NewTermRenderer(WithStyles(style))
	assert.NoError(t, err)

	b, err := r.Render("# everything is uppercase\n## everything is titled\n### everything is lowercase")
	assert.NoError(t, err)

	// expected outcome
	td := assert.UseFileContent(t, "testdata/capitalization.test")

	assert.Equal(t, string(td), b)
}
