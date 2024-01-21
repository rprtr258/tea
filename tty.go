package tea

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/containerd/console"
	isatty "github.com/mattn/go-isatty"
	localereader "github.com/mattn/go-localereader"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

func (p *Program[M]) initInput() {
	// If input's a file, use console to manage it
	f, ok := p.input.(*os.File)
	if !ok {
		return
	}

	c, err := console.ConsoleFromFile(f)
	if err != nil {
		return // ignore error, this was just a test
	}

	p.console = c
}

// On unix systems, RestoreInput closes any TTYs we opened for input. Note that
// we don't do this on Windows as it causes the prompt to not be drawn until
// the terminal receives a keypress rather than appearing promptly after the
// program exits.
func (p *Program[M]) restoreInput() error {
	if p.console != nil {
		if err := p.console.Reset(); err != nil {
			return fmt.Errorf("error restoring console: %w", err)
		}
	}
	return nil
}

func openInputTTY() (*os.File, error) {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, fmt.Errorf("could not open a new TTY: %w", err)
	}
	return f, nil
}

func (p *Program[M]) initTerminal() error {
	p.initInput()

	if p.console != nil {
		if err := p.console.SetRaw(); err != nil {
			return fmt.Errorf("error entering raw mode: %w", err)
		}
	}

	p.renderer.setCursor(false)
	return nil
}

// restoreTerminalState restores the terminal to the state prior to running the
// Tea program.
func (p *Program[M]) restoreTerminalState() error {
	p.renderer.setCursor(true)
	p.renderer.setMouseCellMotion(false)
	p.renderer.setMouseAllMotion(false)

	if p.renderer.altScreen() {
		p.renderer.exitAltScreen()

		// give the terminal a moment to catch up
		time.Sleep(10 * time.Millisecond)
	}

	if p.console != nil {
		err := p.console.Reset()
		if err != nil {
			return fmt.Errorf("error restoring terminal state: %w", err)
		}
	}

	return p.restoreInput()
}

// initCancelReader (re)commences reading inputs.
func (p *Program[M]) initCancelReader() error {
	var err error
	p.cancelReader, err = cancelreader.NewReader(p.input)
	if err != nil {
		return fmt.Errorf("error creating cancelreader: %w", err)
	}

	p.readLoopDone = make(chan struct{})
	go p.readLoop()

	return nil
}

func (p *Program[M]) readLoop() {
	defer close(p.readLoopDone)

	input := localereader.NewReader(p.cancelReader)
	err := readInputs(p.ctx, p.msgs, input)
	if !errors.Is(err, io.EOF) && !errors.Is(err, cancelreader.ErrCanceled) {
		select {
		case <-p.ctx.Done():
		case p.errs <- err:
		}
	}
}

// waitForReadLoop waits for the cancelReader to finish its read loop.
func (p *Program[M]) waitForReadLoop() {
	select {
	case <-p.readLoopDone:
	case <-time.After(500 * time.Millisecond):
		// The read loop hangs, which means the input
		// cancelReader's cancel function has returned true even
		// though it was not able to cancel the read.
	}
}

// checkResize detects the current size of the output and informs the program
// via a MsgWindowSize.
func (p *Program[M]) checkResize() {
	f, ok := p.output.TTY().(*os.File)
	if !ok || !isatty.IsTerminal(f.Fd()) {
		// can't query window size
		return
	}

	w, h, err := term.GetSize(int(f.Fd()))
	if err != nil {
		select {
		case <-p.ctx.Done():
		case p.errs <- err:
		}

		return
	}

	p.Send(MsgWindowSize{
		Width:  w,
		Height: h,
	})
}
