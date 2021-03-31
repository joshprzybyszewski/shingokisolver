package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type looseEndConnector struct {
	iterations int
}

func (lec *looseEndConnector) connect(
	partial *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(`connect`, partial, lec.iterations)

	numLooseEndsByNode := make(map[model.NodeCoord]int, len(partial.looseEnds))
	for _, g := range partial.looseEnds {
		numLooseEndsByNode[g] += 1
	}

	for i, start := range partial.looseEnds {
		if numLooseEndsByNode[start] > 1 {
			// This means that two nodes have a "loose end" that meets
			// up at the same point. Let's just skip trying to find this
			// "loose" end a buddy since it already has one.
			continue
		}

		dfs := newDFSSolverForPartialSolution()
		p, morePartials, sol := dfs.solveForGoals(
			partial.puzzle,
			start,
			partial.looseEnds[i+1:],
		)
		lec.iterations += dfs.iterations()

		switch sol {
		case solvedPuzzle:
			return p
		}

		// we only need to look at the first loose end in the
		// puzzle, so we return the following.
		return lec.iterateMorePartials(partial, start, morePartials)
	}

	return nil
}

func (lec *looseEndConnector) iterateMorePartials(
	partial *partialSolutionItem,
	start model.NodeCoord,
	morePartials map[model.NodeCoord][]*puzzle.Puzzle,
) *puzzle.Puzzle {
	looseEndsWithoutStart := remove(partial.looseEnds, start)

	for hitGoal, slice := range morePartials {
		for _, nextPuzzle := range slice {
			puzz := lec.connect(&partialSolutionItem{
				puzzle:    nextPuzzle,
				looseEnds: remove(looseEndsWithoutStart, hitGoal),
			})
			if puzz != nil {
				return puzz
			}
		}
	}
	return nil
}

func remove(orig []model.NodeCoord, exclude model.NodeCoord) []model.NodeCoord {
	cpy := make([]model.NodeCoord, 0, len(orig)-1)
	for _, le := range orig {
		if le != exclude {
			cpy = append(cpy, le)
		}
	}
	return cpy
}
