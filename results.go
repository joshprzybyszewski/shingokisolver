package main

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

var (
	veryHardPuzzles = []string{
		`893,598`,
		`3,225,837`,
		`5,070,205`,
		`5,155,284`,
		`5,996,280`,
		`9,506,048`,
	}
)

func shouldSkip(pd model.Definition) bool {
	if pd.NumEdges > state.MaxEdges {
		return true
	}
	for _, hardPID := range veryHardPuzzles {
		if strings.Contains(pd.String(), hardPID) {
			// :badpokerface: this puzzle is destroying my machine. I'm skipping
			// it, and that makes me look bad:#
			return true
		}
	}
	return false
}

func runStandardSolver() {
	allPDs := reader.GetAllPuzzles()
	allSummaries := make([]summary, 0, len(allPDs))

	numBySize := make(map[int]map[model.Difficulty]int, 8)

	for _, pd := range allPDs {
		if pd.Difficulty == model.Unknown {
			// don't solve "unknowns" because it bloats my stats
			continue
		}
		if _, ok := numBySize[pd.NumEdges]; !ok {
			numBySize[pd.NumEdges] = make(map[model.Difficulty]int, 3)
		}
		if shouldSkip(pd) {
			continue
		}

		if numBySize[pd.NumEdges][pd.Difficulty] >= sampleSize {
			continue
		} else if pd.NumEdges > 20 && pd.Difficulty == model.Hard && numBySize[pd.NumEdges][pd.Difficulty] >= numHard25s {
			continue
		}

		summ, solved := runSolver(pd)
		if solved {
			allSummaries = append(allSummaries, summ)
		} else if !(*shouldWriteResults) {
			panic(`unsolved puzzle`)
		}

		time.Sleep(100 * time.Millisecond)
		// collect garbage now, which should be that entire puzzle that we solved:#
		runtime.GC()

		numBySize[pd.NumEdges][pd.Difficulty]++
	}

	if *shouldWriteResults {
		updateReadme(allSummaries)
	}
}

func runSolver(
	pd model.Definition,
) (summary, bool) {

	log.Printf("Starting to solve %q...\n", pd.String())

	solve := solvers.Solve
	if *shouldUseConcurrency {
		solve = solvers.SolveConcurrently
	}

	var prevMemStats runtime.MemStats
	runtime.ReadMemStats(&prevMemStats)

	sr, err := solve(
		pd.NumEdges,
		pd.Nodes,
	)

	var rms runtime.MemStats
	runtime.ReadMemStats(&rms)

	unsolvedStr := puzzle.NewPuzzle(
		pd.NumEdges,
		pd.Nodes,
	).BlandString()

	if err != nil {
		log.Printf("Could not solve! %v: %s\n%s\n\n\n", err, sr, unsolvedStr)
		return summary{}, false
	} else {
		log.Printf("Solved: %s\n\n\n", sr)
	}

	return summary{
		Name:       pd.String(),
		NumEdges:   pd.NumEdges,
		Difficulty: pd.Difficulty,
		Duration:   sr.Duration,
		heapSize:   rms.TotalAlloc - prevMemStats.TotalAlloc,
		numGCs:     rms.NumGC - prevMemStats.NumGC,
		pauseNS:    rms.PauseTotalNs - prevMemStats.PauseTotalNs,
		Unsolved:   unsolvedStr,
		Solution:   sr.Puzzle.BlandSolution(),
	}, true
}
