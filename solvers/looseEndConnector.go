package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
) *puzzle.Puzzle {

	looseEnds := puzz.LooseEnds()
	printAllTargetsHit(`connect`, puzz, d.iterations())

	if len(looseEnds) == 0 {
		switch puzz.GetState() {
		case model.Complete:
			return puzz
		default:
			printAllTargetsHit(`had no loose ends, but not complete`, puzz, d.iterations())
			return nil
		}
	}

	start := looseEnds[0]

	p, morePartials, sol := d.solveForGoals(
		puzz.DeepCopy(),
		start,
		looseEnds[1:],
	)

	switch sol {
	case model.Complete:
		return p
	}

	// we only need to look at the first loose end in the
	// puzzle, so we return the following.
	return d.iterateMorePartials(morePartials)
}

func (d *targetSolver) iterateMorePartials(
	morePartials map[model.NodeCoord][]*puzzle.Puzzle,
) *puzzle.Puzzle {

	for _, slice := range morePartials {
		for _, nextPuzzle := range slice {
			puzz := d.connect(
				nextPuzzle,
			)
			if puzz != nil {
				return puzz
			}
		}
	}

	return nil
}
