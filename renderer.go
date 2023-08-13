package tea

// Renderer is the interface for Bubble Tea renderers.
type Renderer interface { //nolint:interfacebloat
	// Write a frame to the renderer. The renderer can write this data to
	// output at its discretion.
	Write(string)

	// start the renderer
	start()
	// stop the renderer, but render the final frame in the buffer, if any
	stop()
	// kill - stop the renderer without doing any final rendering
	kill()

	// repaint - request a full re-render.
	// Note that this will not trigger a render immediately.
	// Rather, this method causes the next render to be a full repaint.
	// Because of this, it's safe to call this method multiple times.
	repaint()
	// clearScreen - clears the terminal
	clearScreen()

	// altScreen - whether or not the alternate screen buffer is enabled
	altScreen() bool
	// enterAltScreen - enable the alternate screen buffer
	enterAltScreen()
	// exitAltScreen - disable the alternate screen buffer
	exitAltScreen()

	// setCursor - if true, show curser, if false, hide cursor
	setCursor(bool)

	// setMouseCellMotion - if true, enable mouse click, release, wheel and
	// motion events if a mouse button is pressed (i.e., drag events).
	// If false, disables Mouse Cell Motion tracking.
	setMouseCellMotion(bool)

	// setMouseAllMotion - if true, enables mouse click, release, wheel and
	// motion events, regardless of whether a mouse button is pressed.
	// Many modern terminals support this, but not all.
	// If false, disables All Motion mouse tracking.
	setMouseAllMotion(bool)

	handleMessages(Msg)
}

// msgRepaint forces a full repaint.
type msgRepaint struct{}
