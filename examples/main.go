package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

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
)

func exampleCommand(name string, main func()) *cli.Command {
	return &cli.Command{
		Name: name,
		Action: func(c *cli.Context) error {
			main()
			return nil
		},
	}
}

func main() {
	if err := (&cli.App{
		Name:  "tea examples",
		Usage: "tea <example>",
		Commands: []*cli.Command{
			exampleCommand("altscreen-toggle", altscreen_toggle.Main),
			exampleCommand("cellbuffer", cellbuffer.Main),
			exampleCommand("chat", chat.Main),
			exampleCommand("composable-views", composable_views.Main),
			exampleCommand("credit-card-form", credit_card_form.Main),
			exampleCommand("debounce", debounce.Main),
			exampleCommand("exec", exec.Main),
			exampleCommand("file-picker", file_picker.Main),
			exampleCommand("fullscreen", fullscreen.Main),
			exampleCommand("help", help.Main),
			exampleCommand("http", http.Main),
			exampleCommand("list-default", list_default.Main),
			exampleCommand("list-fancy", list_fancy.Main),
			exampleCommand("list-simple", list_simple.Main),
			exampleCommand("mouse", mouse.Main),
			exampleCommand("package-manager", package_manager.Main),
			exampleCommand("pager", pager.Main),
			exampleCommand("paginator", paginator.Main),
			exampleCommand("pipe", pipe.Main),
			exampleCommand("prevent-quit", prevent_quit.Main),
			exampleCommand("progress-animated", progress_animated.Main),
			exampleCommand("progress-download", progress_download.Main),
			exampleCommand("progress-static", progress_static.Main),
			exampleCommand("realtime", realtime.Main),
			exampleCommand("result", result.Main),
			exampleCommand("send-msg", send_msg.Main),
			exampleCommand("sequence", sequence.Main),
			exampleCommand("simple", simple.Main),
			exampleCommand("spinner", spinner.Main),
			exampleCommand("spinners", spinners.Main),
			exampleCommand("split-editors", split_editors.Main),
			exampleCommand("stopwatch", stopwatch.Main),
			exampleCommand("table", table.Main),
			exampleCommand("tabs", tabs.Main),
			exampleCommand("textarea", textarea.Main),
			exampleCommand("textinput", textinput.Main),
			exampleCommand("textinputs", textinputs.Main),
			exampleCommand("timer", timer.Main),
			exampleCommand("tui-daemon-combo", tui_daemon_combo.Main),
			exampleCommand("views", views.Main),
			{
				Name:  "tutorials",
				Usage: "Tea tutorials",
				Subcommands: []*cli.Command{
					exampleCommand("basics", basics.Main),
					exampleCommand("commands", commands.Main),
				},
			},
			{
				Name:  "lipgloss",
				Usage: "lipgloss examples",
				Subcommands: []*cli.Command{
					exampleCommand("layout", layout.Main),
					exampleCommand("ssh", ssh.Main),
				},
			},
			{
				Name:  "glamour",
				Usage: "glamour examples",
				Subcommands: []*cli.Command{
					exampleCommand("custom-renderer", custom_renderer.Main),
					exampleCommand("helloworld", helloworld.Main),
					exampleCommand("menu", menu.Main),
				},
			},
		},
	}).Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
