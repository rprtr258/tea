package reactive

import (
	"testing"

	"github.com/rprtr258/assert"
)

func remap(
	x float64,
	x0, x1 float64,
	y0, y1 float64,
) float64 {
	return y0 + (x-x0)*(y1-y0)/(x1-x0)
}

// https://vueuse.org/math/createProjection/
func TestProjection(t *testing.T) {
	useProjector := func(x float64) float64 {
		return remap(x, 0, 10, 0, 100)
	}
	input := Stored[float64](0)
	projected := Computed(input, useProjector)

	assert.Equal(t, 0, projected.Get())

	input.Set(5)
	assert.Equal(t, 50, projected.Get())

	input.Set(10)
	assert.Equal(t, 100, projected.Get())
}
