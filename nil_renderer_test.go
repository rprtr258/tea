package tea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilRenderer(t *testing.T) {
	r := nilRenderer{}
	r.start()
	r.stop()
	r.kill()
	r.Write([]byte("a"))
	r.repaint()
	r.enterAltScreen()
	assert.False(t, r.altScreen())
	r.exitAltScreen()
	r.clearScreen()
	r.setCursor(true)
	r.setCursor(false)
	r.setMouseCellMotion(true)
	r.setMouseCellMotion(false)
	r.setMouseAllMotion(true)
	r.setMouseAllMotion(false)
}
