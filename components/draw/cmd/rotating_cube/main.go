package main

import (
	"math"
	"os"
	"strings"
	"time"

	tg "github.com/nsf/termbox-go" // TODO: remove
	. "github.com/rprtr258/tea/components/draw"
)

const RAD = math.Pi / 180

func round(x float64) int {
	return int(x + 0.5)
}

type Point3D struct {
	x, y, z float64
}

// Rotates the point around the X axis by the given angle in degrees.
func (p Point3D) RotateX(angle float64) Point3D {
	rad := RAD * angle
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	y := p.y*cosa - p.z*sina
	z := p.y*sina + p.z*cosa
	return Point3D{p.x, y, z}
}

// Rotates the point around the Y axis by the given angle in degrees.
func (p Point3D) RotateY(angle float64) Point3D {
	rad := RAD * angle
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	z := p.z*cosa - p.x*sina
	x := p.z*sina + p.x*cosa
	return Point3D{x, p.y, z}
}

// Rotates the point around the Z axis by the given angle in degrees.
func (p Point3D) RotateZ(angle float64) Point3D {
	rad := RAD * angle
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	x := p.x*cosa - p.y*sina
	y := p.x*sina + p.y*cosa
	return Point3D{x, y, p.z}
}

// Transforms this 3D point to 2D using a perspective projection.
func (p *Point3D) Project(win_width, win_height, fov, viewer_distance float64) Point3D {
	factor := fov / (viewer_distance + p.z)
	x := p.x*factor + win_width/2
	y := -p.y*factor + win_height/2
	return Point3D{x, y, 1}
}

func drawFrame(
	c Canvas,
	angleX, angleY, angleZ *float64,
	projection bool,
) {
	vertices := [...]Point3D{
		{-20, +20, -20},
		{+20, +20, -20},
		{+20, -20, -20},
		{-20, -20, -20},
		{-20, +20, +20},
		{+20, +20, +20},
		{+20, -20, +20},
		{-20, -20, +20},
	}

	faces := [...][4]int{
		{0, 1, 2, 3},
		{1, 5, 6, 2},
		{5, 4, 7, 6},
		{4, 0, 3, 7},
		{0, 4, 5, 1},
		{3, 2, 6, 7},
	}

	ts := [len(vertices)]Point3D{} // transformed vertices
	for i, v := range vertices {
		// Rotate the point around X axis, then around Y axis, and finally around Z axis.
		point := v.RotateX(*angleX).RotateY(*angleY).RotateZ(*angleZ)
		if projection {
			// Transform the point from 3D to 2D
			point = point.Project(50, 50, 50, 50)
		}

		ts[i] = point
	}

	for _, f := range faces {
		Line(round(ts[f[0]].x), round(ts[f[0]].y), round(ts[f[1]].x), round(ts[f[1]].y))(c.SetN)
		Line(round(ts[f[1]].x), round(ts[f[1]].y), round(ts[f[2]].x), round(ts[f[2]].y))(c.SetN)
		Line(round(ts[f[2]].x), round(ts[f[2]].y), round(ts[f[3]].x), round(ts[f[3]].y))(c.SetN)
		Line(round(ts[f[3]].x), round(ts[f[3]].y), round(ts[f[0]].x), round(ts[f[0]].y))(c.SetN)
	}

	f := c.Frame(-40, -40, 80, 80)

	const xoffset = 2
	for y, line := range strings.Split(f, "\n") {
		pos := 0
		for _, r := range line { // iterates over runes, not positions
			tg.SetCell(xoffset+pos, y, r, tg.ColorRed|tg.AttrBold, tg.ColorBlack|tg.AttrBold)
			pos++
		}
	}

	tg.Flush()

	*angleX += 2.0
	*angleY += 3.0
	*angleZ += 5.0

	c.Clear()
}

func run(projection bool) {
	tg.Clear(tg.ColorRed|tg.AttrBold, tg.ColorBlack|tg.AttrBold)

	eventQueue := make(chan tg.Event)
	go func() {
		for {
			eventQueue <- tg.PollEvent()
		}
	}()
	drawTick := time.NewTicker(50 * time.Millisecond)

	angleX, angleY, angleZ := 0.0, 0.0, 0.0
	c := NewCanvas()
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == tg.EventKey && ev.Key == tg.KeyCtrlC {
				return
			}
			drawFrame(c, &angleX, &angleY, &angleZ, projection)
		case <-drawTick.C:
			drawFrame(c, &angleX, &angleY, &angleZ, projection)
		}
	}
}

func main() {
	tg.Init()
	defer tg.Close()

	projection := len(os.Args) > 1 && os.Args[1] == "-p"
	run(projection)
}
