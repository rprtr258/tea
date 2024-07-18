package state

type Status int

const (
	StatusSuccess Status = iota
	StatusFail
	StatusNormal
)

// type Transitions [][]State // (source)State, Event, (target)State
type Transitions[State, Event comparable] map[State]map[Event]State

func FSM[State, Event comparable](
	s0 State,
	tr Transitions[State, Event],
	events []Event,
) (State, Status) {
	s := s0
	for _, e := range events {
		to, ok := tr[s][e]
		if !ok {
			return *new(State), StatusFail
		}
		s = to
	}
	return s, StatusSuccess
}
