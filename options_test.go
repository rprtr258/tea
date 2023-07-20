package tea

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("output", func(t *testing.T) {
		var b bytes.Buffer
		p := NewProgram(nil, WithOutput(&b))
		assert.Nil(t, p.output.TTY())
	})

	t.Run("custom input", func(t *testing.T) {
		var b bytes.Buffer
		p := NewProgram(nil, WithInput(&b))
		assert.Equal(t, &b, p.input)
		assert.Equal(t, customInput, p.inputType)
	})

	t.Run("renderer", func(t *testing.T) {
		p := NewProgram(nil, WithoutRenderer())
		assert.IsType(t, (*nilRenderer)(nil), p.renderer)
	})

	t.Run("without signals", func(t *testing.T) {
		p := NewProgram(nil, WithoutSignals())
		assert.True(t, p.ignoreSignals)
	})

	t.Run("filter", func(t *testing.T) {
		p := NewProgram(nil, WithFilter(func(_ Model, msg Msg) Msg { return msg }))
		assert.NotNil(t, p.filter)
	})

	t.Run("input options", func(t *testing.T) {
		exercise := func(t *testing.T, opt ProgramOption, expect inputType) {
			p := NewProgram(nil, opt)
			assert.Equal(t, expect, p.inputType)
		}

		t.Run("tty input", func(t *testing.T) {
			exercise(t, WithInputTTY(), ttyInput)
		})

		t.Run("custom input", func(t *testing.T) {
			var b bytes.Buffer
			exercise(t, WithInput(&b), customInput)
		})
	})

	t.Run("startup options", func(t *testing.T) {
		exercise := func(t *testing.T, opt ProgramOption, expect startupOptions) {
			p := NewProgram(nil, opt)
			assert.True(t, p.startupOptions.has(expect))
		}

		t.Run("alt screen", func(t *testing.T) {
			exercise(t, WithAltScreen(), withAltScreen)
		})

		t.Run("ansi compression", func(t *testing.T) {
			exercise(t, WithANSICompressor(), withANSICompressor)
		})

		t.Run("without catch panics", func(t *testing.T) {
			exercise(t, WithoutCatchPanics(), withoutCatchPanics)
		})

		t.Run("without signal handler", func(t *testing.T) {
			exercise(t, WithoutSignalHandler(), withoutSignalHandler)
		})

		t.Run("mouse cell motion", func(t *testing.T) {
			p := NewProgram(nil, WithMouseAllMotion(), WithMouseCellMotion())
			assert.True(t, p.startupOptions.has(withMouseCellMotion))
			assert.False(t, p.startupOptions.has(withMouseAllMotion))
		})

		t.Run("mouse all motion", func(t *testing.T) {
			p := NewProgram(nil, WithMouseCellMotion(), WithMouseAllMotion())
			assert.True(t, p.startupOptions.has(withMouseAllMotion))
			assert.False(t, p.startupOptions.has(withMouseCellMotion))
		})
	})

	t.Run("multiple", func(t *testing.T) {
		p := NewProgram(nil, WithMouseAllMotion(), WithAltScreen(), WithInputTTY())
		for _, opt := range []startupOptions{withMouseAllMotion, withAltScreen} {
			assert.True(t, p.startupOptions.has(opt))
			assert.Equal(t, ttyInput, p.inputType)
		}
	})
}
