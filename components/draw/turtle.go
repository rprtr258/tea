package draw

import "math"

// Turtle graphics interface
// http://en.wikipedia.org/wiki/Turtle_graphics
type Turtle struct {
	Canvas
	pos_x, pos_y, rotation float64
	brush_on               bool
}

func NewTurtle() *Turtle {
	return &Turtle{
		Canvas:   NewCanvas(),
		pos_x:    0,
		pos_y:    0,
		rotation: 0,
		brush_on: true,
	}
}

// Pull the brush Up.
func (t *Turtle) Up() {
	t.brush_on = false
}

// Push the brush Down.
func (t *Turtle) Down() {
	t.brush_on = true
}

const _rad = math.Pi / 180

// Move the turtle Forward.
//
// :param step: Distance to move Forward.
func (t *Turtle) Forward(step float64) {
	x := t.pos_x + math.Cos(t.rotation*_rad)*step
	y := t.pos_y + math.Sin(t.rotation*_rad)*step
	prev_brush_state := t.brush_on
	t.brush_on = true
	t.Move(x, y)
	t.brush_on = prev_brush_state
}

// """Move the turtle to a coordinate.
//
// :param x: x coordinate
// :param y: y coordinate
// """
func (t *Turtle) Move(x, y float64) {
	if t.brush_on {
		// TODO: iterator
		Line(round(t.pos_x), round(t.pos_y), round(x), round(y))(func(lx, ly int) {
			t.Set(float64(lx), float64(ly))
		})
	}

	t.pos_x = x
	t.pos_y = y
}

// """Rotate the turtle (positive direction).
//
// :param angle: Integer. Rotation angle in degrees.
// """
func (t *Turtle) Right(angle float64) {
	t.rotation += angle
}

// """Rotate the turtle (negative direction).
//
// :param angle: Integer. Rotation angle in degrees.
// """
func (t *Turtle) Left(angle float64) {
	t.rotation -= angle
}

// """Move the turtle backwards.
//
// :param step: Integer. Distance to move backwards.
// """
func (t *Turtle) Back(step float64) {
	t.Forward(-step)
}
