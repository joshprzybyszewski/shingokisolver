package solvers

import (
	"errors"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type solver interface {
	solve() (*puzzle.Puzzle, bool)
}

func newSolver(
	size int,
	nl []model.NodeLocation,
	st SolverType,
) solver {
	switch st {
	case TargetSolverType:
		return newTargetSolver(size, nl)
	default:
		return newTargetSolver(size, nl)
	}
}

type solverWrapper struct {
	s solver
}

func newWrapper(s solver) Solver {
	return &solverWrapper{
		s: s,
	}
}

func (sw *solverWrapper) Solve() (SolvedResults, error) {
	t0 := time.Now()

	puzz, isSolved := sw.s.solve()
	if !isSolved {
		return SolvedResults{
			Duration: time.Since(t0),
		}, errors.New(`unsolvable`)
	}

	return SolvedResults{
		Puzzle:   puzz,
		Duration: time.Since(t0),
	}, nil
}
