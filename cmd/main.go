package main

import (
	"cmp"
	"context"
	"log"
	"os"
	"slices"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/cmd/altscreen_toggle"
	"github.com/rprtr258/tea/cmd/autocomplete"
	"github.com/rprtr258/tea/cmd/cellbuffer"
	"github.com/rprtr258/tea/cmd/chat"
	"github.com/rprtr258/tea/cmd/composable_views"
	"github.com/rprtr258/tea/cmd/credit_card_form"
	"github.com/rprtr258/tea/cmd/debounce"
	"github.com/rprtr258/tea/cmd/exec"
	"github.com/rprtr258/tea/cmd/file_picker"
	"github.com/rprtr258/tea/cmd/file_picker2"
	"github.com/rprtr258/tea/cmd/fullscreen"
	"github.com/rprtr258/tea/cmd/help"
	"github.com/rprtr258/tea/cmd/http"
	"github.com/rprtr258/tea/cmd/list_default"
	"github.com/rprtr258/tea/cmd/list_fancy"
	"github.com/rprtr258/tea/cmd/list_simple"
	"github.com/rprtr258/tea/cmd/markdown/custom_renderer"
	"github.com/rprtr258/tea/cmd/markdown/helloworld"
	"github.com/rprtr258/tea/cmd/markdown/menu"
	"github.com/rprtr258/tea/cmd/mouse"
	"github.com/rprtr258/tea/cmd/package_manager"
	"github.com/rprtr258/tea/cmd/pager"
	"github.com/rprtr258/tea/cmd/paginator"
	"github.com/rprtr258/tea/cmd/pipe"
	"github.com/rprtr258/tea/cmd/plasma"
	"github.com/rprtr258/tea/cmd/prevent_quit"
	"github.com/rprtr258/tea/cmd/progress_animated"
	"github.com/rprtr258/tea/cmd/progress_download"
	"github.com/rprtr258/tea/cmd/progress_static"
	"github.com/rprtr258/tea/cmd/realtime"
	"github.com/rprtr258/tea/cmd/result"
	"github.com/rprtr258/tea/cmd/send_msg"
	"github.com/rprtr258/tea/cmd/sequence"
	"github.com/rprtr258/tea/cmd/simple"
	"github.com/rprtr258/tea/cmd/spinner"
	"github.com/rprtr258/tea/cmd/spinners"
	"github.com/rprtr258/tea/cmd/split_editors"
	"github.com/rprtr258/tea/cmd/stopwatch"
	"github.com/rprtr258/tea/cmd/styles/layout"
	"github.com/rprtr258/tea/cmd/styles/ssh"
	"github.com/rprtr258/tea/cmd/table"
	"github.com/rprtr258/tea/cmd/table_chess"
	"github.com/rprtr258/tea/cmd/table_pokemon"
	"github.com/rprtr258/tea/cmd/tabs"
	"github.com/rprtr258/tea/cmd/textarea"
	"github.com/rprtr258/tea/cmd/textinput"
	"github.com/rprtr258/tea/cmd/textinputs"
	"github.com/rprtr258/tea/cmd/timer"
	"github.com/rprtr258/tea/cmd/tui_daemon_combo"
	"github.com/rprtr258/tea/cmd/tutorials/basics"
	"github.com/rprtr258/tea/cmd/tutorials/commands"
	"github.com/rprtr258/tea/cmd/views"
	"github.com/rprtr258/tea/components/list"
	"github.com/rprtr258/tea/styles"
)

type examples = map[string]func(context.Context) error

var (
	teaExamples = examples{
		"altscreen-toggle":  altscreen_toggle.Main,
		"autcomplete":       autocomplete.Main,
		"cellbuffer":        cellbuffer.Main,
		"chat":              chat.Main,
		"composable-views":  composable_views.Main,
		"credit-card-form":  credit_card_form.Main,
		"debounce":          debounce.Main,
		"exec":              exec.Main,
		"file-picker":       file_picker.Main,
		"file-picker2":      file_picker2.Main,
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
		"plasma":            plasma.Main,
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
		"tablebox/pokemon":  table_pokemon.Main,
		"tablebox/chess":    table_chess.Main,
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
	stylesExamples = examples{
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
	_styleTitle        = styles.Style{} // .MarginLeft(2)
	_styleItem         = styles.Style{}
	_styleItemSelected = styles.Style{}.Foreground(styles.FgColor("170"))
	_stylePagination   = list.DefaultStyle.PaginationStyle // .PaddingLeft(4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                                     { return 1 }
func (d itemDelegate) Spacing() int                                    { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model[item]) []tea.Cmd { return nil }
func (d itemDelegate) Render(vb tea.Viewbox, m *list.Model[item], index int, i item) {
	vb = vb.PaddingLeft(2)
	var style styles.Style
	if index == m.Index() {
		style = _styleItemSelected
		vb.Styled(_styleItemSelected).WriteLine("> ")
	} else {
		style = _styleItem
	}
	vb = vb.PaddingLeft(2)
	vb.Styled(style).WriteLine(i.name)
}

type model struct {
	list     list.Model[item]
	choice   item
	quitting bool
}

func (m *model) Init(tea.Context[*model]) {}

func (m *model) Update(c tea.Context[*model], msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.MsgWindowSize:
		m.list.SetWidth(msg.Width)
		return
	case tea.MsgKey:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			c.Dispatch(tea.Quit)
			return
		case "enter":
			i, ok := m.list.SelectedItem()
			if ok {
				m.choice = i
			}
			c.Dispatch(tea.ExitAltScreen, tea.ClearScreen, tea.Quit)
			return
		}
	}

	ctxList := tea.Of(c, func(m *model) *list.Model[item] { return &m.list })
	m.list.Update(ctxList, msg)
}

func (m *model) View(vb tea.Viewbox) {
	if m.choice.main != nil || m.quitting {
		return
	}

	m.list.View(vb.Padding(tea.PaddingOptions{
		Top:    1,
		Left:   2,
		Bottom: 1,
	}))
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
	slices.SortFunc(items, func(i, j item) int {
		return cmp.Compare(i.name, j.name)
	})

	const (
		listHeight   = 30
		defaultWidth = 20
	)

	l := list.New(items, itemDelegate{}, defaultWidth, min(listHeight, len(items)+8))
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = _styleTitle
	l.Styles.PaginationStyle = _stylePagination

	m, err := tea.NewProgram2(ctx, &model{list: l}).Run()
	if err != nil {
		return err
	}

	if m.M.choice.main == nil {
		return nil
	}

	return m.M.choice.main(ctx)
}

func main() {
	app := &cli.App{
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
				Name:  "styles",
				Usage: "styles examples",
				Action: func(ctx *cli.Context) error {
					return runExamplesList(ctx.Context, "styles examples", stylesExamples)
				},
			},
			{
				Name:  "markdown",
				Usage: "markdown examples",
				Action: func(ctx *cli.Context) error {
					return runExamplesList(ctx.Context, "Markdown examples", glamourExamples)
				},
			},
		},
	}
	for name, f := range teaExamples {
		app.Commands = append(app.Commands, &cli.Command{
			Name:  name,
			Usage: name,
			Action: func(ctx *cli.Context) error {
				return f(ctx.Context)
			},
		})
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
