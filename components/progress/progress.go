package progress

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/ansi"
	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/styles"
)

const (
	fps              = 60
	defaultWidth     = 40
	defaultFrequency = 18.0
	defaultDamping   = 1.0
)

// Option is used to set options in New. For example:
//
//	    progress := New(
//		       WithRamp("#ff0000", "#0000ff"), // TODO: update color example
//		       WithoutPercentage(),
//	    )
type Option func(*Model)

var (
	_defaultColorA, _ = colorful.Hex("#5A56E0")
	_defaultColorB, _ = colorful.Hex("#EE6FF8")
)

// WithDefaultGradient sets a gradient fill with default colors.
func WithDefaultGradient() Option {
	return WithGradient(_defaultColorA, _defaultColorB)
}

// WithGradient sets a gradient fill blending between two colors.
func WithGradient(colorA, colorB colorful.Color) Option {
	return func(m *Model) {
		m.setRamp(colorA, colorB, false)
	}
}

// WithDefaultScaledGradient sets a gradient with default colors, and scales the
// gradient to fit the filled portion of the ramp.
func WithDefaultScaledGradient() Option {
	return WithScaledGradient(_defaultColorA, _defaultColorB)
}

// WithScaledGradient scales the gradient to fit the width of the filled portion of
// the progress bar.
func WithScaledGradient(colorA, colorB colorful.Color) Option {
	return func(m *Model) {
		m.setRamp(colorA, colorB, true)
	}
}

// WithSolidFill sets the progress to use a solid fill with the given color.
func WithSolidFill(color string) Option {
	return func(m *Model) {
		m.FullColor = color
		m.useRamp = false
	}
}

// WithoutPercentage hides the numeric percentage.
func WithoutPercentage() Option {
	return func(m *Model) {
		m.ShowPercentage = false
	}
}

// WithWidth sets the initial width of the progress bar.
// Note that you can also set the width via the Width property,
// which can come in handy if you're waiting for a tea.MsgWindowSize.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.Width = w
	}
}

// WithSpringOptions sets the initial frequency and damping options for the
// progress bar's built-in spring-based animation.
// Frequency corresponds to speed, and damping to bounciness.
func WithSpringOptions(frequency, damping float64) Option {
	return func(m *Model) {
		m.SetSpringOptions(frequency, damping)
		m.springCustomized = true
	}
}

// MsgFrame indicates that an animation step should occur.
type MsgFrame struct {
	tag int
}

// Model stores values we'll use when rendering the progress bar.
type Model struct {
	// An identifier to keep us from receiving frame messages too quickly.
	tag int

	// Total width of the progress bar, including percentage, if set.
	Width int

	// "Filled" sections of the progress bar.
	Full      rune
	FullColor string

	// "Empty" sections of the progress bar.
	Empty      rune
	EmptyColor string

	// Settings for rendering the numeric percentage.
	ShowPercentage  bool
	PercentFormat   string // a fmt string for a float
	PercentageStyle styles.Style

	// Members for animated transitions.
	spring           harmonica.Spring
	springCustomized bool
	percentShown     float64 // percent currently displaying
	targetPercent    float64 // percent to which we're animating
	velocity         float64

	// Gradient settings
	useRamp    bool
	rampColorA colorful.Color
	rampColorB colorful.Color

	// When true, we scale the gradient to fit the width of the filled section
	// of the progress bar. When false, the width of the gradient will be set
	// to the full width of the progress bar.
	scaleRamp bool
}

// New returns a model with default values.
func New(opts ...Option) Model {
	m := Model{
		Width:          defaultWidth,
		Full:           '█',
		FullColor:      "#7571F9",
		Empty:          '░',
		EmptyColor:     "#606060",
		ShowPercentage: true,
		PercentFormat:  " %3.0f%%",
	}
	if !m.springCustomized {
		m.SetSpringOptions(defaultFrequency, defaultDamping)
	}

	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// Update is used to animate the progress bar during transitions. Use
// SetPercent to create the command you'll need to trigger the animation.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) Update(c tea.Context[*Model], msg tea.Msg) {
	switch msg := msg.(type) {
	case MsgFrame:
		if msg.tag != m.tag {
			return
		}

		// If we've more or less reached equilibrium, stop updating.
		if !m.IsAnimating() {
			return
		}

		m.percentShown, m.velocity = m.spring.Update(m.percentShown, m.velocity, m.targetPercent)
		m.nextFrame(c)
	}
}

// SetSpringOptions sets the frequency and damping for the current spring.
// Frequency corresponds to speed, and damping to bounciness.
func (m *Model) SetSpringOptions(frequency, damping float64) {
	m.spring = harmonica.NewSpring(harmonica.FPS(fps), frequency, damping)
}

// Percent returns the current visible percentage on the model. This is only
// relevant when you're animating the progress bar.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) Percent() float64 {
	return m.targetPercent
}

// SetPercent sets the percentage state of the model as well as a command
// necessary for animating the progress bar to this new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) SetPercent(c tea.Context[*Model], p float64) {
	m.targetPercent = max(0, min(1, p))
	m.tag++
	m.nextFrame(c)
}

// IncrPercent increments the percentage by a given amount, returning a command
// necessary to animate the progress bar to the new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) IncrPercent(c tea.Context[*Model], v float64) {
	m.SetPercent(c, m.Percent()+v)
}

// DecrPercent decrements the percentage by a given amount, returning a command
// necessary to animate the progress bar to the new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) DecrPercent(c tea.Context[*Model], v float64) {
	m.SetPercent(c, m.Percent()-v)
}

// View renders an animated progress bar in its current state. To render
// a static progress bar based on your own calculations use ViewAs instead.
func (m *Model) View(vb tea.Viewbox) {
	m.ViewAs(vb, m.percentShown)
}

// ViewAs renders the progress bar with a given percentage.
func (m *Model) ViewAs(vb tea.Viewbox, percent float64) {
	percentView := m.percentageView(percent)
	textWidth := ansi.PrintableRuneWidth(percentView)
	m.barView(vb, percent, textWidth)
	vb.PaddingLeft(m.Width - textWidth).WriteLine(percentView)
}

func (m *Model) nextFrame(c tea.Context[*Model]) {
	d := time.Second / time.Duration(fps)
	// TODO: tea.Tick(time.Second/time.Duration(fps))
	msg := MsgFrame{
		tag: m.tag,
	}
	c.F(func() tea.Msg2[*Model] {
		return func(m *Model) {
			<-time.After(d)
			m.Update(c, msg)
		}
	})
}

func (m *Model) barView(vb tea.Viewbox, percent float64, textWidth int) {
	tw := max(0, m.Width-textWidth)              // total width
	fw := int(math.Round(float64(tw) * percent)) // filled width

	fw = max(0, min(tw, fw))

	if m.useRamp {
		// Gradient fill
		for i := 0; i < fw; i++ {
			p := fun.Switch(true, float64(i)/float64(tw-1)).
				// this is up for debate: in a gradient of width=1, should the
				// single character rendered be the first color, the last color
				// or exactly 50% in between? I opted for 50%
				Case(0.5, fw == 1).
				Case(float64(i)/float64(fw-1), m.scaleRamp).
				End()

			vb = vb.
				PaddingLeft(i).
				Styled(styles.Style{}.Foreground(scuf.FgRGB(m.rampColorA.BlendLuv(m.rampColorB, p).RGB255()))).
				WriteLineX(string(m.Full))
		}
	} else {
		// Solid fill
		vb.
			Styled(styles.Style{}.Foreground(scuf.FgRGB(scuf.MustParseHexRGB(m.FullColor)))).
			WriteLine(strings.Repeat(string(m.Full), fw))
	}

	// Empty fill
	n := max(0, tw-fw)
	vb.
		PaddingLeft(fw).
		Styled(styles.Style{}.Foreground(scuf.FgRGB(scuf.MustParseHexRGB(m.EmptyColor)))).
		WriteLine(strings.Repeat(string(m.Empty), n))
}

func (m *Model) percentageView(percent float64) string {
	if !m.ShowPercentage {
		return ""
	}

	percent = max(0, min(1, percent))
	percentage := fmt.Sprintf(m.PercentFormat, percent*100) //nolint:gomnd
	return m.PercentageStyle.Render(percentage)
}

func (m *Model) setRamp(colorA, colorB colorful.Color, scaled bool) {
	// // In the event of an error colors here will default to black. For
	// // usability's sake, and because such an error is only cosmetic, we're
	// // ignoring the error.
	// a, _ :=
	// b, _ := colorful.Hex(colorB)
	m.useRamp = true
	m.scaleRamp = scaled
	m.rampColorA = colorA
	m.rampColorB = colorB
}

// IsAnimating returns false if the progress bar reached equilibrium and is no longer animating.
func (m *Model) IsAnimating() bool {
	dist := math.Abs(m.percentShown - m.targetPercent)
	return dist >= 0.001 || m.velocity >= 0.01
}
