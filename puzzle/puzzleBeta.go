package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

func (p Puzzle) Alpha() *state.TriEdges {
	return p.edges
}

func (p Puzzle) Beta() *logic.RuleSet {
	return p.rules
}
