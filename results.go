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

func runStandardSolver() {
	allPDs := reader.GetAllPuzzles()
	allSummaries := make([]summary, 0, len(allPDs))

	numBySize := make(map[int]map[model.Difficulty]int, 8)

	for _, pd := range allPDs {
		if _, ok := numBySize[pd.NumEdges]; !ok {
			numBySize[pd.NumEdges] = make(map[model.Difficulty]int, 3)
		}
		if pd.NumEdges > state.MaxEdges || pd.NumEdges != 25 || pd.Difficulty != model.Hard {
			continue
		}
		if strings.Contains(pd.String(), `893,598`) {
			// :badpokerface: this puzzle is destroying my machine. I'm skipping
			// it, and that makes me look bad:#
			continue
		}

		if !strings.Contains(pd.String(), `9,307,442`) {
			// this puzzle is busted?
			continue
		}

		if numBySize[pd.NumEdges][pd.Difficulty] >= sampleSize {
			continue
		} else if pd.NumEdges > 20 && pd.Difficulty == model.Hard && numBySize[pd.NumEdges][pd.Difficulty] >= numHard25s {
			continue
		}

		summ := runSolver(pd)
		allSummaries = append(allSummaries, summ)

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
) summary {

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
		log.Fatalf("Could not solve! %v: %s\n%s\n\n\n", err, sr, unsolvedStr)
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
	}
}
