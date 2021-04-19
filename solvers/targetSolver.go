package solvers

import (
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type targetSolver struct {
	puzzle puzzle.Puzzle
}

func newTargetSolver(
	size int,
	nl []model.NodeLocation,
) Solver {
	if len(nl) == 0 {
		return nil
	}

	return &targetSolver{
		puzzle: puzzle.NewPuzzle(size, nl),
	}
}

func (d *targetSolver) Solve() (SolvedResults, error) {
	t0 := time.Now()

	puzz, isSolved := d.solve()
	if !isSolved {
		return SolvedResults{
			Duration: time.Since(t0),
		}, fmt.Errorf("puzzle unsolvable: %s", d.puzzle.String())
	}

	return SolvedResults{
		Puzzle:   puzz,
		Duration: time.Since(t0),
	}, nil
}

func (d *targetSolver) solve() (puzzle.Puzzle, bool) {

	puzz, s := d.puzzle.DeepCopy().ClaimGimmes()
	switch s {
	case model.Incomplete, model.Complete:
		// printPuzzleUpdate(`ClaimGimmes`, 0, puzz, nil)
	default:
		return puzzle.Puzzle{}, false
	}

	target, state := d.puzzle.GetFirstTarget()
	switch state {
	case model.Incomplete, model.Complete:
		// printPuzzleUpdate(`GetFirstTarget`, 0, puzz, target)
	default:
		return puzzle.Puzzle{}, false
	}

	return d.getSolutionFromDepths(
		0,
		puzz.DeepCopy(),
		target,
	)
}

func (d *targetSolver) getSolutionFromDepths(
	depth int,
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

			return d.flip(
				puzz.DeepCopy(),
			)
		}

		return d.getSolutionFromDepths(
			depth+1,
			puzz.DeepCopy(),
			nextTarget,
		)
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	for _, option := range targeting.Options {
		// then, once we find a completion path, add it to the returned slice
		p, isComplete := d.buildAllTwoArmsForTraversal(
			depth,
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

func (d *targetSolver) buildAllTwoArmsForTraversal(
	depth int,
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) (puzzle.Puzzle, bool) {

	switch puzz.AddEdges(getTwoArmsEdges(
		curTarget.Node.Coord(),
		ta,
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

		return d.flip(
			puzz.DeepCopy(),
		)
	}

	return d.getSolutionFromDepths(
		depth+1,
		puzz.DeepCopy(),
		nextTarget,
	)
}

func getTwoArmsEdges(
	start model.NodeCoord,
	ta model.TwoArms,
) []model.EdgePair {

	allEdges := make([]model.EdgePair, 0, ta.One.Len+ta.Two.Len)

	arm1End := start
	for i := int8(0); i < ta.One.Len; i++ {
		allEdges = append(allEdges, model.NewEdgePair(arm1End, ta.One.Heading))
		arm1End = arm1End.Translate(ta.One.Heading)
	}

	arm2End := start
	for i := int8(0); i < ta.Two.Len; i++ {
		allEdges = append(allEdges, model.NewEdgePair(arm2End, ta.Two.Heading))
		arm2End = arm2End.Translate(ta.Two.Heading)
	}

	return allEdges
}
