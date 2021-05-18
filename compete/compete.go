package compete

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func PopulateCache(
	size int,
	diff model.Difficulty,
	numDesired int,
) {
	if size > 25 {
		// these are daily/weekly/monthly "specials"
		numDesired = 1
	} else if numDesired < 1000 {
		// we have plenty
		return
	}

	for numGot := 0; numGot < numDesired; {
		_, err := getPuzzle(size, diff)
		if err != nil {
			log.Printf("getPuzzle errored: %+v", err)
		}
		numGot++
	}
}

func Run() {
	numPuzzles := 14
	sleepDur := time.Second

	for _, d := range []model.Difficulty{
		model.Easy,
		model.Medium,
		model.Hard,
	} {
		for _, s := range []int{
			// 5,
			// 7,
			10,
			15,
			20,
			25,
			30,
			// 35,
			// 40,
		} {
			for i := 0; i < numPuzzles; i++ {
				log.Printf("starting competition for %s %d edges", d, s)

				getAndSubmitPuzzle(s, d, i)

				// collect garbage now, which should be that entire puzzle that we solved:#
				runtime.GC()

				// between runs so we don't overwhelm
				// their servers or our machine accidentally:#
				time.Sleep(sleepDur)

				// another garbage collect can't hurt
				runtime.GC()
			}
		}
	}
}

func getAndSubmitPuzzle(
	size int,
	diff model.Difficulty,
	i int,
) {
	err := initLogs(fmt.Sprintf("%s_%dx%d_%d/", diff, size, size, i))
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
