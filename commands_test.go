package tea

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type msgString string

func TestEvery(t *testing.T) {
	expected := msgString("every ms")
	msg := Every(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}

func TestTick(t *testing.T) {
	expected := msgString("tick")
	msg := Tick(time.Millisecond, func(t time.Time) Msg {
		return expected
	})()
	assert.Equal(t, expected, msg)
}
