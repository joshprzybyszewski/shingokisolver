package solvers

import (
	"fmt"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

type SolverType int

const (
	BFSSolverType              SolverType = 1
	DFSSolverType              SolverType = 2
	HeadsAndTailsDFSSolverType SolverType = 3
	TargetSolverType           SolverType = 4
)

func (st SolverType) String() string {
	switch st {
	case DFSSolverType:
		return `DFSSolverType`
	case HeadsAndTailsDFSSolverType:
		return `HeadsAndTailsDFSSolverType`
	case BFSSolverType:
		return `BFSSolverType`
	case TargetSolverType:
		return `TargetSolverType`
	default:
		return `Unknown Solver`
	}
}

var (
	AllSolvers = []SolverType{
		TargetSolverType,
		DFSSolverType,
		HeadsAndTailsDFSSolverType,
		BFSSolverType,
	}
)

type Solver interface {
	Solve() (SolvedResults, error)
}

func NewSolver(
	size int,
	nl []model.NodeLocation,
	st SolverType,
) Solver {
	return newWrapper(newSolver(size, nl, st))
}

type SolvedResults struct {
	Puzzle        *puzzle.Puzzle
	NumIterations int
	Duration      time.Duration
}

func (sr SolvedResults) String() string {
	if sr.Puzzle == nil {
		return fmt.Sprintf("(%d iterations in %s) <no solution>\n", sr.NumIterations, sr.Duration.String())
	}
	return fmt.Sprintf("(%d iterations in %s)\n%s\n", sr.NumIterations, sr.Duration.String(), sr.Puzzle.String())
}
