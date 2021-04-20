package compete

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func PopulateCache(
	size int,
	numDesired int,
) {
	for numGot := 0; numGot < numDesired; {
		for _, d := range []difficulty{
			easy,
			medium,
			hard,
		} {
			_, err := getPuzzle(size, d)
			if err != nil {
				log.Printf("getPuzzle errored: %+v", err)
			}
			numGot++
		}
	}
}

func Run() {
	// TODO figure out GC
	// disable garbage collection entirely.
	// dangerous, I know.
	// debug.SetGCPercent(-1)

	for _, s := range []int{
		5,
		7,
		// 10,
		// 15,
		// 20,
		25,
	} {
		for _, d := range []difficulty{
			// easy,
			// medium,
			hard,
		} {
			log.Printf("starting competition for %s %d edges", d, s)

			getAndSubmitPuzzle(s, d)

			// wait 10 seconds between runs so we don't overwhelm
			// their servers or our machine accidentally:#
			time.Sleep(10 * time.Second)

			// collect garbage now, which should be that entire puzzle that we solved:#
			runtime.GC()
		}
	}
}

func getAndSubmitPuzzle(
	size int,
	diff difficulty,
) {
	err := initLogs(fmt.Sprintf("%s_%dx%d/", diff, size, size))
	if err != nil {
		panic(err)
	}
	defer flushLogs()

	wp, err := getPuzzle(
		size,
		diff,
	)
	if err != nil {
		panic(err)
	}

	sr, err := solvers.Solve(
		wp.pd.NumEdges,
		wp.pd.Nodes,
	)
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
