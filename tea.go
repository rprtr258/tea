// Package tea provides a framework for building rich terminal user interfaces
// based on the paradigms of The Elm Architecture. It's well-suited for simple
// and complex terminal applications, either inline, full-window, or a mix of
// both. It's been battle-tested in several large projects and is
// production-ready.
//
// A tutorial is available at https://github.com/rprtr258/tea/tree/master/tutorials
//
// Example programs can be found at https://github.com/rprtr258/tea/tree/master/examples
package tea

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/containerd/console"
	isatty "github.com/mattn/go-isatty"
	"github.com/muesli/cancelreader"
	"github.com/muesli/termenv"
)

// ErrProgramKilled is returned by [Program.Run] when the program got killed.
var ErrProgramKilled = errors.New("program was killed")

// Msg contain data from the result of a IO operation.
// Msgs trigger the update function and, henceforth, the UI.
type Msg any

// Model contains the program's state as well as its core functions.
type Model interface {
	// Init is the first function that will be called. It returns list of initial commands.
	Init(func(...Cmd))

	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the model and/or send a command.
	Update(Msg, func(...Cmd))

	// View renders the program's UI. The view is rendered after every Update.
	View(Viewbox)
}

// Cmd is an IO operation that returns a message when it's complete. If it's
// nil it's considered a no-op. Use it for things like HTTP requests, timers,
// saving and loading from disk, and so on.
//
// Note that there's almost never a reason to use a command to send a message
// to another part of your program. That can almost always be done in the
// update function.
type Cmd func() Msg

type Msg2[M any] func(M)

type Context[M any] struct {
	Dispatch func(...Cmd)

	F func(...func() Msg2[M])
}

func Of[M, N any](c Context[M], ff func(M) N) Context[N] {
	return Context[N]{
		Dispatch: c.Dispatch,
		F: func(cmds ...func() Msg2[N]) {
			for _, cmd := range cmds {
				c.F(func() Msg2[M] {
					return func(m M) {
						cmd()(ff(m))
					}
				})
			}
		},
	}
}

type inputType int

const (
	defaultInput inputType = iota
	ttyInput
	customInput
)

// String implements the stringer interface for [inputType].
// It is inteded to be used in testing.
func (i inputType) String() string {
	return [...]string{
		"default input",
		"tty input",
		"custom input",
	}[i]
}

// Options to customize the program during its initialization. These are
// generally set with ProgramOptions.
//
// The options here are treated as bits.
type startupOptions byte

const (
	withAltScreen startupOptions = 1 << iota
	withMouseCellMotion
	withMouseAllMotion
	withoutSignalHandler

	// Catching panics is incredibly useful for restoring the terminal to a
	// usable state after a panic occurs. When this is set, Tea will
	// recover from panics, print the stack trace, and disable raw mode. This
	// feature is on by default.
	withoutCatchPanics
)

func (s startupOptions) has(option startupOptions) bool {
	return s&option != 0
}

// handlers manages series of channels returned by various processes. It allows
// us to wait for those processes to terminate before exiting the program.
type handlers []<-chan struct{}

// handlersShutdown waits for all handlers to terminate.
func handlersShutdown(h handlers) {
	var wg sync.WaitGroup
	for _, ch := range h {
		wg.Add(1)
		go func(ch <-chan struct{}) {
			<-ch
			wg.Done()
		}(ch)
	}
	wg.Wait()
}

// Program is a terminal user interface.
type Program[M Model] struct {
	model M

	// Configuration options that will set as the program is initializing,
	// treated as bits. These options can be set via various ProgramOptions.
	startupOptions startupOptions

	inputType inputType

	ctx    context.Context //nolint:containedctx // TODO: remove
	cancel context.CancelFunc

	msgs chan Msg
	errs chan error
	done chan struct{}

	// where to send output, this will usually be os.Stdout.
	output        *termenv.Output
	restoreOutput func() error
	renderer      *Renderer
	vb            Viewbox

	// where to read inputs from, this will usually be os.Stdin.
	input        io.Reader
	cancelReader cancelreader.CancelReader
	readLoopDone chan struct{}
	console      console.Console

	// was the altscreen active before releasing the terminal?
	altScreenWasActive bool
	ignoreSignals      bool

	filter func(M, Msg) Msg

	// fps is the frames per second we should set on the renderer, if applicable,
	fps int
}

// MsgQuit signals that the program should quit. You can send a MsgQuit with
// Quit.
type MsgQuit struct{}

// Quit is a special command that tells the Tea program to exit.
func Quit() Msg {
	return MsgQuit{}
}

type model2[M any] interface {
	Init(Context[M])
	Update(Context[M], Msg)
	View(Viewbox)
}

type AdapterModel[M model2[M]] struct {
	M M
}

func (m *AdapterModel[M]) Init(f func(...Cmd)) {
	c := Context[M]{
		Dispatch: f,
		F: func(fs ...func() Msg2[M]) {
			for _, fn := range fs {
				f(func() Msg {
					fn()(m.M)
					return nil // TODO: ???
				})
			}
		},
	}
	m.M.Init(c)
}
func (m *AdapterModel[M]) Update(msg Msg, f func(...Cmd)) {
	c := Context[M]{
		Dispatch: f,
		F: func(fs ...func() Msg2[M]) {
			for _, fn := range fs {
				f(func() Msg {
					fn()(m.M)
					return nil // TODO: ???
				})
			}
		},
	}
	m.M.Update(c, msg)
}
func (m *AdapterModel[M]) View(vb Viewbox) {
	m.M.View(vb)
}

func NewProgram2[M model2[M]](ctx context.Context, model M) *Program[*AdapterModel[M]] {
	return NewProgram(ctx, &AdapterModel[M]{M: model})
}

// NewProgram creates a new Program.
func NewProgram[M Model](ctx context.Context, model M) *Program[M] {
	// Initialize context and teardown channel.
	ctx, cancel := context.WithCancel(ctx)

	output := termenv.DefaultOutput()
	// cache detected color values
	termenv.WithColorCache(true)(output)

	restoreOutput, _ := termenv.EnableVirtualTerminalProcessing(output)

	return &Program[M]{
		model:         model,
		msgs:          make(chan Msg),
		output:        output,
		ctx:           ctx,
		cancel:        cancel,
		restoreOutput: restoreOutput,
		renderer:      newRenderer(output /*p.fps*/, _fpsDefault),
	}
}

func (p *Program[M]) handleSignals() chan struct{} {
	ch := make(chan struct{})

	// Listen for SIGINT and SIGTERM.
	//
	// In most cases ^C will not send an interrupt because the terminal will be
	// in raw mode and ^C will be captured as a keystroke and sent along to
	// Program.Update as a MsgKey. When input is not a TTY, however, ^C will be
	// caught here.
	//
	// SIGTERM is sent by unix utilities (like kill) to terminate a process.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer func() {
			signal.Stop(sig)
			close(ch)
		}()

		for {
			select {
			case <-p.ctx.Done():
				return

			case <-sig:
				if !p.ignoreSignals {
					p.msgs <- MsgQuit{}
					return
				}
			}
		}
	}()

	return ch
}

// listenForResize sends messages (or errors) when the terminal resizes.
// Argument output should be the file descriptor for the terminal; usually
// os.Stdout.
func (p *Program[M]) listenForResize(done chan struct{}) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	defer func() {
		signal.Stop(sig)
		close(done)
	}()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-sig:
		}

		p.checkResize()
	}
}

// handleResize handles terminal resize events.
func (p *Program[M]) handleResize() chan struct{} {
	ch := make(chan struct{})

	if f, ok := p.output.TTY().(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		// Get the initial terminal size and send it to the program.
		go p.checkResize()

		// Listen for window resizes.
		go p.listenForResize(ch)
	} else {
		close(ch)
	}

	return ch
}

// handleCommands runs commands in a goroutine and sends the result to the
// program's message channel.
func (p *Program[M]) handleCommands(cmds chan []Cmd) chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			case <-p.ctx.Done():
				return

			case cmds := <-cmds:
				// Don't wait on these goroutines, otherwise the shutdown
				// latency would get too large as a Cmd can run for some time
				// (e.g. tick commands that sleep for half a second). It's not
				// possible to cancel them so we'll have to leak the goroutine
				// until Cmd returns.
				for _, cmd := range cmds {
					go func() {
						msg := cmd()
						if msg == nil {
							msg = msgRepaint{}
						}
						p.Send(msg) // this can be long.
					}()
				}
			}
		}
	}()

	return ch
}

// eventLoop is the central message loop. It receives and handles the default
// Tea messages, update the model and triggers redraws.
func (p *Program[M]) eventLoop(model M, cmds chan []Cmd) (M, error) {
	for {
		select {
		case <-p.ctx.Done():
			return model, nil

		case err := <-p.errs:
			return model, err

		case msg := <-p.msgs:
			// Filter messages.
			if p.filter != nil {
				msg = p.filter(model, msg)
			}
			if msg == nil {
				continue
			}

			// Handle special internal messages.
			switch msg := msg.(type) {
			case MsgQuit:
				return model, nil

			case msgClearScreen:
				p.renderer.clearScreen()

			case msgEnterAltScreen:
				p.renderer.enterAltScreen()

			case msgExitAltScreen:
				p.renderer.exitAltScreen()

			case msgEnableMouseCellMotion:
				p.renderer.setMouseCellMotion(true)

			case msgEnableMouseAllMotion:
				p.renderer.setMouseAllMotion(true)

			case msgDisableMouse:
				p.renderer.setMouseCellMotion(false)
				p.renderer.setMouseAllMotion(false)

			case msgShowCursor:
				p.renderer.setCursor(true)

			case msgHideCursor:
				p.renderer.setCursor(false)

			case msgExec:
				// NB: this blocks.
				p.exec(msg.cmd, msg.fn)

				// TODO: move to renderer
			case MsgWindowSize:
				p.vb = NewViewbox(msg.Height, msg.Width)
			}

			p.renderer.handleMessages(msg)

			model.Update(msg, func(c ...Cmd) {
				cmds <- c
			}) // run update, process command (if any)
			p.renderer.reset()
			p.vb.clear()
			model.View(p.vb)
			p.renderer.Write(p.vb.Render())
		}
	}
}

// Run initializes the program and runs its event loops, blocking until it gets
// terminated by either [Program.Quit], [Program.Kill], or its signal handler.
// Returns the final model.
func (p *Program[M]) Run() (M, error) {
	myHandlers := handlers{}
	cmds := make(chan []Cmd)
	p.errs = make(chan error)
	p.done = make(chan struct{}, 1)

	defer p.cancel()

	switch p.inputType {
	case defaultInput:
		p.input = os.Stdin

		// The user has not set a custom input, so we need to check whether or
		// not standard input is a terminal. If it's not, we open a new TTY for
		// input. This will allow things to "just work" in cases where data was
		// piped in or redirected to the application.
		//
		// To disable input entirely pass nil to the [WithInput] program option.
		f, isFile := p.input.(*os.File)
		if !isFile {
			break
		}
		if isatty.IsTerminal(f.Fd()) {
			break
		}

		f, err := openInputTTY()
		if err != nil {
			return p.model, err
		}
		defer f.Close() //nolint:errcheck // uuh

		p.input = f
	case ttyInput:
		// Open a new TTY, by request
		f, err := openInputTTY()
		if err != nil {
			return p.model, err
		}
		defer f.Close() //nolint:errcheck // uuh

		p.input = f
	case customInput:
		// (There is nothing extra to do.)
	}

	// Handle signals.
	if !p.startupOptions.has(withoutSignalHandler) {
		myHandlers = append(myHandlers, p.handleSignals())
	}

	// Recover from panics.
	if !p.startupOptions.has(withoutCatchPanics) {
		defer func() {
			if r := recover(); r != nil {
				p.shutdown(true)
				fmt.Printf("Caught panic:\n\n%s\n\nRestoring terminal...\n\n", r)
				debug.PrintStack()
			}
		}()
	}

	// Check if output is a TTY before entering raw mode, hiding the cursor and so on.
	if err := p.initTerminal(); err != nil {
		return p.model, err
	}

	// Honor program startup options.
	if p.startupOptions&withAltScreen != 0 {
		p.renderer.enterAltScreen()
	}
	if p.startupOptions&withMouseCellMotion != 0 {
		p.renderer.setMouseCellMotion(true)
	} else if p.startupOptions&withMouseAllMotion != 0 {
		p.renderer.setMouseAllMotion(true)
	}

	// Initialize the program.
	p.model.Init(func(cmdss ...Cmd) { // TODO: remove
		// initCmds = append(initCmds, cmdss...)
		go func() {
			cmds <- cmdss
		}()
	})

	// Start the renderer.
	p.renderer.start()

	// Render the initial view.
	p.renderer.reset()
	p.vb.clear()
	p.model.View(p.vb)
	p.renderer.Write(p.vb.Render())

	// Subscribe to user input.
	if p.input != nil {
		if err := p.initCancelReader(); err != nil {
			return p.model, err
		}
	}

	myHandlers = append(myHandlers,
		p.handleResize(),       // Handle resize events.
		p.handleCommands(cmds), // Process commands.
	)

	// Run event loop, handle updates and draw.
	var err error
	p.model, err = p.eventLoop(p.model, cmds)
	killed := p.ctx.Err() != nil
	if killed {
		err = ErrProgramKilled
	} else {
		// Ensure we rendered the final state of the model.
		p.renderer.reset()
		p.vb.clear()
		p.model.View(p.vb)
		p.renderer.Write(p.vb.Render())
	}

	// Tear down.
	p.cancel()

	// Check if the cancel reader has been setup before waiting and closing.
	if p.cancelReader != nil {
		// Wait for input loop to finish.
		if p.cancelReader.Cancel() {
			p.waitForReadLoop()
		}
		_ = p.cancelReader.Close()
	}

	// Wait for all handlers to finish.
	handlersShutdown(myHandlers)

	// Restore terminal state.
	p.shutdown(killed)

	return p.model, err
}

// Send sends a message to the main update function, effectively allowing
// messages to be injected from outside the program for interoperability
// purposes.
//
// If the program hasn't started yet this will be a blocking operation.
// If the program has already been terminated this will be a no-op, so it's safe
// to send messages after the program has exited.
func (p *Program[M]) Send(msg Msg) {
	select {
	case <-p.ctx.Done():
	case p.msgs <- msg:
	}
}

// Quit is a convenience function for quitting Tea programs. Use it
// when you need to shut down a Tea program from the outside.
//
// If you wish to quit from within a Tea program use the Quit command.
//
// If the program is not running this will be a no-op, so it's safe to call
// if the program is unstarted or has already exited.
func (p *Program[M]) Quit() {
	p.Send(Quit())
}

// Kill stops the program immediately and restores the former terminal state.
// The final render that you would normally see when quitting will be skipped.
// [program.Run] returns a [ErrProgramKilled] error.
func (p *Program[M]) Kill() {
	p.cancel()
}

// Wait waits/blocks until the underlying Program finished shutting down.
func (p *Program[M]) Wait() {
	<-p.done
}

// shutdown performs operations to free up resources and restore the terminal
// to its original state.
func (p *Program[M]) shutdown(kill bool) {
	if p.renderer != nil {
		p.renderer.stop(!kill)
	}

	_ = p.restoreTerminalState()
	if p.restoreOutput != nil {
		_ = p.restoreOutput()
	}
	p.done <- struct{}{}
}

// ReleaseTerminal restores the original terminal state and cancels the input reader.
// You can return control to the Program with RestoreTerminal.
func (p *Program[M]) ReleaseTerminal() error {
	p.ignoreSignals = true
	p.cancelReader.Cancel()
	p.waitForReadLoop()

	if p.renderer != nil {
		p.renderer.stop(true)
	}

	p.altScreenWasActive = p.renderer.altScreen()
	return p.restoreTerminalState()
}

// RestoreTerminal reinitializes the Program's input reader, restores the
// terminal to the former state when the program was running, and repaints.
// Use it to reinitialize a Program after running ReleaseTerminal.
func (p *Program[M]) RestoreTerminal() error {
	p.ignoreSignals = false

	if err := p.initTerminal(); err != nil {
		return err
	}
	if err := p.initCancelReader(); err != nil {
		return err
	}

	if p.altScreenWasActive {
		p.renderer.enterAltScreen()
	} else {
		// entering alt screen already causes a repaint.
		go p.Send(msgRepaint{})
	}
	p.renderer.start()

	// If the output is a terminal, it may have been resized while another
	// process was at the foreground, in which case we may not have received SIGWINCH.
	// Detect any size change now and propagate the new size as needed.
	go p.checkResize()

	return nil
}
