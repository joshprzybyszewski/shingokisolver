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

func newDFSSolverForPartialSolution() *dfsSolver {
	return &dfsSolver{}
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

func (d *dfsSolver) solveForGoals(
	input *puzzle,
	start nodeCoord,
	goals []nodeCoord,
) (*puzzle, map[nodeCoord][]*puzzle, dfsGoalSolution) {

	if input == nil {
		return nil, nil, badState
	}

	ret := make(map[nodeCoord][]*puzzle, len(goals))
	for _, g := range goals {
		ret[g] = nil
	}

	p, state := d.takeNextStepIntoDepthTowardsGoals(
		ret,
		&dfsSolverStep{
			puzzle: input,
			coord:  start,
		},
	)
	if state == solvedPuzzle {
		return p, nil, state
	}
	return nil, ret, state
}

func (d *dfsSolver) takeNextStepIntoDepthTowardsGoals(
	found map[nodeCoord][]*puzzle,
	q *dfsSolverStep,
) (*puzzle, dfsGoalSolution) {
	if isIncomplete, err := q.puzzle.IsIncomplete(q.coord); err != nil {
		return nil, badState
	} else if !isIncomplete {
		return q.puzzle, solvedPuzzle
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
		if step == nil {
			continue
		}

		g, sol := d.takeNextStepIntoDepthTowardsGoals(found, step)
		switch sol {
		case solvedPuzzle:
			return g, sol
		}
		if slice, ok := found[step.coord]; ok {
			found[step.coord] = append(slice, g)
		}
	}

	return q.puzzle, mayContinue
}
