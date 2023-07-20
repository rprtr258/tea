package main

import (
	"log"
	"os"

	"examples/altscreen_toggle"
	"examples/cellbuffer"
	"examples/chat"
	"examples/composable_views"
	"examples/credit_card_form"
	"examples/debounce"
	"examples/exec"
	"examples/file_picker"
	"examples/fullscreen"
	"examples/glamour"
	"examples/help"
	"examples/http"
	"examples/list_default"
	"examples/list_fancy"
	"examples/list_simple"
	"examples/mouse"
	"examples/package_manager"
	"examples/pager"
	"examples/paginator"
	"examples/pipe"
	"examples/prevent_quit"
	"examples/progress_animated"
	"examples/progress_download"
	"examples/progress_static"
	"examples/realtime"
	"examples/result"
	"examples/send_msg"
	"examples/sequence"
	"examples/simple"
	"examples/spinner"
	"examples/spinners"
	"examples/split_editors"
	"examples/stopwatch"
	"examples/table"
	"examples/tabs"
	"examples/textarea"
	"examples/textinput"
	"examples/textinputs"
	"examples/timer"
	"examples/tui_daemon_combo"
	"examples/views"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <example>")
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
		log.Fatal("Unknown example")
	}
}
