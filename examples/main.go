package main

import (
	"log"
	"os"

	altscreentoggle "examples/altscreen-toggle"
	"examples/cellbuffer"
	"examples/chat"
	composableviews "examples/composable-views"
	creditcardform "examples/credit-card-form"
	"examples/debounce"
	"examples/exec"
	filepicker "examples/file-picker"
	"examples/fullscreen"
	"examples/glamour"
	"examples/help"
	"examples/http"
	listdefault "examples/list-default"
	listfancy "examples/list-fancy"
	listsimple "examples/list-simple"
	"examples/mouse"
	packagemanager "examples/package-manager"
	"examples/pager"
	"examples/paginator"
	"examples/pipe"
	preventquit "examples/prevent-quit"
	progressanimated "examples/progress-animated"
	progressdownload "examples/progress-download"
	progressstatic "examples/progress-static"
	"examples/realtime"
	"examples/result"
	sendmsg "examples/send-msg"
	"examples/sequence"
	"examples/simple"
	"examples/spinner"
	"examples/spinners"
	spliteditors "examples/split-editors"
	"examples/stopwatch"
	"examples/table"
	"examples/tabs"
	"examples/textarea"
	"examples/textinput"
	"examples/textinputs"
	"examples/timer"
	tuidaemoncombo "examples/tui-daemon-combo"
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
		altscreentoggle.Main()
	case "cellbuffer":
		cellbuffer.Main()
	case "chat":
		chat.Main()
	case "composable-views":
		composableviews.Main()
	case "credit-card-form":
		creditcardform.Main()
	case "debounce":
		debounce.Main()
	case "exec":
		exec.Main()
	case "file-picker":
		filepicker.Main()
	case "fullscreen":
		fullscreen.Main()
	case "glamour":
		glamour.Main()
	case "help":
		help.Main()
	case "http":
		http.Main()
	case "list-default":
		listdefault.Main()
	case "list-fancy":
		listfancy.Main()
	case "list-simple":
		listsimple.Main()
	case "mouse":
		mouse.Main()
	case "package-manager":
		packagemanager.Main()
	case "pager":
		pager.Main()
	case "paginator":
		paginator.Main()
	case "pipe":
		pipe.Main()
	case "prevent-quit":
		preventquit.Main()
	case "progress-animated":
		progressanimated.Main()
	case "progress-download":
		progressdownload.Main()
	case "progress-static":
		progressstatic.Main()
	case "realtime":
		realtime.Main()
	case "result":
		result.Main()
	case "send-msg":
		sendmsg.Main()
	case "sequence":
		sequence.Main()
	case "simple":
		simple.Main()
	case "spinner":
		spinner.Main()
	case "spinners":
		spinners.Main()
	case "split-editors":
		spliteditors.Main()
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
		tuidaemoncombo.Main()
	case "views":
		views.Main()
	default:
		log.Fatal("Unknown example")
	}
}
