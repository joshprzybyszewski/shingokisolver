package compete

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func Run() {
	wp, err := getPuzzle(5)
	if err != nil {
		panic(err)
	}

	log.Printf("wp: %+v\n", wp)

	s := solvers.NewSolver(
		wp.pd.NumEdges,
		wp.pd.Nodes,
		solvers.TargetSolverType,
	)

	sr, err := s.Solve()
	if err != nil {
		p := puzzle.NewPuzzle(
			wp.pd.NumEdges,
			wp.pd.Nodes,
		)
		log.Printf("derp. Couldn't solve. %v\n%s\n", err, p)
		return
	}

	submitAnswer(sr)
}
