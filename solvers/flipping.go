package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func flip(
	puzz puzzle.Puzzle,
) (puzzle.Puzzle, bool) {

	printPuzzleUpdate(`flip`, puzz, model.InvalidTarget)

	ep, ok := puzz.GetUnknownEdge()
	if !ok {
		switch puzz.GetState(model.InvalidNodeCoord) {
		case model.Complete:
			return puzz, true
		default:
			return puzzle.Puzzle{}, false
		}
	}

	puzzWithEdge, state := puzzle.AddEdge(
		puzz,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		if puzzWithEdge.GetState(ep.NodeCoord) == model.Complete {
			return puzzWithEdge, true
		}

		res, isComplete := flip(puzzWithEdge)
		if isComplete {
			return res, true
		}
	}

	puzzWithoutEdge, state := puzzle.AvoidEdge(
		puzz,
		ep,
	)
	switch state {
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
