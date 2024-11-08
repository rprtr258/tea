package tea

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"
	"syscall"
	"time"

	// "github.com/muesli/ansi/compressor"

	"github.com/muesli/termenv"
	"github.com/rprtr258/fun"
	"github.com/samber/lo"
)

const (
	// defaultFramerate specifies the maximum interval at which we should
	// update the view.
	_fpsDefault = 60
	_fpsMax     = 120
)

// msgRepaint forces a full repaint.
type msgRepaint struct{}

// Renderer is a framerate-based terminal renderer, updating the view
// at a given framerate to avoid overloading the terminal emulator.
//
// In cases where very high performance is needed the renderer can be told
// to exclude ranges of lines, allowing them to be written to directly.
type Renderer struct {
	mu  *sync.Mutex
	out *termenv.Output

	buf                bytes.Buffer
	queuedMessageLines []string
	frameDuration      time.Duration
	ticker             *time.Ticker
	done               chan struct{}
	lastRender         string
	linesRendered      int
	once               sync.Once

	// cursor visibility state
	cursorHidden bool

	// essentially whether or not we're using the full size of the terminal
	altScreenActive bool

	// renderer dimensions; usually the size of the window
	width  int
	height int

	// lines explicitly set not to render
	ignoreLines map[int]struct{}
}

// newRenderer creates a new renderer. Normally you'll want to initialize it
// with os.Stdout as the first argument.
func newRenderer(out *termenv.Output, fps int) *Renderer {
	fps = min(fun.IF(fps >= 1, fps, _fpsDefault), _fpsMax)
	return &Renderer{
		out:                out,
		mu:                 &sync.Mutex{},
		done:               make(chan struct{}),
		frameDuration:      time.Second / time.Duration(fps),
		queuedMessageLines: []string{},
	}
}

// start starts the renderer.
func (r *Renderer) start() {
	if r.ticker == nil {
		r.ticker = time.NewTicker(r.frameDuration)
	} else {
		// If the ticker already exists, it has been stopped and we need to
		// reset it.
		r.ticker.Reset(r.frameDuration)
	}

	// Since the renderer can be restarted after a stop, we need to reset
	// the done channel and its corresponding sync.Once.
	r.once = sync.Once{}

	go r.listen()
}

// stop permanently halts the renderer, rendering the final frame if flush.
func (r *Renderer) stop(flush bool) {
	// Stop the renderer before acquiring the mutex to avoid a deadlock.
	r.once.Do(func() {
		r.done <- struct{}{}
	})

	if flush {
		// flush locks the mutex
		r.flush()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.out.ClearLine()
}

// listen waits for ticks on the ticker, or a signal to stop the renderer.
func (r *Renderer) listen() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			return

		case <-r.ticker.C:
			r.flush()
		}
	}
}

// flush renders the buffer.
func (r *Renderer) flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.buf.Len() == 0 || r.buf.String() == r.lastRender {
		// Nothing to do
		return
	}

	// r.out.MoveCursor(0, 0)
	// r.out.Write(r.buf.Bytes())
	// os.Stdout.Write(r.buf.Bytes())
	syscall.Write(1, []byte("\x1b[0;0H"))
	for _, chunk := range lo.Chunk(r.buf.Bytes(), 16*1024) {
		syscall.Write(1, chunk)
	}
}

func (r *Renderer) reset() {
	r.buf.Reset()
}

// write writes to the internal buffer. The buffer will be outputted via the
// ticker which calls flush().
func (r *Renderer) Write(s []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// _, _ = r.buf.Write(s)
	r.buf = *bytes.NewBuffer(s)
}

func (r *Renderer) repaint() {
	r.lastRender = ""
}

func (r *Renderer) clearScreen() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.out.ClearScreen()
	r.out.MoveCursor(1, 1)

	r.repaint()
}

func (r *Renderer) altScreen() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.altScreenActive
}

func (r *Renderer) enterAltScreen() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.altScreenActive {
		return
	}

	r.altScreenActive = true
	r.out.AltScreen()

	// Ensure that the terminal is cleared, even when it doesn't support
	// alt screen (or alt screen support is disabled, like GNU screen by
	// default).
	//
	// Note: we can't use r.clearScreen() here because the mutex is already
	// locked.
	r.out.ClearScreen()
	r.out.MoveCursor(1, 1)

	// cmd.exe and other terminals keep separate cursor states for the AltScreen
	// and the main buffer. We have to explicitly reset the cursor visibility
	// whenever we enter AltScreen.
	if r.cursorHidden {
		r.out.HideCursor()
	} else {
		r.out.ShowCursor()
	}

	r.repaint()
}

func (r *Renderer) exitAltScreen() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.altScreenActive {
		return
	}

	r.altScreenActive = false
	r.out.ExitAltScreen()

	// cmd.exe and other terminals keep separate cursor states for the AltScreen
	// and the main buffer. We have to explicitly reset the cursor visibility
	// whenever we exit AltScreen.
	if r.cursorHidden {
		r.out.HideCursor()
	} else {
		r.out.ShowCursor()
	}

	r.repaint()
}

func (r *Renderer) setCursor(show bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cursorHidden = !show
	if r.cursorHidden {
		r.out.HideCursor()
	} else {
		r.out.ShowCursor()
	}
}

func (r *Renderer) setMouseCellMotion(enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if enabled {
		r.out.EnableMouseCellMotion()
	} else {
		r.out.DisableMouseCellMotion()
	}
}

func (r *Renderer) setMouseAllMotion(enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if enabled {
		r.out.EnableMouseAllMotion()
	} else {
		r.out.DisableMouseAllMotion()
	}
}

// setIgnoredLines specifies lines not to be touched by the standard Tea renderer.
func (r *Renderer) setIgnoredLines(from, to int) {
	// Lock if we're going to be clearing some lines since we don't want
	// anything jacking our cursor.
	if r.linesRendered > 0 {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	if r.ignoreLines == nil {
		r.ignoreLines = make(map[int]struct{})
	}
	for i := from; i < to; i++ {
		r.ignoreLines[i] = struct{}{}
	}

	// Erase ignored lines
	if r.linesRendered > 0 {
		buf := &bytes.Buffer{}
		out := termenv.NewOutput(buf)

		for i := r.linesRendered - 1; i >= 0; i-- {
			if _, exists := r.ignoreLines[i]; exists {
				out.ClearLine()
			}
			out.CursorUp(1)
		}
		out.MoveCursor(r.linesRendered, 0) // put cursor back
		_, _ = r.out.Write(buf.Bytes())
	}
}

// clearIgnoredLines returns control of any ignored lines to the standard Tea renderer.
// That is, any lines previously set to be ignored can be rendered to again.
func (r *Renderer) clearIgnoredLines() {
	r.ignoreLines = nil
}

// insertTop effectively scrolls up. It inserts lines at the top of a given
// area designated to be a scrollable region, pushing everything else down.
// This is roughly how ncurses does it.
//
// To call this function use command ScrollUp().
//
// For this to work renderer.ignoreLines must be set to ignore the scrollable
// region since we are bypassing the normal Tea renderer here.
//
// Because this method relies on the terminal dimensions, it's only valid for
// full-window applications (generally those that use the alternate screen
// buffer).
//
// This method bypasses the normal rendering buffer and is philosophically
// different than the normal way we approach rendering in Tea. It's for
// use in high-performance rendering, such as a pager that could potentially
// be rendering very complicated ansi. In cases where the content is simpler
// standard Tea rendering should suffice.
func (r *Renderer) insertTop(lines []string, topBoundary, bottomBoundary int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	buf := &bytes.Buffer{}
	out := termenv.NewOutput(buf)

	out.ChangeScrollingRegion(topBoundary, bottomBoundary)
	out.MoveCursor(topBoundary, 0)
	out.InsertLines(len(lines))
	_, _ = out.WriteString(strings.Join(lines, "\r\n"))
	out.ChangeScrollingRegion(0, r.height)

	// Move cursor back to where the main rendering routine expects it to be
	out.MoveCursor(r.linesRendered, 0)

	_, _ = r.out.Write(buf.Bytes())
}

// insertBottom effectively scrolls down. It inserts lines at the bottom of
// a given area designated to be a scrollable region, pushing everything else
// up. This is roughly how ncurses does it.
//
// To call this function use the command ScrollDown().
//
// See note in insertTop() for caveats, how this function only makes sense for
// full-window applications, and how it differs from the normal way we do
// rendering in Tea.
func (r *Renderer) insertBottom(lines []string, topBoundary, bottomBoundary int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	buf := &bytes.Buffer{}
	out := termenv.NewOutput(buf)

	out.ChangeScrollingRegion(topBoundary, bottomBoundary)
	out.MoveCursor(bottomBoundary, 0)
	_, _ = out.WriteString("\r\n" + strings.Join(lines, "\r\n"))
	out.ChangeScrollingRegion(0, r.height)

	// Move cursor back to where the main rendering routine expects it to be
	out.MoveCursor(r.linesRendered, 0)

	_, _ = r.out.Write(buf.Bytes())
}

// handleMessages handles internal messages for the renderer.
func (r *Renderer) handleMessages(msg Msg) {
	switch msg := msg.(type) {
	case msgRepaint:
		// Force a repaint by clearing the render cache as we slide into a
		// render.
		r.mu.Lock()
		r.repaint()
		r.mu.Unlock()
	case MsgWindowSize:
		r.mu.Lock()
		r.width = msg.Width
		r.height = msg.Height
		r.repaint()
		r.mu.Unlock()
	case msgClearScrollArea:
		r.clearIgnoredLines()

		// Force a repaint on the area where the scrollable stuff was in this
		// update cycle
		r.mu.Lock()
		r.repaint()
		r.mu.Unlock()
	case msgSyncScrollArea:
		// Re-render scrolling area
		r.clearIgnoredLines()
		r.setIgnoredLines(msg.topBoundary, msg.bottomBoundary)
		r.insertTop(msg.lines, msg.topBoundary, msg.bottomBoundary)

		// Force non-scrolling stuff to repaint in this update cycle
		r.mu.Lock()
		r.repaint()
		r.mu.Unlock()
	case msgScrollUp:
		r.insertTop(msg.lines, msg.topBoundary, msg.bottomBoundary)
	case msgScrollDown:
		r.insertBottom(msg.lines, msg.topBoundary, msg.bottomBoundary)
	case msgPrintLine:
		if !r.altScreenActive {
			lines := strings.Split(msg.messageBody, "\n")
			r.mu.Lock()
			r.queuedMessageLines = append(r.queuedMessageLines, lines...)
			r.repaint()
			r.mu.Unlock()
		}
	}
}

// HIGH-PERFORMANCE RENDERING STUFF

type msgSyncScrollArea struct {
	lines          []string
	topBoundary    int
	bottomBoundary int
}

// SyncScrollArea performs a paint of the entire region designated to be the
// scrollable area. This is required to initialize the scrollable region and
// should also be called on resize (MsgWindowSize).
//
// For high-performance, scroll-based rendering only.
func SyncScrollArea(lines []string, topBoundary, bottomBoundary int) Cmd {
	return func() Msg {
		return msgSyncScrollArea{
			lines:          lines,
			topBoundary:    topBoundary,
			bottomBoundary: bottomBoundary,
		}
	}
}

type msgClearScrollArea struct{}

// ClearScrollArea deallocates the scrollable region and returns the control of
// those lines to the main rendering routine.
//
// For high-performance, scroll-based rendering only.
func ClearScrollArea() Msg {
	return msgClearScrollArea{}
}

type msgScrollUp struct {
	lines          []string
	topBoundary    int
	bottomBoundary int
}

// ScrollUp adds lines to the top of the scrollable region, pushing existing
// lines below down. Lines that are pushed out the scrollable region disappear
// from view.
//
// For high-performance, scroll-based rendering only.
func ScrollUp(newLines []string, topBoundary, bottomBoundary int) Cmd {
	return func() Msg {
		return msgScrollUp{
			lines:          newLines,
			topBoundary:    topBoundary,
			bottomBoundary: bottomBoundary,
		}
	}
}

type msgScrollDown struct {
	lines          []string
	topBoundary    int
	bottomBoundary int
}

// ScrollDown adds lines to the bottom of the scrollable region, pushing
// existing lines above up. Lines that are pushed out of the scrollable region
// disappear from view.
//
// For high-performance, scroll-based rendering only.
func ScrollDown(newLines []string, topBoundary, bottomBoundary int) Cmd {
	return func() Msg {
		return msgScrollDown{
			lines:          newLines,
			topBoundary:    topBoundary,
			bottomBoundary: bottomBoundary,
		}
	}
}

type msgPrintLine struct {
	messageBody string
}

// Println prints above the Program. This output is unmanaged by the program and
// will persist across renders by the Program.
//
// Unlike fmt.Println (but similar to log.Println) the message will be print on
// its own line.
//
// If the altscreen is active no output will be printed.
func Println(args ...any) Cmd {
	return func() Msg {
		log.Println(args...)
		return msgPrintLine{
			messageBody: fmt.Sprint(args...),
		}
	}
}

// Printf prints above the Program. It takes a format template followed by
// values similar to fmt.Printf. This output is unmanaged by the program and
// will persist across renders by the Program.
//
// Unlike fmt.Printf (but similar to log.Printf) the message will be print on
// its own line.
//
// If the altscreen is active no output will be printed.
func Printf(format string, args ...any) Cmd {
	return func() Msg {
		log.Printf(format, args...)
		return msgPrintLine{
			messageBody: fmt.Sprintf(format, args...),
		}
	}
}
