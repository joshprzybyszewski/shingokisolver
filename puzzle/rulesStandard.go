package puzzle

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func getStandardNodeRules(
	ep model.EdgePair,
	otherStartEdges []model.EdgePair,
	otherEndEdges []model.EdgePair,
) []func(model.GetEdger) model.EdgeState {
	if len(otherStartEdges) != 3 || len(otherEndEdges) != 3 {
		panic(fmt.Sprintf(`unexpected input: %+v, %+v`, otherStartEdges, otherEndEdges))
	}

	return []func(model.GetEdger) model.EdgeState{
		getStandardEvalFor(otherStartEdges),
		getStandardEvalFor(otherEndEdges),
	}
}

func getStandardEvalFor(otherInputs []model.EdgePair) func(ge model.GetEdger) model.EdgeState {
	return func(ge model.GetEdger) model.EdgeState {
		numNonExisting := 0
		numExisting := 0

		for _, otherEP := range otherInputs {
			switch s := ge.GetEdge(otherEP); s {
			case model.EdgeExists:
				numExisting++
			case model.EdgeAvoided, model.EdgeOutOfBounds:
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
}
