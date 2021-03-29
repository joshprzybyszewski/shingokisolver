package solvers

import (
	"errors"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type solver interface {
	solve() (*puzzle.Puzzle, bool)
	iterations() int
}

func newSolver(
	size int,
	nl []model.NodeLocation,
	st SolverType,
) solver {
	switch st {
	case BFSSolverType:
		return newBFSSolver(size, nl)
	case DFSSolverType:
		return newDFSSolver(size, nl)
	case HeadsAndTailsDFSSolverType:
		return newHeadsAndTailsDFSSolver(size, nl)
	case TargetSolverType:
		return newTargetSolver(size, nl)
	default:
		return newDFSSolver(size, nl)
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

	g, isSolved := sw.s.solve()
	if !isSolved {
		return SolvedResults{
			NumIterations: sw.s.iterations(),
			Duration:      time.Since(t0),
		}, errors.New(`bfsSolver unsolvable`)
	}

	return SolvedResults{
		Drawing:       g.String(),
		NumIterations: sw.s.iterations(),
		Duration:      time.Since(t0),
	}, nil
}
