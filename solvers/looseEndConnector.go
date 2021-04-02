package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) connect(
	puzz *puzzle.Puzzle,
	looseEnds []model.NodeCoord,
) *puzzle.Puzzle {

	printAllTargetsHit(`connect`, puzz, looseEnds, d.iterations())

	looseEndsDeduped := dedupeLooseEnds(looseEnds)
	if len(looseEndsDeduped) == 0 {
		switch puzz.GetState() {
		case model.Complete:
			return puzz
		default:
			return nil
		}
	}

	start := looseEndsDeduped[0]

	p, morePartials, sol := d.solveForGoals(
		puzz.DeepCopy(),
		start,
		looseEndsDeduped[1:],
	)

	switch sol {
	case model.Complete:
		return p
	}

	// we only need to look at the first loose end in the
	// puzzle, so we return the following.
	return d.iterateMorePartials(start, looseEndsDeduped[1:], morePartials)
}

// eliminates loose ends that don't actually exist
// leaves the looseEnds slice in the order that it had previously
func dedupeLooseEnds(looseEnds []model.NodeCoord) []model.NodeCoord {

	numExisting := make(map[model.NodeCoord]int, len(looseEnds))
	for _, le := range looseEnds {
		numExisting[le] += 1
	}

	looseEndsDeduped := make([]model.NodeCoord, 0, len(looseEnds))
	for _, le := range looseEnds {
		if existing := numExisting[le]; existing%2 == 0 {
			looseEndsDeduped = append(looseEndsDeduped, le)
		}
	}

	return looseEndsDeduped
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
