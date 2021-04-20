package solvers

import (
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func solveWithTargets(
	size int,
	nl []model.NodeLocation,
) (SolvedResults, error) {

	return solvePuzzleByTargets(
		puzzle.NewPuzzle(size, nl),
	)
}

func solvePuzzleByTargets(
	puzz puzzle.Puzzle,
) (sr SolvedResults, _ error) {
	defer func(t0 time.Time) {
		sr.Duration = time.Since(t0)
	}(time.Now())

	solution, isSolved := doSolve(puzz)
	if !isSolved {
		return SolvedResults{}, fmt.Errorf("puzzle unsolvable: %s", puzz.String())
	}

	return SolvedResults{
		Puzzle: solution,
	}, nil
}

func doSolve(
	puzz puzzle.Puzzle,
) (puzzle.Puzzle, bool) {

	puzz, s := puzz.DeepCopy().ClaimGimmes()
	switch s {
	case model.Incomplete, model.Complete:
		// printPuzzleUpdate(`ClaimGimmes`, 0, puzz, nil)
	default:
		return puzzle.Puzzle{}, false
	}

	target, state := puzz.GetFirstTarget()
	switch state {
	case model.Incomplete, model.Complete:
		// printPuzzleUpdate(`GetFirstTarget`, 0, puzz, target)
	default:
		return puzzle.Puzzle{}, false
	}

	return getSolutionFromDepths(
		puzz.DeepCopy(),
		target,
	)
}

func getSolutionFromDepths(
	puzz puzzle.Puzzle,
	targeting model.Target,
) (puzzle.Puzzle, bool) {

	nc := targeting.Node.Coord()

	// printPuzzleUpdate(`getSolutionFromDepths`, depth, puzz, targeting)

	switch puzz.GetNodeState(nc) {
	case model.Violation:
		return puzzle.Puzzle{}, false

	case model.Complete:
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends
		nextTarget, state := puzz.GetNextTarget(targeting)
		switch state {
		case model.Violation:
			return puzzle.Puzzle{}, false
		case model.NodesComplete:
			switch puzz.GetState(nc) {
			case model.Complete:
				return puzz, true
			}

			return flip(
				puzz.DeepCopy(),
			)
		}

		return getSolutionFromDepths(
			puzz.DeepCopy(),
			nextTarget,
		)
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	for _, option := range targeting.Options {
		// then, once we find a completion path, add it to the returned slice
		p, isComplete := buildAllTwoArmsForTraversal(
			puzz.DeepCopy(),
			targeting,
			option,
		)
		if isComplete {
			return p, true
		}
	}

	return puzzle.Puzzle{}, false
}

func buildAllTwoArmsForTraversal(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) (puzzle.Puzzle, bool) {

	switch puzz.AddEdges(ta.GetAllEdges(
		curTarget.Node.Coord(),
	)...) {
	case model.Duplicate, model.Incomplete, model.Complete:
	default:
		return puzzle.Puzzle{}, false
	}

	switch puzz.GetState(curTarget.Node.Coord()) {
	case model.Complete:
		return puzz, true
	case model.Incomplete:
		// continue
	default:
		return puzzle.Puzzle{}, false
	}

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Violation:
		return puzzle.Puzzle{}, false
	case model.NodesComplete:
		switch puzz.GetState(curTarget.Node.Coord()) {
		case model.Complete:
			return puzz, true
		}

		return flip(
			puzz.DeepCopy(),
		)
	}

	return getSolutionFromDepths(
		puzz.DeepCopy(),
		nextTarget,
	)
}
