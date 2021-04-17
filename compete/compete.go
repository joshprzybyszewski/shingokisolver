package compete

import (
	"log"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func Run() {
	for _, s := range []int{
		// 5, 7, 10, 15, 20,
		25,
	} {
		for _, d := range []difficulty{
			easy, medium, hard,
		} {
			log.Printf("starting competition for %s %d edges", d, s)

			getAndSubmitPuzzle(s, d)

			// wait 10 seconds between runs so we don't overwhelm
			// their servers or our machine accidentally:#
			time.Sleep(10 * time.Second)
		}
	}
}

func getAndSubmitPuzzle(
	size int,
	diff difficulty,
) {
	initLogs()
	defer flushLogs()

	wp, err := getPuzzle(
		size,
		diff,
	)
	if err != nil {
		panic(err)
	}

	s := solvers.NewSolver(
		wp.pd.NumEdges,
		wp.pd.Nodes,
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

	err = submitAnswer(wp, sr)
	if err != nil {
		log.Printf("submitAnswer errored: %v", err)
	}

	log.Printf("Solved in %s:\n%s\n ", sr.Duration, sr.Puzzle.Solution())
}
