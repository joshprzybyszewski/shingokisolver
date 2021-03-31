package compete

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func Run() {
	pd := getPuzzle()

	s := solvers.NewSolver(
		pd.NumEdges,
		pd.Nodes,
		solvers.TargetSolverType,
	)

	sr, err := s.Solve()
	if err != nil {
		p := puzzle.NewPuzzle(
			pd.NumEdges,
			pd.Nodes,
		)
		log.Printf("derp. Couldn't solve. %v\n%s\n", err, p)
		return
	}

	submitAnswer(sr)
}
