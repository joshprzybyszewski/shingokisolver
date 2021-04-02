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
	return `TODO`
}
