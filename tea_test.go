package tea

import (
	"bytes"
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rprtr258/assert"
)

type msgIncrement struct{}

type testModel struct {
	viewCalled atomic.Bool
	counter    atomic.Int32
}

func (m *testModel) Init(func(...Cmd)) {}

func (m *testModel) Update(msg Msg, f func(...Cmd)) {
	switch msg.(type) {
	case msgIncrement:
		m.counter.Add(1)
	case MsgKey:
		f(Quit)
	}
}

func (m *testModel) View(vb Viewbox) {
	m.viewCalled.Store(true)
	vb.WriteLine("success")
}

func TestTeaModel(t *testing.T) {
	in := bytes.NewBuffer([]byte("q"))

	var buf bytes.Buffer
	_, err := NewProgram(context.Background(), &testModel{}).WithInput(in).WithOutput(&buf).Run()
	assert.NoError(t, err)

	assert.True(t, len(buf.Bytes()) > 0)
}

func TestTeaQuit(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	p := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf)
	go func() {
		for {
			time.Sleep(time.Millisecond)
			if m.viewCalled.Load() {
				p.Quit()
				return
			}
		}
	}()

	_, err := p.Run()
	assert.NoError(t, err)
}

func TestTeaWithFilter(t *testing.T) {
	for preventCount := uint32(0); preventCount < 3; preventCount++ {
		preventCount := preventCount

		m := &testModel{}
		shutdowns := uint32(0)
		p := NewProgram(context.Background(), m).
			WithInput(&bytes.Buffer{}).
			WithOutput(&bytes.Buffer{}).
			WithFilter(func(_ *testModel, msg Msg) Msg {
				if _, ok := msg.(MsgQuit); !ok {
					return msg
				}
				if shutdowns < preventCount {
					atomic.AddUint32(&shutdowns, 1)
					return nil
				}
				return msg
			})

		go func() {
			for atomic.LoadUint32(&shutdowns) <= preventCount {
				time.Sleep(time.Millisecond)
				p.Quit()
			}
		}()

		_, err := p.Run()
		assert.NoError(t, err)
		assert.Equal(t, preventCount, shutdowns)
	}
}

func TestTeaKill(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	p := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf)
	go func() {
		for {
			time.Sleep(time.Millisecond)
			if m.viewCalled.Load() {
				p.Kill()
				return
			}
		}
	}()

	_, err := p.Run()
	assert.Equal(t, err, ErrProgramKilled)
}

func TestTeaContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	p := NewProgram(ctx, m).WithInput(&in).WithOutput(&buf)
	go func() {
		for {
			time.Sleep(time.Millisecond)
			if m.viewCalled.Load() {
				cancel()
				return
			}
		}
	}()

	_, err := p.Run()
	assert.Equal(t, err, ErrProgramKilled)
}

func TestMsgBatch(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	inc := func() Msg {
		return msgIncrement{}
	}

	m := &testModel{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewProgram(ctx, m).WithInput(&in).WithOutput(&buf)
	go func() {
		p.Send(inc())
		p.Send(inc())

		for {
			time.Sleep(time.Millisecond)
			if m.counter.Load() >= 2 {
				p.Quit()
				return
			}
		}
	}()

	_, err := p.Run()
	assert.NoError(t, err)
	assert.Equal(t, int32(2), m.counter.Load())
}

func TestTeaSend(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	p := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf)

	// sending before the program is started is a blocking operation
	go p.Send(Quit())

	_, err := p.Run()
	assert.NoError(t, err)

	// sending a message after program has quit is a no-op
	p.Send(Quit())
}

func TestTeaNoRun(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	assert.NotZero(t, NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf))
}
