package spinner

import (
	"sync"
	"time"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/lipgloss"
)

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var (
	lastID int
	idMtx  sync.Mutex
)

// Return the next ID we should use on the Model.
func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

// Spinner is a set of frames used in animating the spinner.
type Spinner struct {
	Frames []string
	FPS    time.Duration
}

// Some spinners to choose from. You could also make your own.
var (
	Line = Spinner{
		Frames: []string{"|", "/", "-", "\\"},
		FPS:    time.Second / 10, //nolint:gomnd
	}
	Dot = Spinner{
		Frames: []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		FPS:    time.Second / 10, //nolint:gomnd
	}
	MiniDot = Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    time.Second / 12, //nolint:gomnd
	}
	Jump = Spinner{
		Frames: []string{"⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"},
		FPS:    time.Second / 10, //nolint:gomnd
	}
	Pulse = Spinner{
		Frames: []string{"█", "▓", "▒", "░"},
		FPS:    time.Second / 8, //nolint:gomnd
	}
	Points = Spinner{
		Frames: []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●"},
		FPS:    time.Second / 7, //nolint:gomnd
	}
	Globe = Spinner{
		Frames: []string{"🌍", "🌎", "🌏"},
		FPS:    time.Second / 4, //nolint:gomnd
	}
	Moon = Spinner{
		Frames: []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
		FPS:    time.Second / 8, //nolint:gomnd
	}
	Monkey = Spinner{
		Frames: []string{"🙈", "🙉", "🙊"},
		FPS:    time.Second / 3, //nolint:gomnd
	}
	Meter = Spinner{
		Frames: []string{
			"▱▱▱",
			"▰▱▱",
			"▰▰▱",
			"▰▰▰",
			"▰▰▱",
			"▰▱▱",
			"▱▱▱",
		},
		FPS: time.Second / 7, //nolint:gomnd
	}
	Hamburger = Spinner{
		Frames: []string{"☱", "☲", "☴", "☲"},
		FPS:    time.Second / 3, //nolint:gomnd
	}
	Ellipsis = Spinner{
		Frames: []string{"", ".", "..", "..."},
		FPS:    time.Second / 3, //nolint:gomnd
	}
	Circle = Spinner{
		Frames: []string{
			"⠈⠁", "⠈⠑", "⠈⠱", "⠈⡱", "⢀⡱", "⢄⡱", "⢄⡱", "⢆⡱",
			"⢎⡱", "⢎⡰", "⢎⡠", "⢎⡀", "⢎⠁", "⠎⠁", "⠊⠁",
		},
		FPS: 100 * time.Millisecond,
	}
)

// Model contains the state for the spinner. Use New to create new models
// rather than using Model as a struct literal.
type Model struct {
	// Spinner settings to use. See type Spinner.
	Spinner Spinner

	// Style sets the styling for the spinner. Most of the time you'll just
	// want foreground and background coloring, and potentially some padding.
	//
	// For an introduction to styling with Lip Gloss see:
	// https://github.com/rprtr258/tea/lipgloss
	Style lipgloss.Style

	frame int
	id    int
	tag   int
}

// ID returns the spinner's unique ID.
func (m *Model) ID() int {
	return m.id
}

// New returns a model with default values.
func New(opts ...Option) Model {
	m := Model{
		Spinner: Line,
		id:      nextID(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// MsgTick indicates that the timer has ticked and we should render a frame.
type MsgTick struct {
	Time time.Time
	tag  int
	ID   int
}

// Update is the Tea update function.
func (m *Model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case MsgTick:
		// If an ID is set, and the ID doesn't belong to this spinner, reject
		// the message.
		if msg.ID > 0 && msg.ID != m.id {
			return nil
		}

		// If a tag is set, and it's not the one we expect, reject the message.
		// This prevents the spinner from receiving too many messages and
		// thus spinning too fast.
		if msg.tag > 0 && msg.tag != m.tag {
			return nil
		}

		m.frame++
		if m.frame >= len(m.Spinner.Frames) {
			m.frame = 0
		}

		m.tag++
		return []tea.Cmd{m.tick(m.id, m.tag)}
	default:
		return nil
	}
}

// View renders the model's view.
func (m *Model) View() string {
	if m.frame >= len(m.Spinner.Frames) {
		return "(error)"
	}

	return m.Style.Render(m.Spinner.Frames[m.frame])
}

// CmdTick is the command used to advance the spinner one frame. Use this command
// to effectively start the spinner.
func (m *Model) CmdTick() tea.Msg {
	return MsgTick{
		// The time at which the tick occurred.
		Time: time.Now(),

		// The ID of the spinner that this message belongs to. This can be
		// helpful when routing messages, however bear in mind that spinners
		// will ignore messages that don't contain ID by default.
		ID: m.id,

		tag: m.tag,
	}
}

func (m *Model) tick(id, tag int) tea.Cmd {
	return tea.Tick(m.Spinner.FPS, func(t time.Time) tea.Msg {
		return MsgTick{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

// Option is used to set options in New. For example:
//
//	spinner := New(WithSpinner(Dot))
type Option func(*Model)

// WithSpinner is an option to set the spinner.
func WithSpinner(spinner Spinner) Option {
	return func(m *Model) {
		m.Spinner = spinner
	}
}

// WithStyle is an option to set the spinner style.
func WithStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.Style = style
	}
}
