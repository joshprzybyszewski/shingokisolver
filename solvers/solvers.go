package solvers

import (
	"errors"
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func Solve(
	size int,
	nl []model.NodeLocation,
) (SolvedResults, error) {
	if len(nl) == 0 {
		return SolvedResults{}, errors.New(`cannot solve a puzzle with no constraints`)
	}

	return solveWithTargets(size, nl)
}

type SolvedResults struct {
	Puzzle puzzle.Puzzle

	Duration time.Duration
}

func (sr SolvedResults) String() string {
	if sr.Puzzle.GetState(model.InvalidNodeCoord) != model.Complete {
		return fmt.Sprintf("Took %s. <no solution>\n", sr.Duration)
	}
	return fmt.Sprintf("Took %s\n%s\n", sr.Duration, sr.Puzzle.Solution())
}
