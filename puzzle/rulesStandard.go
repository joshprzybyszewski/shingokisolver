package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func getStandardNodeRules(
	ep edgePair,
	otherStartEdges, otherEndEdges []edgePair,
) func(ge getEdger) model.EdgeState {

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
		if numAvoided > 2 {
			return model.EdgeErrored
		}

		shouldAvoid := numExisting == 2
		shouldExist := numAvoided+numOutOfBounds == 2 && numExisting == 1

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

		return model.EdgeUnknown
	}
}
