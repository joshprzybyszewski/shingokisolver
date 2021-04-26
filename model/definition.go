package model

import "fmt"

type Definition struct {
	Description string
	Nodes       []NodeLocation

	NumEdges   int
	Difficulty Difficulty
}

func (d Definition) String() string {
	return fmt.Sprintf("%dx%d %s (%s)", d.NumEdges, d.NumEdges, d.Difficulty, d.Description)
}
