package solvers

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
) *puzzle.Puzzle {

	looseEnd, state := puzz.GetLooseEnd()
	// TODO remove Sprintf
	printAllTargetsHit(
		fmt.Sprintf(`connect(%+v, %s)`, looseEnd, state),
		puzz,
		d.iterations(),
	)

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
