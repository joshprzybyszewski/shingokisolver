package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func firstFlip(
	puzz puzzle.Puzzle,
) (puzzle.Puzzle, bool) {
	puzzle.SetNodesComplete(&puzz)

	printPuzzleUpdate(`firstFlip`, puzz, model.InvalidTarget)

	ep, ok := puzz.GetUnknownEdge()
	if !ok {
		switch puzz.GetState() {
		case model.Complete:
			return puzz, true
		default:
			return puzzle.Puzzle{}, false
		}
	}

	return flip(puzz, ep)
}

func flip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) (puzzle.Puzzle, bool) {

	printPuzzleUpdate(`flip`, puzz, model.InvalidTarget)

	var nextUnknown model.EdgePair

	puzzWithEdge, state := puzzle.AddEdge(
		puzz,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		nextUnknown, state = puzzWithEdge.GetStateOfLoop()
		if state == model.Complete {
			return puzzWithEdge, true
		}

		res, isComplete := flip(puzzWithEdge, nextUnknown)
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
		nextUnknown, state = puzzWithoutEdge.GetStateOfLoop()
		if state == model.Complete {
			return puzzWithoutEdge, true
		}

		res, isComplete := flip(puzzWithoutEdge, nextUnknown)
		if isComplete {
			return res, true
		}
	}

	return puzzle.Puzzle{}, false
}
