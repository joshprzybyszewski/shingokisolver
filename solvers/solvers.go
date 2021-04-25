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
	defer cs.logMeta()

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
	if state := sr.Puzzle.GetState(); state != model.Complete {
		log.Printf("Took %s. %s <no solution>\n%s\n", sr.Duration, state, sr.Puzzle.String())
		return fmt.Sprintf("Took %s. %s <no solution>\n%s\n", sr.Duration, state, sr.Puzzle.String())
	}

	return fmt.Sprintf("Took %s\n%s\n", sr.Duration, sr.Puzzle.Solution())
}
