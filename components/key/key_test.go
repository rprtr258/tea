package key

import (
	"testing"

	"github.com/rprtr258/assert"
)

func TestBinding_Enabled(t *testing.T) {
	binding := NewBinding(
		WithKeys("k", "up"),
		WithHelp("↑/k", "move up"),
	)
	assert.True(t, binding.Enabled())

	binding.SetEnabled(false)
	assert.False(t, binding.Enabled())

	binding.SetEnabled(true)
	binding.Unbind()
	assert.False(t, binding.Enabled())
}
