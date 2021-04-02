package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
) *puzzle.Puzzle {

	looseEnd, ok := puzz.GetLooseEnd()
	printAllTargetsHit(`connect`, puzz, d.iterations())

	if !ok {
		switch puzz.GetState() {
		case model.Complete:
			return puzz
		default:
			printAllTargetsHit(`had no loose ends, but not complete`, puzz, d.iterations())
			return nil
		}
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
