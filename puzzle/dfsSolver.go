package puzzle

import (
	"fmt"
)

type dfsSolverStep struct {
	puzzle *puzzle
	coord  nodeCoord
}

type dfsSolver struct {
	puzzle *puzzle

	goal *nodeCoord

	numProcessed int
}

func newDFSSolver(
	size int,
	nl []NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &dfsSolver{
		puzzle: newPuzzle(size, nl),
	}
}

func newDFSSolverForPartialSolution(
	partiallySolved *puzzle,
) *dfsSolver {
	if partiallySolved == nil {
		return nil
	}

	return &dfsSolver{
		puzzle: partiallySolved,
	}
}

func (d *dfsSolver) iterations() int {
	return d.numProcessed
}

func (d *dfsSolver) solve() (*puzzle, bool) {
	var bestCoord nodeCoord
	bestVal := int8(-1)
	for nc, n := range d.puzzle.nodes {
		if n.val > bestVal {
			bestCoord = nc
			bestVal = n.val
		}
	}
	return d.takeNextStepIntoDepth(&dfsSolverStep{
		puzzle: d.puzzle,
		coord:  bestCoord,
	})
}

func (d *dfsSolver) takeNextStepIntoDepth(
	q *dfsSolverStep,
) (*puzzle, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.row, q.coord.col, d.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.puzzle.IsIncomplete(q.coord); err != nil {
		return q.puzzle, false
	} else if !isIncomplete {
		return q.puzzle, true
	}

	for _, step := range []*dfsSolverStep{
		d.getNextStep(q.puzzle, headUp, q.coord),
		d.getNextStep(q.puzzle, headRight, q.coord),
		d.getNextStep(q.puzzle, headDown, q.coord),
		d.getNextStep(q.puzzle, headLeft, q.coord),
	} {
		g, isSolved := d.takeNextStepIntoDepth(step)
		if isSolved {
			return g, true
		}
	}

	return q.puzzle, false
}

func (d *dfsSolver) getNextStep(
	g *puzzle,
	move cardinal,
	nc nodeCoord,
) *dfsSolverStep {
	newCoord, newPuzzle, err := g.AddEdge(move, nc)
	if err != nil {
		return nil
	}

	if newPuzzle.isRangeInvalidWithBoundsCheck(
		newCoord.row-2,
		newCoord.row+2,
		newCoord.col-2,
		newCoord.col+2,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil
	}

	return &dfsSolverStep{
		puzzle: newPuzzle,
		coord:  newCoord,
	}
}

type dfsGoalSolution int

const (
	mayContinue  dfsGoalSolution = 0
	badState     dfsGoalSolution = 1
	solvedPuzzle dfsGoalSolution = 2
	foundGoal    dfsGoalSolution = 3
)

func (d *dfsSolver) solveForGoal(
	start, goal nodeCoord,
) (*puzzle, dfsGoalSolution) {
	d.goal = &goal
	// TODO if I could pass in a slice of goals, then return a map[goal][]*puzzle
	// then maybe that'd allow me to iterative connect points?

	return d.takeNextStepIntoDepthTowardsGoal(&dfsSolverStep{
		puzzle: d.puzzle,
		coord:  start,
	})
}

func (d *dfsSolver) takeNextStepIntoDepthTowardsGoal(
	q *dfsSolverStep,
) (*puzzle, dfsGoalSolution) {
	if q == nil {
		return nil, badState
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.row, q.coord.col, d.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	for _, step := range []*dfsSolverStep{
		d.getNextStep(q.puzzle, headUp, q.coord),
		d.getNextStep(q.puzzle, headRight, q.coord),
		d.getNextStep(q.puzzle, headDown, q.coord),
		d.getNextStep(q.puzzle, headLeft, q.coord),
	} {
		if isIncomplete, err := step.puzzle.IsIncomplete(step.coord); err != nil {
			continue
		} else if !isIncomplete {
			return step.puzzle, solvedPuzzle
		}

		g, sol := d.takeNextStepIntoDepthTowardsGoal(step)
		switch sol {
		case solvedPuzzle, foundGoal, badState:
			return g, sol
		}
		if step.coord == *d.goal {
			return g, foundGoal
		}
	}

	return q.puzzle, mayContinue
}
