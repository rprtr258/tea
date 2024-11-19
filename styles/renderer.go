package styles

import (
	"io"

	"github.com/muesli/termenv"
)

// Renderer is a styles terminal renderer.
type Renderer struct {
	Output            *termenv.Output
	hasDarkBackground *bool
}

// We're manually creating the struct here to avoid initializing the output and
// query the terminal multiple times.
var _renderer = &Renderer{
	Output: termenv.DefaultOutput(),
}

// NewRenderer creates a new Renderer.
// w will be used to determine the terminal's color capabilities.
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{
		Output: termenv.NewOutput(w),
	}
}

// HasDarkBackground returns whether or not the renderer will render to a dark
// background. A dark background can either be auto-detected, or set explicitly
// on the renderer.
func (r *Renderer) HasDarkBackground() bool {
	if r.hasDarkBackground != nil {
		return *r.hasDarkBackground
	}
	return r.Output.HasDarkBackground()
}

// SetHasDarkBackground sets the background color detection value on the
// renderer. This function exists mostly for testing purposes so that you can
// assure you're testing against a specific background color setting.
//
// Outside of testing you likely won't want to use this function as the
// backgrounds value will be automatically detected and cached against the
// terminal's current background color setting.
//
// This function is thread-safe.
func (r *Renderer) SetHasDarkBackground(b bool) {
	r.hasDarkBackground = &b
}
