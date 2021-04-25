package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

var (
	runConcurrently = flag.Bool(`concurrency`, false, `set to true to enable concurrency in the solver`)

	addPprof           = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	shouldWriteResults = flag.Bool(`results`, false, `set to update the results in the READMEs`)
	runCompetitive     = flag.Bool(`competitive`, false, `set to true to compete online!`)
)

func main() {
	flag.Parse()

	if *addPprof {
		runProfiler()
		return
	}

	if *runCompetitive {
		compete.Run()
		return
	}

	runStandardSolver()
}

func runProfiler() {
	log.Printf("Starting a pprof...")
	f, err := os.Create(`solverProfile.pprof`)
	if err != nil {
		log.Fatal(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// TODO figure out GC
	// disable garbage collection entirely.
	// dangerous, I know.
	// debug.SetGCPercent(-1)
	t0 := time.Now()

	for i, pd := range reader.GetAllPuzzles() {
		if pd.NumEdges <= 20 { //  || !strings.Contains(pd.String(), `2,589,287`) {
			// if !strings.Contains(pd.String(), `5,434,778`) {
			continue
		}
		// if !strings.Contains(pd.String(), `90,104`) {
		// 	continue
		// }

		runSolver(pd)
		// go runSolver(pd)
		// time.Sleep(30 * time.Second)

		if time.Since(t0) > 5*60*time.Second || i > 50 {
			return
		}
	}
}

func runStandardSolver() {
	// TODO figure out GC
	// disable garbage collection entirely.
	// dangerous, I know.
	// debug.SetGCPercent(-1)

	allPDs := reader.GetAllPuzzles()
	allSummaries := make([]summary, 0, len(allPDs))

	numBySize := make(map[int]map[model.Difficulty]int, 8)

	for _, pd := range allPDs {
		if _, ok := numBySize[pd.NumEdges]; !ok {
			numBySize[pd.NumEdges] = make(map[model.Difficulty]int, 3)
		}
		if pd.NumEdges > state.MaxEdges {
			continue
		}
		// if !strings.Contains(pd.String(), `7,626,434`) {
		// 	continue
		// }

		if numBySize[pd.NumEdges][pd.Difficulty] >= sampleSize {
			continue
		} else if pd.NumEdges > 20 && pd.Difficulty == model.Hard && numBySize[pd.NumEdges][pd.Difficulty] >= 3 {
			continue
		}

		summ := runSolver(pd)
		allSummaries = append(allSummaries, summ)

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
	if *runConcurrently {
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
