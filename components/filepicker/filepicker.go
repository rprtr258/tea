package filepicker

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unsafe"

	"github.com/dustin/go-humanize"
	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/key"
	"github.com/rprtr258/tea/styles"
)

func pop(data []int) ([]int, int) {
	return data[:len(data)-1], data[len(data)-1]
}

// TODO: accept fs.FS, test w/ S3
// New returns a new filepicker model with default styling and key bindings.
func New() Model {
	return Model{
		CurrentDirectory: ".",
		Cursor:           ">",
		FileAllowed:      true,
		AutoHeight:       true,
		KeyMap:           DefaultKeyMap,
		Styles:           DefaultStyles,
	}
}

type msgError struct{ err error }

type msgReadDir struct {
	id      uintptr
	entries []os.DirEntry
}

const (
	_marginBottom  = 5
	_fileSizeWidth = 8
	_paddingLeft   = 2
)

// KeyMap defines key bindings for each user action.
type KeyMap struct {
	GoToTop  key.Binding
	GoToLast key.Binding
	Down     key.Binding
	Up       key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Back     key.Binding
	Open     key.Binding
	Select   key.Binding
}

// DefaultKeyMap defines the default keybindings.
var DefaultKeyMap = KeyMap{
	GoToTop:  key.Binding{Keys: []string{"g"}, Help: key.Help{"g", "first"}},
	GoToLast: key.Binding{Keys: []string{"G"}, Help: key.Help{"G", "last"}},
	Down:     key.Binding{Keys: []string{"j", "down", "ctrl+n"}, Help: key.Help{"j", "down"}},
	Up:       key.Binding{Keys: []string{"k", "up", "ctrl+p"}, Help: key.Help{"k", "up"}},
	PageUp:   key.Binding{Keys: []string{"K", "pgup"}, Help: key.Help{"pgup", "page up"}},
	PageDown: key.Binding{Keys: []string{"J", "pgdown"}, Help: key.Help{"pgdown", "page down"}},
	Back:     key.Binding{Keys: []string{"h", "backspace", "left", "esc"}, Help: key.Help{"h", "back"}},
	Open:     key.Binding{Keys: []string{"l", "right", "enter"}, Help: key.Help{"l", "open"}},
	Select:   key.Binding{Keys: []string{"enter"}, Help: key.Help{"enter", "select"}},
}

// Styles defines the possible customizations for styles in the file picker.
type Styles struct {
	DisabledCursor   styles.Style
	Cursor           styles.Style
	Symlink          styles.Style
	Directory        styles.Style
	File             styles.Style
	DisabledFile     styles.Style
	Permission       styles.Style
	Selected         styles.Style
	DisabledSelected styles.Style
	FileSize         styles.Style
	EmptyDirectory   styles.Style
}

// DefaultStyles is default styling for the file picker, with a given Lip Gloss renderer.
var DefaultStyles = Styles{
	DisabledCursor:   styles.Style{}.Foreground(scuf.FgANSI(247)),
	Cursor:           styles.Style{}.Foreground(scuf.FgANSI(212)),
	Symlink:          styles.Style{}.Foreground(scuf.FgANSI(36)),
	Directory:        styles.Style{}.Foreground(scuf.FgANSI(99)),
	File:             styles.Style{},
	DisabledFile:     styles.Style{}.Foreground(scuf.FgANSI(243)),
	DisabledSelected: styles.Style{}.Foreground(scuf.FgANSI(247)),
	Permission:       styles.Style{}.Foreground(scuf.FgANSI(244)),
	Selected:         styles.Style{}.Foreground(scuf.FgANSI(212)).Bold(true),
	FileSize:         styles.Style{}.Foreground(scuf.FgANSI(240)).Align(styles.Right),
	EmptyDirectory:   styles.Style{}.Foreground(scuf.FgANSI(240)),
}

// Model represents a file picker.
type Model struct {
	// Path is the path which the user has selected with the file picker.
	Path string

	// CurrentDirectory is the directory that the user is currently in.
	CurrentDirectory string

	// AllowedTypes specifies which file types the user may select.
	// If empty the user may select any file.
	AllowedTypes []string

	KeyMap      KeyMap
	files       []os.DirEntry
	ShowHidden  bool
	DirAllowed  bool
	FileAllowed bool

	fileSelected  string
	selected      int
	selectedStack []int

	min      int
	max      int
	maxStack []int
	minStack []int

	Height     int
	AutoHeight bool

	Cursor string
	Styles Styles
}

func (m *Model) pushView() {
	m.minStack = append(m.minStack, m.min)
	m.maxStack = append(m.maxStack, m.max)
	m.selectedStack = append(m.selectedStack, m.selected)
}

func (m *Model) popView() (selected, min, max int) {
	m.selectedStack, selected = pop(m.selectedStack)
	m.minStack, min = pop(m.selectedStack)
	m.maxStack, max = pop(m.selectedStack)
	return
}

// isHidden reports whether a file is hidden or not.
func isHidden(file string) bool {
	return strings.HasPrefix(file, ".")
}

func (m *Model) cmdReadDir(path string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return msgError{err}
		}

		slices.SortFunc(dirEntries, func(i, j fs.DirEntry) int {
			switch {
			case i.IsDir() == j.IsDir():
				return cmp.Compare(i.Name(), j.Name())
			case i.IsDir():
				return -1
			default:
				return 1
			}
		})

		return msgReadDir{
			id: uintptr(unsafe.Pointer(m)),
			entries: fun.
				If(showHidden, dirEntries).
				ElseF(func() []fs.DirEntry {
					return fun.Filter(
						func(dirEntry os.DirEntry) bool {
							return !isHidden(dirEntry.Name())
						}, dirEntries...)
				}),
		}
	}
}

// Init initializes the file picker model.
func (m *Model) Init(yield func(...tea.Cmd)) {
	yield(m.cmdReadDir(m.CurrentDirectory, m.ShowHidden))
}

// Update handles user interactions within the file picker model.
func (m *Model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case msgReadDir:
		if msg.id != uintptr(unsafe.Pointer(m)) {
			break
		}
		m.files = msg.entries
		m.max = m.Height - 1
	case tea.MsgWindowSize:
		if m.AutoHeight {
			m.Height = msg.Height - _marginBottom
		}
		m.max = m.Height - 1
	case tea.MsgKey:
		switch {
		case key.Matches(msg, m.KeyMap.GoToTop):
			m.selected = 0
			m.min = 0
			m.max = m.Height - 1
		case key.Matches(msg, m.KeyMap.GoToLast):
			m.selected = len(m.files) - 1
			m.min = len(m.files) - m.Height
			m.max = len(m.files) - 1
		case key.Matches(msg, m.KeyMap.Down):
			m.selected = min(m.selected+1, len(m.files)-1)
			if m.selected > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.KeyMap.Up):
			m.selected = max(m.selected-1, 0)
			if m.selected < m.min {
				m.min--
				m.max--
			}
		case key.Matches(msg, m.KeyMap.PageDown):
			m.selected = min(m.selected+m.Height, len(m.files)-1)
			m.min += m.Height
			m.max += m.Height

			if m.max >= len(m.files) {
				m.max = len(m.files) - 1
				m.min = m.max - m.Height
			}
		case key.Matches(msg, m.KeyMap.PageUp):
			m.selected = max(m.selected-m.Height, 0)
			m.min -= m.Height
			m.max -= m.Height

			if m.min < 0 {
				m.min = 0
				m.max = m.min + m.Height
			}
		case key.Matches(msg, m.KeyMap.Back):
			m.CurrentDirectory = filepath.Dir(m.CurrentDirectory)
			if len(m.selectedStack) > 0 {
				m.selected, m.min, m.max = m.popView()
			} else {
				m.selected = 0
				m.min = 0
				m.max = m.Height - 1
			}
			return []tea.Cmd{m.cmdReadDir(m.CurrentDirectory, m.ShowHidden)}
		case key.Matches(msg, m.KeyMap.Open):
			if len(m.files) == 0 {
				break
			}

			f := m.files[m.selected]
			info, err := f.Info()
			if err != nil {
				break
			}
			isSymlink := info.Mode()&os.ModeSymlink != 0
			isDir := f.IsDir()

			if isSymlink {
				symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.CurrentDirectory, f.Name()))
				info, err := os.Stat(symlinkPath)
				if err != nil {
					break
				}
				if info.IsDir() {
					isDir = true
				}
			}

			if !isDir && m.FileAllowed || isDir && m.DirAllowed {
				if key.Matches(msg, m.KeyMap.Select) {
					// Select the current path as the selection
					m.Path = filepath.Join(m.CurrentDirectory, f.Name())
				}
			}

			if !isDir {
				break
			}

			m.CurrentDirectory = filepath.Join(m.CurrentDirectory, f.Name())
			m.pushView()
			m.selected = 0
			m.min = 0
			m.max = m.Height - 1
			return []tea.Cmd{m.cmdReadDir(m.CurrentDirectory, m.ShowHidden)}
		}
	}
	return nil
}

// View returns the view of the file picker.
func (m *Model) View(vb tea.Viewbox) {
	if len(m.files) == 0 {
		vb.Styled(m.Styles.EmptyDirectory).WriteLine("Bummer. No Files Found.")
		return
	}

	for i, f := range m.files {
		if i < m.min || i > m.max {
			continue
		}

		info, _ := f.Info()
		isSymlink := info.Mode()&os.ModeSymlink != 0
		size := humanize.Bytes(uint64(info.Size()))
		size = strings.Repeat(" ", _fileSizeWidth-len(size)) + size
		name := f.Name()

		disabled := !m.canSelect(name) && !f.IsDir()
		fileName := fun.If(!isSymlink, name).ElseF(func() string {
			symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.CurrentDirectory, name))
			return fmt.Sprintf("%s → %s", name, symlinkPath)
		})

		vb0 := vb.Row(i)
		if m.selected == i {
			styleCursor := fun.IF(disabled, m.Styles.DisabledSelected, m.Styles.Cursor)
			styleSelected := fun.IF(disabled, m.Styles.DisabledSelected, m.Styles.Selected)
			vb0.Styled(styleCursor).WriteLine(m.Cursor)
			vb0 = vb0.PaddingLeft(2).Styled(styleSelected)
			vb0.WriteLine(info.Mode().String())
			vb0 = vb0.PaddingLeft(1 + 3*3 + 1) // type(1) + perms(3*3) + space(1)
			vb0.WriteLine(size)
			vb0 = vb0.PaddingLeft(_fileSizeWidth + 1)
			vb0.WriteLine(fileName)
			continue
		}

		vb0 = vb0.PaddingLeft(2)

		vb0.Styled(m.Styles.Permission).WriteLine(info.Mode().String())
		vb0 = vb0.PaddingLeft(1 + 3*3 + 1) // type(1) + perms(3*3) + space(1)
		vb0.Styled(m.Styles.FileSize).WriteLine(size)
		vb0 = vb0.PaddingLeft(_fileSizeWidth + 1)
		vb0.
			Styled(fun.Switch(true, m.Styles.File).
				Case(f.IsDir(), m.Styles.Directory).
				Case(isSymlink, m.Styles.Symlink).
				Case(disabled, m.Styles.DisabledFile).
				End()).
			WriteLine(fileName)
	}
}

// DidSelectFile returns whether a user has selected a file (on this msg).
func (m *Model) DidSelectFile(msg tea.Msg) (bool, string) {
	didSelect, path := m.didSelectFile(msg)
	return didSelect && m.canSelect(path), path
}

// DidSelectDisabledFile returns whether a user tried to select a disabled file
// (on this msg). This is necessary only if you would like to warn the user that
// they tried to select a disabled file.
func (m *Model) DidSelectDisabledFile(msg tea.Msg) (bool, string) {
	didSelect, path := m.didSelectFile(msg)
	if didSelect && !m.canSelect(path) {
		return true, path
	}
	return false, ""
}

func (m *Model) didSelectFile(msg tea.Msg) (bool, string) {
	if len(m.files) == 0 {
		return false, ""
	}

	switch msg := msg.(type) {
	case tea.MsgKey:
		// If the msg does not match the Select keymap then this could not have been a selection.
		if !key.Matches(msg, m.KeyMap.Select) {
			return false, ""
		}

		// The key press was a selection, let's confirm whether the current file could
		// be selected or used for navigating deeper into the stack.
		f := m.files[m.selected]
		info, err := f.Info()
		if err != nil {
			return false, ""
		}

		isDir := f.IsDir()
		if info.Mode()&os.ModeSymlink != 0 { // is symlink
			symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.CurrentDirectory, f.Name()))
			info, err := os.Stat(symlinkPath)
			if err != nil {
				return false, ""
			}

			if info.IsDir() {
				isDir = true
			}
		}

		if !isDir && m.FileAllowed || isDir && m.DirAllowed && m.Path != "" {
			return true, m.Path
		}

		// If the msg was not a MsgKey, then the file could not have been selected this iteration.
		// Only a MsgKey can select a file.
		return false, ""
	default:
		return false, ""
	}
}

func (m *Model) canSelect(file string) bool {
	if len(m.AllowedTypes) == 0 {
		return true
	}

	for _, ext := range m.AllowedTypes {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}

	return false
}
