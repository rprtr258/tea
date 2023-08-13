package tea

type nilRenderer struct{}

func (nilRenderer) start()                  {}
func (nilRenderer) stop()                   {}
func (nilRenderer) kill()                   {}
func (nilRenderer) Write(string)            {}
func (nilRenderer) repaint()                {}
func (nilRenderer) clearScreen()            {}
func (nilRenderer) altScreen() bool         { return false }
func (nilRenderer) enterAltScreen()         {}
func (nilRenderer) exitAltScreen()          {}
func (nilRenderer) setCursor(bool)          {}
func (nilRenderer) setMouseCellMotion(bool) {}
func (nilRenderer) setMouseAllMotion(bool)  {}
func (nilRenderer) handleMessages(Msg)      {}
