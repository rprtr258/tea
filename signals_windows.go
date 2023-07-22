//go:build windows
// +build windows

package tea

// listenForResize is not available on windows because windows does not
// implement syscall.SIGWINCH.
func (p *Program[M]) listenForResize(done chan struct{}) {
	close(done)
}
