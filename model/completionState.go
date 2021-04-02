package model

type State uint8

const (
	Complete   State = 1
	Incomplete State = 2
	Violation  State = 3
	Unexpected State = 4
	Duplicate  State = 5
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
	default:
		return `unknown State`
	}
}
