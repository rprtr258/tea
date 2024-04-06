package simple

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rprtr258/assert"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/teatest"
)

func TestApp(t *testing.T) {
	m := model(10)
	tm := teatest.NewTestModelFixture(t, &m,
		teatest.WithInitialTermSize(70, 30),
	)

	time.Sleep(time.Second + time.Millisecond*200)
	tm.Type("I'm typing things, but it'll be ignored by my program")
	tm.Send("ignored msg")
	tm.Send(tea.MsgKey{
		Type: tea.KeyEnter,
	})

	assert.NoError(t, tm.Quit())

	out := readBts(t, tm.FinalOutput(t))
	assert.True(t, regexp.MustCompile(`This program will exit in \d+ seconds`).Match(out))
	// TODO: get back
	// teatest.RequireEqualOutput(t, out)

	assert.Equal(t, model(9), *tm.FinalModel(t))
}

func TestAppInteractive(t *testing.T) {
	m := model(10)
	tm := teatest.NewTestModelFixture(t, &m,
		teatest.WithInitialTermSize(70, 30),
	)

	time.Sleep(time.Second + 200*time.Millisecond)
	tm.Send("ignored msg")

	assert.True(t, strings.Contains(string(readBts(t, tm.Output())), "This program will exit in 9 seconds"))

	teatest.WaitFor(t, tm.Output(), func(out []byte) bool {
		return bytes.Contains(out, []byte("This program will exit in 7 seconds"))
	}, teatest.WithDuration(3*time.Second))

	tm.Send(tea.MsgKey{
		Type: tea.KeyEnter,
	})

	assert.Equal(t, model(7), *tm.FinalModel(t))
}

func readBts(t *testing.T, r io.Reader) []byte {
	t.Helper()
	bts, err := io.ReadAll(r)
	assert.NoError(t, err)
	return bts
}
