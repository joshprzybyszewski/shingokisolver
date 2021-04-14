package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) solveFromLooseEnd(
	input *puzzle.Puzzle,
	start model.NodeCoord,
) (*puzzle.Puzzle, model.State) {

	if input == nil {
		return nil, model.Unexpected
	}
	inputCopy := input.DeepCopy()

	switch state := d.sendOutDFSPath(
		inputCopy,
		start,
		model.HeadNowhere,
	); state {
	case model.Incomplete:
		retPuzz := d.connect(inputCopy)
		if retPuzz != nil {
			return retPuzz, model.Complete
		}
		return nil, state
	default:
		return nil, state
	}
}

func (d *targetSolver) sendOutDFSPath(
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
	avoid model.Cardinal,
	curPath ...puzzle.EdgePair,
) model.State {

	for _, nextHeading := range model.AllCardinals {
		if nextHeading == avoid {
			continue
		}
		ep := puzzle.NewEdgePair(fromCoord, nextHeading)
		if puzz.GetEdgeState(ep) != model.EdgeUnknown {
			continue
		} else if ep.IsIn(curPath...) {
			continue
		}

		switch state := d.sendOutDFSPath(
			puzz.DeepCopy(),
			fromCoord.Translate(nextHeading),
			nextHeading.Opposite(),
			append(curPath, ep)...,
		); state {
		case model.Incomplete:
			// if the puzzle isn't complete, allow it to continue
		default:
			return state
		}
	}

	// add all edges
	return puzz.AddEdges(curPath...)
}
