package tea

import (
	"bytes"
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type msgIncrement struct{}

type testModel struct {
	viewCalled atomic.Bool
	counter    atomic.Int32
}

func (m *testModel) Init() []Cmd {
	return nil
}

func (m *testModel) Update(msg Msg) []Cmd {
	switch msg.(type) {
	case msgIncrement:
		m.counter.Add(1)

	case MsgKey:
		return []Cmd{Quit}
	}

	return nil
}

func (m *testModel) View(r Renderer) {
	m.viewCalled.Store(true)
	r.Write("success\n")
}

func TestTeaModel(t *testing.T) {
	in := bytes.NewBuffer([]byte("q"))

	var buf bytes.Buffer
	_, err := NewProgram(context.Background(), &testModel{}).WithInput(in).WithOutput(&buf).Run()
	assert.NoError(t, err)

	assert.NotEmpty(t, buf.Bytes())
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

		assert.NoError(t, p.Start())
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
			if m.counter.Load() >= 2 { //nolint:lostcancel
				p.Quit()
				return
			}
		}
	}()

	_, err := p.Run()
	assert.NoError(t, err)
	assert.Equal(t, int32(2), m.counter.Load())
}

func TestMsgSequence(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	inc := func() Msg {
		return msgIncrement{}
	}

	m := &testModel{}
	p := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf)
	go p.Send(msgSequence{inc, inc, Quit})

	_, err := p.Run()
	assert.NoError(t, err)
	assert.Equal(t, int32(2), m.counter.Load())
}

func TestMsgSequenceWithMsgBatch(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	inc := func() Msg {
		return msgIncrement{}
	}

	m := &testModel{}
	p := NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf)
	go p.Send(msgSequence{inc, inc, inc, Quit})

	_, err := p.Run()
	assert.NoError(t, err)
	assert.Equal(t, int32(3), m.counter.Load())
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
	assert.NotNil(t, NewProgram(context.Background(), m).WithInput(&in).WithOutput(&buf))
}
