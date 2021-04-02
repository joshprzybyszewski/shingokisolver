package solvers

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) getNextStep(
	puzz *puzzle.Puzzle,
	move model.Cardinal,
	nc model.NodeCoord,
) (*puzzle.Puzzle, model.NodeCoord) {

	newCoord, newP, err := puzz.AddEdge(move, nc)
	if err != nil {
		return nil, model.NodeCoord{}
	}

	if newP.IsRangeInvalid(
		newCoord.Row-1,
		newCoord.Row+1,
		newCoord.Col-1,
		newCoord.Col+1,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil, model.NodeCoord{}
	}

	return newP, newCoord
}

type dfsGoalSolution int

const (
	mayContinue  dfsGoalSolution = 0
	badState     dfsGoalSolution = 1
	solvedPuzzle dfsGoalSolution = 2
	foundGoal    dfsGoalSolution = 3
)

func (d *targetSolver) solveForGoals(
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
		input.DeepCopy(),
		start,
	)
	return p, ret, state
}

func (d *targetSolver) takeNextStepIntoDepthTowardsGoals(
	puzzlesByTargetedLooseEnd map[model.NodeCoord][]*puzzle.Puzzle,
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
) (*puzzle.Puzzle, dfsGoalSolution) {
	if puzz == nil {
		return nil, mayContinue
	}

	if len(puzzlesByTargetedLooseEnd) <= 1 {
		// there's only one target node. that means we need to check for violations
		// and walk the path to verify it's complete
		if isIncomplete, err := puzz.IsIncomplete(fromCoord); err != nil {
			return nil, badState
		} else if !isIncomplete {
			return puzz, solvedPuzzle
		}
	} else {
		// there's more than 1 possible goal nodes. instead of walking the
		// whole path, let's just check to see if the puzzle has been violated
		if violates, err := puzz.IsInViolation(fromCoord); err != nil || violates {
			return nil, badState
		}
	}

	if slice, ok := puzzlesByTargetedLooseEnd[fromCoord]; ok {
		puzzlesByTargetedLooseEnd[fromCoord] = append(slice, puzz.DeepCopy())
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n",
			fromCoord.Row, fromCoord.Col,
			d.numProcessed,
			puzz.String(),
		)
		fmt.Scanf("hello there")
	}

	for _, nextHeading := range model.AllCardinals {
		nextPuzz, nextCoord := d.getNextStep(
			puzz.DeepCopy(),
			nextHeading,
			fromCoord,
		)
		g, sol := d.takeNextStepIntoDepthTowardsGoals(
			puzzlesByTargetedLooseEnd,
			nextPuzz,
			nextCoord,
		)

		switch sol {
		case solvedPuzzle:
			return g, sol
		}
	}

	return nil, mayContinue
}
