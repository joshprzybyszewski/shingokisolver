package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

func (p Puzzle) Alpha() *state.TriEdges {
	return &p.edges
}

func (p Puzzle) Beta() *logic.RuleSet {
	return p.rules
}
