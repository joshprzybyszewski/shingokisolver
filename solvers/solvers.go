package solvers

import (
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type Solver interface {
	Solve() (SolvedResults, error)
}

func NewSolver(
	size int,
	nl []model.NodeLocation,
) Solver {
	if len(nl) == 0 {
		panic(`cannot solve a puzzle with no nodes!`)
	}

	return newTargetSolver(size, nl)
}

type SolvedResults struct {
	Puzzle *puzzle.Puzzle

	Duration time.Duration
}

func (sr SolvedResults) String() string {
	if sr.Puzzle == nil {
		return fmt.Sprintf("Took %s. <no solution>\n", sr.Duration)
	}
	return fmt.Sprintf("Took %s\n%s\n", sr.Duration, sr.Puzzle.Solution())
}
