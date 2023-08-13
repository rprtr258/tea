package tea

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	isatty "github.com/mattn/go-isatty"
	localereader "github.com/mattn/go-localereader"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

func (p *Program[M]) initTerminal() error {
	err := p.initInput()
	if err != nil {
		return err
	}

	if p.console != nil {
		err = p.console.SetRaw()
		if err != nil {
			return fmt.Errorf("error entering raw mode: %w", err)
		}
	}

	p.renderer.setCursor(false)
	return nil
}

// restoreTerminalState restores the terminal to the state prior to running the
// Bubble Tea program.
func (p *Program[M]) restoreTerminalState() error {
	if p.renderer != nil {
		p.renderer.setCursor(true)
		p.renderer.setMouseCellMotion(false)
		p.renderer.setMouseAllMotion(false)

		if p.renderer.altScreen() {
			p.renderer.exitAltScreen()

			// give the terminal a moment to catch up
			time.Sleep(time.Millisecond * 10) //nolint:gomnd
		}
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
	case <-time.After(500 * time.Millisecond): //nolint:gomnd
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
