package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
) *puzzle.Puzzle {

	looseEnd, state := puzz.GetLooseEnd()

	switch state {
	case model.Complete:
		return puzz
	case model.NodesComplete:
		if puzz.GetState() == model.Complete {
			return puzz
		}
		return nil
	case model.Incomplete:
		// keep on goin
	default:
		printAllTargetsHit(`GetLooseEnd returned a state`, puzz, d.iterations())
		return nil
	}

	p, sol := d.solveFromLooseEnd(
		puzz.DeepCopy(),
		looseEnd,
	)

	switch sol {
	case model.Complete:
		return p
	}

	return nil
}
