package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
	looseEnds []model.NodeCoord,
) *puzzle.Puzzle {

	printAllTargetsHit(
		`connect`,
		puzz,
		looseEnds,
		d.iterations(),
	)

	looseEndsDeduped := getLooseEndsWithoutDuplicates(looseEnds)
	if len(looseEndsDeduped) == 0 {
		isIncomplete, err := puzz.IsIncomplete(model.NodeCoord{})
		if err != nil || isIncomplete {
			return nil
		}
		// otherwise, if we have no loose ends, then we can't do anything!
		return puzz
	}

	dfs := newDFSSolverForPartialSolution()
	defer func() {
		d.numProcessed += dfs.iterations()
	}()

	start := looseEndsDeduped[0]

	p, morePartials, sol := dfs.solveForGoals(
		puzz.DeepCopy(),
		start,
		looseEndsDeduped[1:],
	)

	switch sol {
	case solvedPuzzle:
		return p
	}

	// we only need to look at the first loose end in the
	// puzzle, so we return the following.
	return d.iterateMorePartials(start, looseEndsDeduped[1:], morePartials)
}

func (d *targetSolver) iterateMorePartials(
	start model.NodeCoord,
	otherLooseEnds []model.NodeCoord,
	morePartials map[model.NodeCoord][]*puzzle.Puzzle,
) *puzzle.Puzzle {
	for hitGoal, slice := range morePartials {
		for _, nextPuzzle := range slice {

			// TODO remove this print
			printAllTargetsHit(
				`iterateMorePartials`,
				nextPuzzle,
				copyAndRemove(otherLooseEnds, hitGoal),
				d.iterations(),
			)

			puzz := d.connect(
				nextPuzzle,
				copyAndRemove(otherLooseEnds, hitGoal),
			)
			if puzz != nil {
				return puzz
			}
		}
	}
	return nil
}

func copyAndRemove(orig []model.NodeCoord, exclude model.NodeCoord) []model.NodeCoord {
	cpy := make([]model.NodeCoord, 0, len(orig)-1)
	for _, le := range orig {
		if le != exclude {
			cpy = append(cpy, le)
		}
	}
	return cpy
}
