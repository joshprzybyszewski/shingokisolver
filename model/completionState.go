package model

type State uint8

const (
	Complete      State = 1
	NodesComplete State = 2
	Incomplete    State = 3
	Violation     State = 4
	Unexpected    State = 5
	Duplicate     State = 6
	Ok            State = 7
)

func (s State) String() string {
	switch s {
	case Complete:
		return `Complete`
	case Incomplete:
		return `Incomplete`
	case Violation:
		return `Violation`
	case Unexpected:
		return `Unexpected`
	case Duplicate:
		return `Duplicate`
	case Ok:
		return `Ok`
	default:
		return `unknown State`
	}
}
