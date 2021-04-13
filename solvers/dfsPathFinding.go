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

	p, state := d.dfsOutFrom(
		input.DeepCopy(),
		start,
	)
	return p, state
}

func (d *targetSolver) dfsOutFrom(
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
) (*puzzle.Puzzle, model.State) {
	if puzz == nil {
		return nil, model.Unexpected
	}
	printAllTargetsHit(`dfsOutFrom`, puzz, d.iterations())

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
		switch nextPuzz.AddEdge(fromCoord, nextHeading) {
		case model.Violation, model.Unexpected, model.Duplicate:
			continue
		}
		nextCoord := fromCoord.Translate(nextHeading)

		if oe, ok := nextPuzz.GetOutgoingEdgesFrom(nextCoord); ok && oe.GetNumOutgoingDirections() == 2 {
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
