package teatest

import (
	"fmt"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/rprtr258/assert"
)

func TestRequireEqualOutputUpdate(t *testing.T) {
	enableUpdate(t)
	RequireEqualOutput(t, []byte("test"))
}

func TestWaitForErrorReader(t *testing.T) {
	err := doWaitFor(iotest.ErrReader(fmt.Errorf("fake")), func([]byte) bool {
		return true
	}, WithDuration(time.Millisecond), WithCheckInterval(10*time.Microsecond))
	assert.Equal(t, err.Error(), "WaitFor: fake")
}

func TestWaitForTimeout(t *testing.T) {
	err := doWaitFor(strings.NewReader("nope"), func([]byte) bool {
		return false
	}, WithDuration(time.Millisecond), WithCheckInterval(10*time.Microsecond))
	assert.Equal(t, err.Error(), "WaitFor: condition not met after 1ms")
}

func enableUpdate(tb testing.TB) {
	tb.Helper()
	previous := update
	*update = true
	tb.Cleanup(func() {
		update = previous
	})
}
