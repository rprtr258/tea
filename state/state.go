// inspired by https://github.com/zlumer/tesm
package state

type State[Event, Cmd any] interface {
	Handle(Event) (State[Event, Cmd], []Cmd)
}

type StateFunc[Event, Cmd any] func(Event) (State[Event, Cmd], []Cmd)

func (f StateFunc[Event, Cmd]) Handle(ev Event) (State[Event, Cmd], []Cmd) {
	return f(ev)
}

func call[T any](f func(T), t T) {
	if f == nil {
		return
	}
	f(t)
}

type h2[A, B any] struct {
	A func(A)
	B func(B)
}

type Sum2[A, B any] func(h2[A, B])

type new2[A, B any] struct {
	A func(A) Sum2[A, B]
	B func(B) Sum2[A, B]
}

func New2[A, B any]() new2[A, B] {
	return new2[A, B]{
		A: func(a A) Sum2[A, B] { return func(h h2[A, B]) { call(h.A, a) } },
		B: func(b B) Sum2[A, B] { return func(h h2[A, B]) { call(h.B, b) } },
	}
}

type h3[A, B, C any] struct {
	A func(A)
	B func(B)
	C func(C)
}

type Sum3[A, B, C any] func(h3[A, B, C])

type new3[A, B, C any] struct {
	A func(A) Sum3[A, B, C]
	B func(B) Sum3[A, B, C]
	C func(C) Sum3[A, B, C]
}

func New3[A, B, C any]() new3[A, B, C] {
	return new3[A, B, C]{
		A: func(a A) Sum3[A, B, C] { return func(h h3[A, B, C]) { call(h.A, a) } },
		B: func(b B) Sum3[A, B, C] { return func(h h3[A, B, C]) { call(h.B, b) } },
		C: func(c C) Sum3[A, B, C] { return func(h h3[A, B, C]) { call(h.C, c) } },
	}
}

type h4[A, B, C, D any] struct {
	A func(A)
	B func(B)
	C func(C)
	D func(D)
}

type Sum4[A, B, C, D any] func(h4[A, B, C, D])

type new4[A, B, C, D any] struct {
	A func(A) Sum4[A, B, C, D]
	B func(B) Sum4[A, B, C, D]
	C func(C) Sum4[A, B, C, D]
	D func(D) Sum4[A, B, C, D]
}

func New4[A, B, C, D any]() new4[A, B, C, D] {
	return new4[A, B, C, D]{
		A: func(a A) Sum4[A, B, C, D] { return func(h h4[A, B, C, D]) { call(h.A, a) } },
		B: func(b B) Sum4[A, B, C, D] { return func(h h4[A, B, C, D]) { call(h.B, b) } },
		C: func(c C) Sum4[A, B, C, D] { return func(h h4[A, B, C, D]) { call(h.C, c) } },
		D: func(d D) Sum4[A, B, C, D] { return func(h h4[A, B, C, D]) { call(h.D, d) } },
	}
}

// type Hook[State, Event, Cmd any] struct {
// 	state  State
// 	update func(State, Event) (State, []Cmd)
// 	cmds   []Cmd
// }

// func NewHook[State, Event, Cmd any](
// 	state State,
// 	update func(State, Event) (State, []Cmd),
// ) Hook[State, Event, Cmd] {
// 	return Hook[State, Event, Cmd]{state, update, nil}
// }

// func (h *Hook[S, E, C]) Send(events ...E) {
// 	for _, event := range events {
// 		state, cmds := h.update(h.state, event)
// 		h.state = state
// 		h.cmds = append(h.cmds, cmds...)
// 	}
// }

// func (h *Hook[S, E, C]) State() S {
// 	return h.state
// }
