package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func getStandardNodeRules(
	ep edgePair,
	otherStartEdges, otherEndEdges []edgePair,
) func(ge getEdger) model.EdgeState {

	if len(otherStartEdges) != 3 || len(otherEndEdges) != 3 {
		panic(`unexpected input`)
	}

	return func(ge getEdger) model.EdgeState {
		numAvoided := 0
		numOutOfBounds := 0
		numExisting := 0

		for _, otherEP := range otherStartEdges {
			switch s := ge.GetEdge(otherEP); s {
			case model.EdgeExists:
				numExisting++
			case model.EdgeAvoided:
				numAvoided++
			case model.EdgeOutOfBounds:
				numOutOfBounds++
			}
		}

		if numExisting > 2 {
			return model.EdgeErrored
		}

		shouldAvoid := numExisting == 2 || numAvoided+numOutOfBounds == 3
		shouldExist := numAvoided+numOutOfBounds == 2 && numExisting == 1

		if shouldAvoid && shouldExist {
			return model.EdgeErrored
		}

		numAvoided = 0
		numOutOfBounds = 0
		numExisting = 0

		for _, otherEP := range otherEndEdges {
			switch s := ge.GetEdge(otherEP); s {
			case model.EdgeExists:
				numExisting++
			case model.EdgeAvoided:
				numAvoided++
			case model.EdgeOutOfBounds:
				numOutOfBounds++
			}
		}

		if numExisting > 2 {
			return model.EdgeErrored
		}
		if numAvoided > 2 {
			return model.EdgeErrored
		}

		if numExisting == 2 {
			if shouldExist {
				return model.EdgeErrored
			}
			return model.EdgeAvoided
		}

		if numAvoided+numOutOfBounds == 2 && numExisting == 1 {
			if shouldAvoid {
				return model.EdgeErrored
			}
			return model.EdgeExists
		}

		if shouldExist {
			return model.EdgeExists
		}

		if shouldAvoid {
			return model.EdgeAvoided
		}

		return model.EdgeUnknown
	}
}
