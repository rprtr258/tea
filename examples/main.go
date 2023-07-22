package main

import (
	"fmt"
	"os"
	"sort"

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

var (
	_examples = map[string]func(){
		"altscreen-toggle":  altscreen_toggle.Main,
		"cellbuffer":        cellbuffer.Main,
		"chat":              chat.Main,
		"composable-views":  composable_views.Main,
		"credit-card-form":  credit_card_form.Main,
		"debounce":          debounce.Main,
		"exec":              exec.Main,
		"file-picker":       file_picker.Main,
		"fullscreen":        fullscreen.Main,
		"glamour":           glamour.Main,
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
	_examplesNames = examplesNames()
)

func examplesNames() []string {
	names := make([]string, 0, len(_examples))
	for name := range _examples {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func usage() {
	fmt.Println("Usage: go run main.go <example>")
	for _, example := range _examplesNames {
		fmt.Println(example)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	exampleMain, ok := _examples[os.Args[1]]
	if !ok {
		usage()
	}

	exampleMain()
}
