package solvers

import (
	"fmt"

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

	p, state := d.dfsOutFrom(
		input.DeepCopy(),
		start,
	)
	return p, state
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

func (d *targetSolver) dfsOutFrom(
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
) (*puzzle.Puzzle, model.State) {
	if puzz == nil {
		return nil, model.Unexpected
	}
	// TODO remove Sprintf
	printAllTargetsHit(
		fmt.Sprintf(`dfsOutFrom(%+v)`, fromCoord),
		puzz,
		d.iterations(),
	)

	switch s := puzz.GetState(); s {
	case model.Complete:
		return puzz, s
	case model.Incomplete:
		// continue in the func
	default:
		return nil, s
	}

	for _, nextHeading := range model.AllCardinals {
		nextPuzz := puzz.DeepCopy()

		d.numProcessed++
		shouldContinue := true
		switch nextPuzz.AddEdge(fromCoord, nextHeading) {
		case model.Incomplete, model.Complete:
			shouldContinue = false
		}
		if shouldContinue {
			continue
		}

		nextCoord := fromCoord.Translate(nextHeading)

		if nextPuzz.HasTwoOutgoingEdges(nextCoord) {
			// we connected to an existing path.
			// iterate down from another loose end
			retPuzz := d.connect(nextPuzz)
			if retPuzz != nil {
				return retPuzz, model.Complete
			}
		} else {
			retPuzz, s := d.dfsOutFrom(
				nextPuzz,
				nextCoord,
			)
			switch s {
			case model.Complete:
				return retPuzz, s
			}
		}
	}

	return nil, model.Incomplete
}
