package draw

import (
	"fmt"
	"math"
)

/*
http://www.alanwood.net/unicode/braille_patterns.html

dots:

	+---+
	|1 4|
	|2 5|
	|3 6|
	|7 8|
	+---+
*/
var pixel_map = [4][2]int{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

// braille unicode characters starts at 0x2800
var braille_char_offset = 0x2800

// TODO: convert to test
func init() {
	for _, tc := range []struct {
		xx, yy int
		x, y   int
	}{
		{178, -3, 89, -1},
	} {
		x, y := get_pos(tc.xx, tc.yy)
		if x != tc.x || y != tc.y {
			panic(fmt.Sprint("get_pos error", tc, x, y))
		}
	}
}

// Convert x,y to cols, rows
func get_pos(x, y int) (int, int) {
	return int(math.Floor(float64(x) / 2)), int(math.Floor(float64(y) / 4))
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func getPixel(y, x int) int {
	cy := mod(y, 4)
	cx := mod(x, 2)
	return pixel_map[cy][cx]
}

type Canvas struct {
	chars map[int]map[int]int
}

// Make a new canvas
func NewCanvas() Canvas {
	return Canvas{
		chars: map[int]map[int]int{},
	}
}

func (c Canvas) MaxY() int {
	res := 0
	for k := range c.chars {
		res = max(res, k)
	}
	return res * 4
}

func (c Canvas) MinY() int {
	res := 0
	for k := range c.chars {
		res = min(res, k)
	}
	return res * 4
}

func (c Canvas) MaxX() int {
	res := 0
	for _, v := range c.chars {
		for k := range v {
			res = max(res, k)
		}
	}
	return res * 2
}

func (c Canvas) MinX() int {
	res := 0
	for _, v := range c.chars {
		for k := range v {
			res = min(res, k)
		}
	}
	return res * 2
}

// Clear all pixels
func (c *Canvas) Clear() {
	for y := range c.chars {
		delete(c.chars, y)
	}
}

// Set a pixel of c
func (c *Canvas) SetN(x, y int) {
	px, py := get_pos(x, y)
	if c.chars[py] == nil {
		c.chars[py] = make(map[int]int)
	}
	c.chars[py][px] |= getPixel(y, x)
}

// Set a pixel of c
func (c *Canvas) Set(x, y float64) {
	c.SetN(round(x), round(y))
}

func abs(x int) int {
	return max(-x, x)
}

// Unset a pixel of c
func (c *Canvas) UnSet(x, y int) {
	px, py := get_pos(x, y)
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]int)
	}
	c.chars[py][px] &^= getPixel(y, x)
}

// Toggle a point
func (c *Canvas) Toggle(x, y int) {
	px, py := get_pos(x, y)
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]int)
	}
	c.chars[py][px] ^= getPixel(y, x)
}

// Set text to the given coordinates
func (c *Canvas) SetText(x, y int, text string) {
	px, py := get_pos(x, y)
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]int)
	}
	for i, char := range text {
		c.chars[py][px+i] = int(char) - braille_char_offset
	}
}

// get pixel at the given coordinates
func (c Canvas) get(x, y int) bool {
	px, py := get_pos(x, y)
	char := c.chars[py][px]
	return (char & pixel_map[y%4][x%2]) != 0
}

// Get character at the given screen coordinates
func (c Canvas) GetScreenCharacter(x, y int) rune {
	return rune(c.chars[y][x] + braille_char_offset)
}

// Get character for the given pixel
func (c Canvas) GetCharacter(x, y int) rune {
	return c.GetScreenCharacter(x/2, y/4)
}

// Retrieve the rows from a given view
func (c Canvas) Rows(minX, minY, maxX, maxY int) func(func(string)) {
	mincol, minrow := get_pos(minX, minY)
	maxcol, maxrow := get_pos(maxX, maxY)

	return func(yield func(string)) {
		for rownum := minrow; rownum <= maxrow; rownum++ {
			row := ""
			for x := mincol; x < (maxcol + 1); x++ {
				char := c.chars[rownum][x]
				row += string(rune(char + braille_char_offset))
			}
			yield(row)
		}
	}
}

// Retrieve a string representation of the frame at the given parameters
func (c Canvas) Frame(minX, minY, maxX, maxY int) string {
	var ret string
	c.Rows(minX, minY, maxX, maxY)(func(row string) {
		ret += row + "\n"
	})
	return ret
}

func (c Canvas) String() string {
	return c.Frame(c.MinX(), c.MinY(), c.MaxX(), c.MaxY())
}
