package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/bubbles/list"
	"github.com/rprtr258/tea/examples/altscreen_toggle"
	"github.com/rprtr258/tea/examples/cellbuffer"
	"github.com/rprtr258/tea/examples/chat"
	"github.com/rprtr258/tea/examples/composable_views"
	"github.com/rprtr258/tea/examples/credit_card_form"
	"github.com/rprtr258/tea/examples/debounce"
	"github.com/rprtr258/tea/examples/exec"
	"github.com/rprtr258/tea/examples/file_picker"
	"github.com/rprtr258/tea/examples/fullscreen"
	"github.com/rprtr258/tea/examples/glamour/custom_renderer"
	"github.com/rprtr258/tea/examples/glamour/helloworld"
	"github.com/rprtr258/tea/examples/glamour/menu"
	"github.com/rprtr258/tea/examples/help"
	"github.com/rprtr258/tea/examples/http"
	"github.com/rprtr258/tea/examples/lipgloss/layout"
	"github.com/rprtr258/tea/examples/lipgloss/ssh"
	"github.com/rprtr258/tea/examples/list_default"
	"github.com/rprtr258/tea/examples/list_fancy"
	"github.com/rprtr258/tea/examples/list_simple"
	"github.com/rprtr258/tea/examples/mouse"
	"github.com/rprtr258/tea/examples/package_manager"
	"github.com/rprtr258/tea/examples/pager"
	"github.com/rprtr258/tea/examples/paginator"
	"github.com/rprtr258/tea/examples/pipe"
	"github.com/rprtr258/tea/examples/prevent_quit"
	"github.com/rprtr258/tea/examples/progress_animated"
	"github.com/rprtr258/tea/examples/progress_download"
	"github.com/rprtr258/tea/examples/progress_static"
	"github.com/rprtr258/tea/examples/realtime"
	"github.com/rprtr258/tea/examples/result"
	"github.com/rprtr258/tea/examples/send_msg"
	"github.com/rprtr258/tea/examples/sequence"
	"github.com/rprtr258/tea/examples/simple"
	"github.com/rprtr258/tea/examples/spinner"
	"github.com/rprtr258/tea/examples/spinners"
	"github.com/rprtr258/tea/examples/split_editors"
	"github.com/rprtr258/tea/examples/stopwatch"
	"github.com/rprtr258/tea/examples/table"
	"github.com/rprtr258/tea/examples/tabs"
	"github.com/rprtr258/tea/examples/textarea"
	"github.com/rprtr258/tea/examples/textinput"
	"github.com/rprtr258/tea/examples/textinputs"
	"github.com/rprtr258/tea/examples/timer"
	"github.com/rprtr258/tea/examples/tui_daemon_combo"
	"github.com/rprtr258/tea/examples/tutorials/basics"
	"github.com/rprtr258/tea/examples/tutorials/commands"
	"github.com/rprtr258/tea/examples/views"
	"github.com/rprtr258/tea/lipgloss"
)

type examples map[string]func(context.Context) error

var (
	teaExamples = examples{
		"altscreen-toggle":  altscreen_toggle.Main,
		"cellbuffer":        cellbuffer.Main,
		"chat":              chat.Main,
		"composable-views":  composable_views.Main,
		"credit-card-form":  credit_card_form.Main,
		"debounce":          debounce.Main,
		"exec":              exec.Main,
		"file-picker":       file_picker.Main,
		"fullscreen":        fullscreen.Main,
		"help":              help.Main,
		"http":              http.Main,
		"list-default":      list_default.Main,
		"list-fancy":        list_fancy.Main,
		"list-simple":       list_simple.Main,
		"mouse":             mouse.Main,
		"package-manager":   package_manager.Main,
		"pager":             pager.Main,
		"paginator":         paginator.Main,
		"pipe":              pipe.Main,
		"prevent-quit":      prevent_quit.Main,
		"progress-animated": progress_animated.Main,
		"progress-download": progress_download.Main,
		"progress-static":   progress_static.Main,
		"realtime":          realtime.Main,
		"result":            result.Main,
		"send-msg":          send_msg.Main,
		"sequence":          sequence.Main,
		"simple":            simple.Main,
		"spinner":           spinner.Main,
		"spinners":          spinners.Main,
		"split-editors":     split_editors.Main,
		"stopwatch":         stopwatch.Main,
		"table":             table.Main,
		"tabs":              tabs.Main,
		"textarea":          textarea.Main,
		"textinput":         textinput.Main,
		"textinputs":        textinputs.Main,
		"timer":             timer.Main,
		"tui-daemon-combo":  tui_daemon_combo.Main,
		"views":             views.Main,
	}
	tutorials = examples{
		"basics":   basics.Main,
		"commands": commands.Main,
	}
	lipglossExamples = examples{
		"layout": layout.Main,
		"ssh":    ssh.Main,
	}
	glamourExamples = examples{
		"custom-renderer": custom_renderer.Main,
		"helloworld":      helloworld.Main,
		"menu":            menu.Main,
	}
)

type item struct {
	name string
	main func(context.Context) error
}

func (i item) FilterValue() string { return i.name }

var (
	_styleTitle        = lipgloss.NewStyle().MarginLeft(2)
	_styleItem         = lipgloss.NewStyle().PaddingLeft(4)
	_styleItemSelected = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	_stylePagination   = list.DefaultStyle.PaginationStyle.PaddingLeft(4)
	_styleHelp         = list.DefaultStyle.HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                                     { return 1 }
func (d itemDelegate) Spacing() int                                    { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model[item]) []tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m *list.Model[item], index int, i item) {
	var res string
	if index == m.Index() {
		res = _styleItemSelected.Render("> " + i.name)
	} else {
		res = _styleItem.Render(i.name)
	}

	fmt.Fprint(w, res)
}

type model struct {
	list     list.Model[item]
	choice   item
	quitting bool
}

func (m *model) Init() []tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.list.SetWidth(msg.Width)
		return nil
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return []tea.Cmd{tea.Quit}
		case "enter":
			i, ok := m.list.SelectedItem()
			if ok {
				m.choice = i
			}
			return []tea.Cmd{tea.ExitAltScreen, tea.ClearScreen, tea.Quit}
		}
	}

	return m.list.Update(msg)
}

func (m *model) View(r tea.Renderer) {
	if m.choice.main != nil || m.quitting {
		return
	}

	r.Write("\n" + m.list.View())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runExamplesList(ctx context.Context, title string, examples examples) error {
	items := lo.MapToSlice(
		examples,
		func(name string, main func(context.Context) error) item {
			return item{
				name: name,
				main: main,
			}
		},
	)
	sort.Slice(items, func(i, j int) bool {
		return items[i].name < items[j].name
	})

	const (
		listHeight   = 30
		defaultWidth = 20
	)

	l := list.New[item](items, itemDelegate{}, defaultWidth, min(listHeight, len(items)+8))
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = _styleTitle
	l.Styles.PaginationStyle = _stylePagination
	l.Styles.HelpStyle = _styleHelp

	m, err := tea.NewProgram(ctx, &model{list: l}).Run()
	if err != nil {
		return err
	}

	if m.choice.main == nil {
		return nil
	}

	return m.choice.main(ctx)
}

func main() {
	if err := (&cli.App{
		Name:  "tea examples",
		Usage: "tea <example>",
		Action: func(ctx *cli.Context) error {
			return runExamplesList(ctx.Context, "Tea examples", teaExamples)
		},
		Commands: []*cli.Command{
			{
				Name:  "tutorials",
				Usage: "Tea tutorials",
				Action: func(ctx *cli.Context) error {
					return runExamplesList(ctx.Context, "Tea tutorials", tutorials)
				},
			},
			{
				Name:  "lipgloss",
				Usage: "lipgloss examples",
				Action: func(ctx *cli.Context) error {
					return runExamplesList(ctx.Context, "Lipgloss examples", lipglossExamples)
				},
			},
			{
				Name:  "glamour",
				Usage: "glamour examples",
				Action: func(ctx *cli.Context) error {
					return runExamplesList(ctx.Context, "Glamour examples", glamourExamples)
				},
			},
		},
	}).Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
