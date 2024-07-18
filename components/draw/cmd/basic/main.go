package main

import (
	"fmt"
	"math"

	drawille "github.com/rprtr258/tea/components/draw"
)

const RAD = math.Pi / 180

func main() {
	c := drawille.NewCanvas()

	for x := 0; x < 1800; x++ {
		c.Set(float64(x)/10, 10*math.Sin(float64(x)*RAD))
	}
	fmt.Print(c)

	c.Clear()

	for x := 0; x < 1800; x += 10 {
		c.Set(float64(x)/10, 10+10*math.Sin(float64(x)*RAD))
		c.Set(float64(x)/10, 10+10*math.Cos(float64(x)*RAD))
	}
	fmt.Print(c)

	c.Clear()

	for x := 0; x < 3600; x += 20 {
		c.Set(float64(x)/20, 4+4*math.Sin(float64(x)*RAD))
	}
	fmt.Print(c)

	c.Clear()

	for x := 0; x < 360; x += 4 {
		c.Set(float64(x)/4, 30+30*math.Sin(float64(x)*RAD))
	}

	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			c.SetN(x, y)
			c.Toggle(x+30, y+30)
			c.Toggle(x+60, y)
		}
	}
	fmt.Print(c)
}
