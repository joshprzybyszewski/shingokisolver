package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func flip(
	puzz puzzle.Puzzle,
) (puzzle.Puzzle, bool) {

	ep, ok := puzz.GetUnknownEdge()
	if !ok {
		switch puzz.GetState(model.InvalidNodeCoord) {
		case model.Complete:
			return puzz, true
		default:
			return puzzle.Puzzle{}, false
		}
	}

	switch puzzWithEdge := puzz.DeepCopy(); puzzWithEdge.AddEdges(ep) {
	case model.Complete, model.Incomplete:
		if puzzWithEdge.GetState(ep.NodeCoord) == model.Complete {
			return puzzWithEdge, true
		}

		res, isComplete := flip(puzzWithEdge)
		if isComplete {
			return res, true
		}
	}

	switch puzzWithoutEdge := puzz.DeepCopy(); puzzWithoutEdge.AvoidEdge(ep) {
	case model.Complete, model.Incomplete:
		if puzzWithoutEdge.GetState(ep.NodeCoord) == model.Complete {
			return puzzWithoutEdge, true
		}
		res, isComplete := flip(puzzWithoutEdge)
		if isComplete {
			return res, true
		}
	}

	return puzzle.Puzzle{}, false
}
