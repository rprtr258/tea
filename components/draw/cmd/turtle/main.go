package main

import (
	"fmt"
	"math"

	drawille "github.com/rprtr258/tea/components/draw"
)

const RAD = math.Pi / 180

func main() {
	t := drawille.NewTurtle()

	for range [36]struct{}{} {
		t.Right(10)
		for range [36]struct{}{} {
			t.Right(10)
			t.Forward(8)
		}
	}
	fmt.Print(t)
}
