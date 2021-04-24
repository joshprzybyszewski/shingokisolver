package solvers

import (
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

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

	puzz, target, state := claimGimmes(puzz)
	switch state {
	case model.Complete:
		return puzz, true
	case model.NodesComplete:
		// ok, the nodes are complete, but we're not solved.
		// don't "aim at target" below, but start flipping edges
		// now.
		return firstFlip(
			puzz,
		)
	case model.Incomplete:
		// continue on
	default:
		printPuzzleUpdate(`claimGimmes issue!`, puzz, target)
		return puzzle.Puzzle{}, false
	}

	// Now we're going to start descending through the nodes, aiming at our next
	// target.
	return solveAimingAtTarget(
		puzz,
		target,
	)
}

func claimGimmes(
	puzz puzzle.Puzzle,
) (puzzle.Puzzle, model.Target, model.State) {

	// claim all of the gimmes we can
	puzz, state := puzzle.ClaimGimmes(puzz)
	switch state {
	case model.Incomplete, model.Complete:
		printPuzzleUpdate(`ClaimGimmes`, puzz, model.Target{})
	default:
		return puzzle.Puzzle{}, model.Target{}, state
	}

	// Get the first node to target in the puzzle
	target, state := puzz.GetFirstTarget()
	switch state {
	case model.Complete:
		// already solved!
		return puzz, model.Target{}, state
	case model.Incomplete:
		printPuzzleUpdate(`GetFirstTarget`, puzz, target)
	default:
		return puzzle.Puzzle{}, model.Target{}, state
	}

	// it's likely that a lot of the nodes will be able to be claimed right
	// away (that is, they only have one option to be satisfied). Think a w2
	// that's on the side, or a b2 that's in a corner. Let's make sure that
	// they all have the edges they need, and iterate through all of them
	var ok bool
	for len(target.Options) == 1 {
		puzz, ok = addTwoArms(puzz, target.Node.Coord(), target.Options[0])
		if !ok {
			printPuzzleUpdate(`addTwoArms issue!`, puzz, target)
			return puzzle.Puzzle{}, model.Target{}, model.Violation
		}

		target, state = puzz.GetNextTarget(target)
		switch state {
		case model.Complete:
			// already solved!
			return puzz, model.Target{}, state
		case model.NodesComplete:
			return puzz, model.Target{}, state
		case model.Violation:
			printPuzzleUpdate(`GetNextTarget issue!`, puzz, target)
			return puzzle.Puzzle{}, model.Target{}, state
		}
	}

	return puzz, target, model.Incomplete
}

func solveAimingAtTarget(
	puzz puzzle.Puzzle,
	targeting model.Target,
) (puzzle.Puzzle, bool) {

	printPuzzleUpdate(`solveAimingAtTarget`, puzz, targeting)

	// Check to see if this node has already been completed.
	switch puzz.GetNodeState(targeting.Node.Coord()) {
	case model.Violation:
		printPuzzleUpdate(`solveAimingAtTarget GetNodeState issue!`, puzz, targeting)
		return puzzle.Puzzle{}, false

	case model.Complete:
		return descendToNextTarget(
			puzz,
			targeting,
		)
	}

	// TODO concurrency!
	// for each of the TwoArm options, we're going to try setting the edges
	// and then descending further into our targets
	for _, option := range targeting.Options {
		// then, once we find a completion path, add it to the returned slice
		p, isComplete := buildTwoArmsToDescend(
			puzz,
			targeting,
			option,
		)
		if isComplete {
			return p, true
		}
	}

	return puzzle.Puzzle{}, false
}

func buildTwoArmsToDescend(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) (puzzle.Puzzle, bool) {

	puzz, ok := addTwoArms(puzz, curTarget.Node.Coord(), ta)
	if !ok {
		return puzzle.Puzzle{}, false
	}

	return descendToNextTarget(puzz, curTarget)
}

func descendToNextTarget(
	puzz puzzle.Puzzle,
	curTarget model.Target,
) (puzzle.Puzzle, bool) {

	printPuzzleUpdate(`descendToNextTarget`, puzz, curTarget)

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Violation, model.Unexpected:
		return puzzle.Puzzle{}, false
	case model.Complete:
		// It's solved!
		return puzz, true
	case model.NodesComplete:
		// At this point, we know that the nodes are all "complete".
		// This means that we need to transition to "flipping edge state"
		// until we can find a complete puzzle
		return firstFlip(
			puzz,
		)
	default:
		return solveAimingAtTarget(
			puzz,
			nextTarget,
		)
	}
}

func addTwoArms(
	inPuzz puzzle.Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) (puzzle.Puzzle, bool) {

	outPuzz, state := puzzle.AddTwoArms(inPuzz, start, ta)
	if state != model.Incomplete {
		return puzzle.Puzzle{}, false
	}
	return outPuzz, true
}
