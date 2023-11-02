package styles

import (
	"os"
	"testing"

	"github.com/rprtr258/assert"
)

func TestRendererHasDarkBackground(t *testing.T) {
	r1 := NewRenderer(os.Stdout)
	r1.SetHasDarkBackground(false)
	assert.False(t, r1.HasDarkBackground())

	r2 := NewRenderer(os.Stdout)
	r2.SetHasDarkBackground(true)
	assert.True(t, r2.HasDarkBackground())
}
