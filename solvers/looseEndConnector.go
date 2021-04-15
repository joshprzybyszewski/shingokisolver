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
	case model.NodesComplete:
		if puzz.GetState() == model.Complete {
			return puzz
		}
		return nil
	case model.Incomplete:
		return d.solveFromLooseEnd(
			puzz.DeepCopy(),
			looseEnd,
		)
	default:
		printAllTargetsHit(`GetLooseEnd returned a state`, puzz, d.iterations())
		return nil
	}
}
