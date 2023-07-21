package main

import (
	"log"
	"os"

	"github.com/rprtr258/tea/examples/altscreen_toggle"
	"github.com/rprtr258/tea/examples/cellbuffer"
	"github.com/rprtr258/tea/examples/chat"
	"github.com/rprtr258/tea/examples/composable_views"
	"github.com/rprtr258/tea/examples/credit_card_form"
	"github.com/rprtr258/tea/examples/debounce"
	"github.com/rprtr258/tea/examples/exec"
	"github.com/rprtr258/tea/examples/file_picker"
	"github.com/rprtr258/tea/examples/fullscreen"
	"github.com/rprtr258/tea/examples/glamour"
	"github.com/rprtr258/tea/examples/help"
	"github.com/rprtr258/tea/examples/http"
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
	"github.com/rprtr258/tea/examples/views"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run main.go <example>")
	}

	switch os.Args[1] {
	case "--help":
		log.Println("TODO: print examples list")
	case "altscreen-toggle":
		altscreen_toggle.Main()
	case "cellbuffer":
		cellbuffer.Main()
	case "chat":
		chat.Main()
	case "composable-views":
		composable_views.Main()
	case "credit-card-form":
		credit_card_form.Main()
	case "debounce":
		debounce.Main()
	case "exec":
		exec.Main()
	case "file-picker":
		file_picker.Main()
	case "fullscreen":
		fullscreen.Main()
	case "glamour":
		glamour.Main()
	case "help":
		help.Main()
	case "http":
		http.Main()
	case "list-default":
		list_default.Main()
	case "list-fancy":
		list_fancy.Main()
	case "list-simple":
		list_simple.Main()
	case "mouse":
		mouse.Main()
	case "package-manager":
		package_manager.Main()
	case "pager":
		pager.Main()
	case "paginator":
		paginator.Main()
	case "pipe":
		pipe.Main()
	case "prevent-quit":
		prevent_quit.Main()
	case "progress-animated":
		progress_animated.Main()
	case "progress-download":
		progress_download.Main()
	case "progress-static":
		progress_static.Main()
	case "realtime":
		realtime.Main()
	case "result":
		result.Main()
	case "send-msg":
		send_msg.Main()
	case "sequence":
		sequence.Main()
	case "simple":
		simple.Main()
	case "spinner":
		spinner.Main()
	case "spinners":
		spinners.Main()
	case "split-editors":
		split_editors.Main()
	case "stopwatch":
		stopwatch.Main()
	case "table":
		table.Main()
	case "tabs":
		tabs.Main()
	case "textarea":
		textarea.Main()
	case "textinput":
		textinput.Main()
	case "textinputs":
		textinputs.Main()
	case "timer":
		timer.Main()
	case "tui-daemon-combo":
		tui_daemon_combo.Main()
	case "views":
		views.Main()
	default:
		log.Fatalln("Unknown example")
	}
}
