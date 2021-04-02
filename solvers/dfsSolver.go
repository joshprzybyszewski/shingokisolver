package solvers

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type dfsSolverStep struct {
	puzzle *puzzle.Puzzle
	coord  model.NodeCoord
}

type dfsSolver struct {
	puzzle *puzzle.Puzzle

	goal *model.NodeCoord

	numProcessed int
}

func newDFSSolver(
	size int,
	nl []model.NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &dfsSolver{
		puzzle: puzzle.NewPuzzle(size, nl),
	}
}

func newDFSSolverForPartialSolution() *dfsSolver {
	return &dfsSolver{}
}

func (d *dfsSolver) iterations() int {
	return d.numProcessed
}

func (d *dfsSolver) solve() (*puzzle.Puzzle, bool) {
	return d.takeNextStepIntoDepth(&dfsSolverStep{
		puzzle: d.puzzle,
		coord:  d.puzzle.GetCoordForHighestValueNode(),
	})
}

func (d *dfsSolver) takeNextStepIntoDepth(
	q *dfsSolverStep,
) (*puzzle.Puzzle, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.Row, q.coord.Col, d.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.puzzle.IsIncomplete(q.coord); err != nil {
		return q.puzzle, false
	} else if !isIncomplete {
		return q.puzzle, true
	}

	for _, step := range []*dfsSolverStep{
		d.getNextStep(q.puzzle, model.HeadUp, q.coord),
		d.getNextStep(q.puzzle, model.HeadRight, q.coord),
		d.getNextStep(q.puzzle, model.HeadDown, q.coord),
		d.getNextStep(q.puzzle, model.HeadLeft, q.coord),
	} {
		g, isSolved := d.takeNextStepIntoDepth(step)
		if isSolved {
			return g, true
		}
	}

	return q.puzzle, false
}

func (d *dfsSolver) getNextStep(
	puzz *puzzle.Puzzle,
	move model.Cardinal,
	nc model.NodeCoord,
) *dfsSolverStep {

	newCoord, newP, err := puzz.DeepCopy().AddEdge(move, nc)
	if err != nil {
		return nil
	}

	if newP.IsRangeInvalid(
		newCoord.Row-1,
		newCoord.Row+1,
		newCoord.Col-1,
		newCoord.Col+1,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil
	}

	return &dfsSolverStep{
		puzzle: newP,
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
	input *puzzle.Puzzle,
	start model.NodeCoord,
	goals []model.NodeCoord,
) (*puzzle.Puzzle, map[model.NodeCoord][]*puzzle.Puzzle, dfsGoalSolution) {

	if input == nil {
		return nil, nil, badState
	}

	ret := make(map[model.NodeCoord][]*puzzle.Puzzle, len(goals))
	for _, g := range goals {
		ret[g] = nil
	}
	delete(ret, start)

	p, state := d.takeNextStepIntoDepthTowardsGoals(
		ret,
		dfsSolverStep{
			puzzle: input,
			coord:  start,
		},
	)
	return p, ret, state
}

func (d *dfsSolver) takeNextStepIntoDepthTowardsGoals(
	puzzlesByTargetedLooseEnd map[model.NodeCoord][]*puzzle.Puzzle,
	step dfsSolverStep,
) (*puzzle.Puzzle, dfsGoalSolution) {
	if step.puzzle == nil {
		return nil, mayContinue
	}

	if len(puzzlesByTargetedLooseEnd) <= 1 {
		// there's only one target node. that means we need to check for violations
		// and walk the path to verify it's complete
		if isIncomplete, err := step.puzzle.IsIncomplete(step.coord); err != nil {
			return nil, badState
		} else if !isIncomplete {
			return step.puzzle, solvedPuzzle
		}
	} else {
		// there's more than 1 possible goal nodes. instead of walking the
		// whole path, let's just check to see if the puzzle has been violated
		if violates, err := step.puzzle.IsInViolation(step.coord); err != nil || violates {
			return nil, badState
		}
	}

	if slice, ok := puzzlesByTargetedLooseEnd[step.coord]; ok {
		puzzlesByTargetedLooseEnd[step.coord] = append(slice, step.puzzle.DeepCopy())
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n",
			step.coord.Row, step.coord.Col,
			d.numProcessed,
			step.puzzle.String(),
		)
		fmt.Scanf("hello there")
	}

	for _, nextHeading := range model.AllCardinals {
		nextStep := d.getNextStep(step.puzzle, nextHeading, step.coord)
		if nextStep == nil {
			continue
		}
		g, sol := d.takeNextStepIntoDepthTowardsGoals(
			puzzlesByTargetedLooseEnd,
			*nextStep,
		)

		switch sol {
		case solvedPuzzle:
			return g, sol
		}
	}

	return nil, mayContinue
}
