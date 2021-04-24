package solvers

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func SolveConcurrently(
	size int,
	nl []model.NodeLocation,
) (SolvedResults, error) {
	if len(nl) == 0 {
		return SolvedResults{}, errors.New(`cannot solve a puzzle with no constraints`)
	}

	cs := concurrentSolver{}
	defer func(cs *concurrentSolver) {
		log.Printf("solved with %d queued payloads and %d processed.", cs.numPayloads, cs.numProcessed)
	}(&cs)

	return cs.solve(puzzle.NewPuzzle(size, nl))
}

func Solve(
	size int,
	nl []model.NodeLocation,
) (SolvedResults, error) {
	if len(nl) == 0 {
		return SolvedResults{}, errors.New(`cannot solve a puzzle with no constraints`)
	}

	return solvePuzzleByTargets(puzzle.NewPuzzle(size, nl))
}

type SolvedResults struct {
	Puzzle puzzle.Puzzle

	Duration time.Duration
}

func (sr SolvedResults) String() string {
	if sr.Puzzle.GetState() != model.Complete {
		return fmt.Sprintf("Took %s. <no solution>\n", sr.Duration)
	}
	return fmt.Sprintf("Took %s\n%s\n", sr.Duration, sr.Puzzle.Solution())
}
