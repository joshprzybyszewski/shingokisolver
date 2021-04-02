package solvers

import (
	"fmt"
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) getNextStep(
	puzz *puzzle.Puzzle,
	move model.Cardinal,
	nc model.NodeCoord,
) (model.NodeCoord, *puzzle.Puzzle) {

	newCoord, newP, nextState := puzz.AddEdge(move, nc)

	switch nextState {
	case model.Violation, model.Unexpected:
		return model.NodeCoord{}, nil
	case model.Duplicate:
		return model.NodeCoord{}, nil
	}

	return newCoord, newP
}

func (d *targetSolver) solveForGoals(
	input *puzzle.Puzzle,
	start model.NodeCoord,
	goals []model.NodeCoord,
) (*puzzle.Puzzle, map[model.NodeCoord][]*puzzle.Puzzle, model.State) {

	printAllTargetsHit(fmt.Sprintf(`solveForGoals(%+v)`, start), input, d.iterations())
	if input == nil {
		return nil, nil, model.Unexpected
	}

	puzzlesByTargetedLooseEnd := make(map[model.NodeCoord][]*puzzle.Puzzle, len(goals))
	for _, g := range goals {
		puzzlesByTargetedLooseEnd[g] = nil
	}
	delete(puzzlesByTargetedLooseEnd, start)

	p, state := d.takeNextStepIntoDepthTowardsGoals(
		puzzlesByTargetedLooseEnd,
		input.DeepCopy(),
		start,
	)
	return p, puzzlesByTargetedLooseEnd, state
}

func (d *targetSolver) takeNextStepIntoDepthTowardsGoals(
	puzzlesByTargetedLooseEnd map[model.NodeCoord][]*puzzle.Puzzle,
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
) (*puzzle.Puzzle, model.State) {
	if puzz == nil {
		return nil, model.Unexpected
	}

	switch s := puzz.GetState(); s {
	case model.Complete:
		return puzz, s
	case model.Incomplete:
		// continue in the func
	default:
		return nil, s
	}

	if slice, ok := puzzlesByTargetedLooseEnd[fromCoord]; ok {
		puzzlesByTargetedLooseEnd[fromCoord] = append(slice, puzz.DeepCopy())
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		log.Printf("takeNextStepIntoDepthTowardsGoals about to process (%+v): %d\n%s\n",
			fromCoord,
			d.numProcessed,
			puzz.String(),
		)
		fmt.Scanf("hello there")
	}

	for _, nextHeading := range model.AllCardinals {
		nextCoord, nextPuzz := d.getNextStep(
			puzz.DeepCopy(),
			nextHeading,
			fromCoord,
		)

		retPuzz, s := d.takeNextStepIntoDepthTowardsGoals(
			puzzlesByTargetedLooseEnd,
			nextPuzz,
			nextCoord,
		)

		switch s {
		case model.Complete:
			return retPuzz, s
		}
	}

	return nil, model.Incomplete
}
