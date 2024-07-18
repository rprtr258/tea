package draw

import (
	"fmt"
	"math"
)

func ternary[T any](p bool, a, b T) T {
	if p {
		return a
	}
	return b
}

func radians(d float64) float64 {
	return d * math.Pi / 180
}

func Line(x0, y0, x1, y1 int) func(func(x, y int)) {
	dx := abs(x1 - x0)
	sx := ternary(x0 < x1, 1, -1)
	dy := -abs(y1 - y0)
	sy := ternary(y0 < y1, 1, -1)

	return func(yield func(x, y int)) {
		for err := dx + dy; ; {
			yield(x0, y0)
			if x0 == x1 && y0 == y1 {
				break
			}
			e2 := 2 * err
			if e2 >= dy {
				if x0 == x1 {
					break
				}
				err += dy
				x0 += sx
			}
			if e2 <= dx {
				if y0 == y1 {
					break
				}
				err += dx
				y0 += sy
			}
		}
	}
}

// TODO: convert to test
func init() {
	for _, tc := range []struct {
		x float64
		n int
	}{
		{-1.2246467991473515e-14, 0},
		{178.5, 178},
		{179.5, 180},
		{-3.2556815445715848, -3},
	} {
		if round(tc.x) != tc.n {
			panic(fmt.Sprint("round error", tc.x, tc.n, round(tc.x)))
		}
	}
}

func round(x float64) int {
	y := math.Floor(x + 0.5)
	if math.Mod(y, 2) == 1 && math.Mod(x, 1) == 0.5 {
		y--
	}
	return int(y)
}

func Polygon(center_x, center_y, sides, radius float64) func(func(x, y int)) {
	degree := 360 / sides
	return func(yield func(x, y int)) {
		for n := 0; n < int(sides); n++ {
			a := float64(n) * degree
			b := float64(n+1) * degree

			x1 := center_x + math.Cos(radians(a))*(radius/2+1)
			y1 := center_y + math.Sin(radians(a))*(radius/2+1)
			x2 := center_x + math.Cos(radians(b))*(radius/2+1)
			y2 := center_y + math.Sin(radians(b))*(radius/2+1)

			Line(round(x1), round(y1), round(x2), round(y2))(yield)
		}
	}
}
