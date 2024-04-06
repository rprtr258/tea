package ssh

// This example demonstrates how to use a custom Lip Gloss renderer with Wish,
// a package for building custom SSH servers.
//
// The big advantage to using custom renderers here is that we can accurately
// detect the background color and color profile for each client and render
// against that accordingly.
//
// For details on wish see: https://github.com/charmbracelet/wish/

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/creack/pty"
	"github.com/muesli/termenv"

	lipgloss "github.com/rprtr258/tea/styles"
)

// Available styles.
type styles struct {
	bold          lipgloss.Style
	faint         lipgloss.Style
	italic        lipgloss.Style
	underline     lipgloss.Style
	strikethrough lipgloss.Style
	red           lipgloss.Style
	green         lipgloss.Style
	yellow        lipgloss.Style
	blue          lipgloss.Style
	magenta       lipgloss.Style
	cyan          lipgloss.Style
	gray          lipgloss.Style
}

// Create new styles against a given renderer.
func makeStyles() styles {
	return styles{
		bold:          lipgloss.Style{}.SetString("bold").Bold(true),
		faint:         lipgloss.Style{}.SetString("faint").Faint(),
		italic:        lipgloss.Style{}.SetString("italic").Italic(),
		underline:     lipgloss.Style{}.SetString("underline").Underline(),
		strikethrough: lipgloss.Style{}.SetString("strikethrough").Strikethrough(true),
		red:           lipgloss.Style{}.SetString("red").Foreground(lipgloss.FgColor("#E88388")),
		green:         lipgloss.Style{}.SetString("green").Foreground(lipgloss.FgColor("#A8CC8C")),
		yellow:        lipgloss.Style{}.SetString("yellow").Foreground(lipgloss.FgColor("#DBAB79")),
		blue:          lipgloss.Style{}.SetString("blue").Foreground(lipgloss.FgColor("#71BEF2")),
		magenta:       lipgloss.Style{}.SetString("magenta").Foreground(lipgloss.FgColor("#D290E4")),
		cyan:          lipgloss.Style{}.SetString("cyan").Foreground(lipgloss.FgColor("#66C2CD")),
		gray:          lipgloss.Style{}.SetString("gray").Foreground(lipgloss.FgColor("#B9BFCA")),
	}
}

// Bridge Wish and Termenv so we can query for a user's terminal capabilities.
type sshOutput struct {
	ssh.Session
	tty *os.File
}

func (s *sshOutput) Write(p []byte) (int, error) {
	return s.Session.Write(p)
}

func (s *sshOutput) Read(p []byte) (int, error) {
	return s.Session.Read(p)
}

func (s *sshOutput) Fd() uintptr {
	return s.tty.Fd()
}

type sshEnviron struct {
	environ []string
}

func (s *sshEnviron) Getenv(key string) string {
	for _, v := range s.environ {
		if strings.HasPrefix(v, key+"=") {
			return v[len(key)+1:]
		}
	}
	return ""
}

func (s *sshEnviron) Environ() []string {
	return s.environ
}

// Create a termenv.Output from the session.
func outputFromSession(sess ssh.Session) *termenv.Output {
	sshPty, _, _ := sess.Pty()
	_, tty, err := pty.Open()
	if err != nil {
		log.Fatalln(err.Error())
	}
	o := &sshOutput{
		Session: sess,
		tty:     tty,
	}
	environ := sess.Environ()
	environ = append(environ, fmt.Sprintf("TERM=%s", sshPty.Term))
	e := &sshEnviron{environ: environ}
	// We need to use unsafe mode here because the ssh session is not running
	// locally and we already know that the session is a TTY.
	return termenv.NewOutput(o, termenv.WithUnsafe(), termenv.WithEnvironment(e))
}

// Handle SSH requests.
func handler(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		// Get client's output.
		clientOutput := outputFromSession(sess)

		pty, _, active := sess.Pty()
		if !active {
			next(sess)
			return
		}
		width := pty.Window.Width

		// Initialize new renderer for the client.
		renderer := lipgloss.NewRenderer(sess)
		renderer.Output = clientOutput

		// Initialize new styles against the renderer.
		styles := makeStyles()

		str := strings.Builder{}

		fmt.Fprintf(&str, "\n\n%s %s %s %s %s",
			styles.bold,
			styles.faint,
			styles.italic,
			styles.underline,
			styles.strikethrough,
		)

		fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s",
			styles.red,
			styles.green,
			styles.yellow,
			styles.blue,
			styles.magenta,
			styles.cyan,
			styles.gray,
		)

		fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s\n\n",
			styles.red,
			styles.green,
			styles.yellow,
			styles.blue,
			styles.magenta,
			styles.cyan,
			styles.gray,
		)

		fmt.Fprintf(&str, "%s %t %s\n\n",
			styles.bold.Copy().UnsetString().Render("Has dark background?"),
			renderer.HasDarkBackground(),
			"ABOBA", // TODO: print renderer.Output.BackgroundColor().Hex
		)

		block := renderer.Place(width,
			10 /*lipgloss.Height(str.String())*/, lipgloss.Center, lipgloss.Center, str.String(),
			lipgloss.WithWhitespaceChars("/"),
			lipgloss.WithWhitespaceForeground(lipgloss.FgAdaptiveColor("250", "236")),
		)

		// Render to client.
		wish.WriteString(sess, block)

		next(sess)
	}
}

func Main(context.Context) error {
	port := 3456
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf(":%d", port)),
		wish.WithHostKeyPath("ssh_example"),
		wish.WithMiddleware(handler, lm.Middleware()),
	)
	if err != nil {
		return err
	}

	log.Printf("SSH server listening on port %d", port)
	log.Printf("To connect from your local machine run: ssh localhost -p %d", port)
	return s.ListenAndServe()
}
