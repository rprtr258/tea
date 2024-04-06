package file_picker2 //nolint:revive,stylecheck

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/iter"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/headless/hierachy"
	"github.com/rprtr258/tea/styles"
)

type dirEntry struct {
	fullpath string
	name     string
	perms    fs.FileMode
}

func readDir(path string) hierachy.Node[dirEntry] {
	// info, _ := os.Stat(path)
	subdirs, _ := os.ReadDir(path)
	return hierachy.Node[dirEntry]{
		Value: dirEntry{
			fullpath: path,
			name:     filepath.Base(path),
			// perms: info.Mode(),
		},
		Children: fun.FilterMap[hierachy.Node[dirEntry]](func(e fs.DirEntry) (hierachy.Node[dirEntry], bool) {
			if strings.HasPrefix(e.Name(), ".git") {
				return hierachy.Node[dirEntry]{}, false
			}
			return readDir(filepath.Join(path, e.Name())), true
		}, subdirs...),
	}
}

type model struct {
	tree *hierachy.Hierachy[dirEntry]
}

func (m *model) Init(yield func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			yield(tea.Quit)
		case "j":
			m.tree.GoNextOrUp()
		case "k":
			m.tree.GoPrevOrUp()
		case "h":
			if m.tree.IsCollapsed() {
				m.tree.GoUp()
			} else {
				m.tree.ToggleCollapsed()
			}
		case "l":
			if m.tree.IsCollapsed() {
				m.tree.ToggleCollapsed()
			} else {
				m.tree.GoDown()
			}
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	selected := 0
	m.tree.Iter(func(i hierachy.IterItem[dirEntry]) bool {
		selected++
		return !i.IsSelected
	})

	height := vb.Height
	iter.Skip(m.tree.Iter, max(0, selected-height/2))(func(i hierachy.IterItem[dirEntry]) bool {
		vbItem := vb.PaddingLeft(i.Depth * 2)
		if i.IsSelected {
			vbItem = vbItem.Styled(styles.Style{}.Foreground(styles.FgColor("170")))
		}
		if i.HasChildren {
			vbArrow := vbItem.Styled(styles.Style{}.Foreground(styles.FgColor("169")))
			if i.IsCollapsed {
				vbArrow.WriteLineX(">")
			} else {
				vbArrow.WriteLineX("∨")
			}
			vbItem = vbItem.PaddingLeft(1)
		}
		vbItem.WriteLine(i.Value.name)

		vb = vb.PaddingTop(1)
		height--
		return height != 0
	})
}

func Main(ctx context.Context) error {
	m := &model{
		tree: hierachy.New(readDir(".")),
	}
	_, err := tea.NewProgram(ctx, m).WithOutput(os.Stderr).Run()
	return err
}
