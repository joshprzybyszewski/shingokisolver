package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

var _ evaluator = standardInput{}

type standardInput struct {
	otherInputs []model.EdgePair
}

func (si standardInput) evaluate(ge model.GetEdger) model.EdgeState {
	numNonExisting := 0
	numExisting := 0

	for _, otherEP := range si.otherInputs {
		if ge.IsEdge(otherEP) {
			numExisting++
		} else if ge.IsAvoided(otherEP) {
			numNonExisting++
		}
	}

	switch {
	case numExisting > 2:
		// branched!
		return model.EdgeErrored
	case numExisting == 2, numNonExisting == 3:
		// either the two inputs for this node have been defined,
		// or we know all three other inputs are not edges
		return model.EdgeAvoided
	case numExisting == 1 && numNonExisting == 2:
		// there exists one input to this node, and nobody else can be its pair
		// therefore we should be.
		return model.EdgeExists
	}

	return model.EdgeUnknown
}
