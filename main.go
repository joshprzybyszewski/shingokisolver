package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

var (
	addPprof            = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	includeProgressLogs = flag.Bool(`includeProcessLogs`, false, `set to see each solver's progress logs`)
	shouldWriteResults  = flag.Bool(`results`, false, `set to update the results in the READMEs`)
	runCompetitive      = flag.Bool(`competitive`, false, `set to true to get a puzzle from the internet and submit a response`)
)

func main() {
	flag.Parse()

	if *includeProgressLogs {
		puzzle.AddProgressLogs()
		solvers.AddProgressLogs()
	}

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

	for _, pd := range reader.GetAllPuzzles() {
		if pd.NumEdges > 15 { //  || !strings.Contains(pd.String(), `2,589,287`) {
			// if !strings.Contains(pd.String(), `5,434,778`) {
			continue
		}

		runSolver(pd)
		// go runSolver(pd)
		// time.Sleep(30 * time.Second)

		if time.Since(t0) > 30*time.Second {
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

	numBySize := make(map[int]int, 8)
	maxPerSize := 10

	for _, pd := range allPDs {
		if pd.NumEdges > 20 {
			continue
		}

		if numBySize[pd.NumEdges] > maxPerSize {
			continue
		}

		// if !strings.Contains(pd.String(), `5,937,602`) {
		// 	continue
		// }

		summ := runSolver(pd)
		allSummaries = append(allSummaries, summ)

		// collect garbage now, which should be that entire puzzle that we solved:#
		runtime.GC()

		numBySize[pd.NumEdges] += 1
	}

	if *shouldWriteResults {
		updateReadme(allSummaries)
	}
}

func runSolver(
	pd reader.PuzzleDef,
) summary {

	log.Printf("Starting to solve %q...\n", pd.String())

	var prevMemStats runtime.MemStats
	runtime.ReadMemStats(&prevMemStats)

	sr, err := solvers.Solve(
		pd.NumEdges,
		pd.Nodes,
	)

	var rms runtime.MemStats
	runtime.ReadMemStats(&rms)

	unsolvedStr := puzzle.NewPuzzle(
		pd.NumEdges,
		pd.Nodes,
	).String()

	if err != nil {
		log.Fatalf("Could not solve! %v: %s\n%s\n\n\n", err, sr, unsolvedStr)
	} else {
		log.Printf("Solved: %s\n\n\n", sr)
	}

	return summary{
		Name:     pd.String(),
		NumEdges: pd.NumEdges,
		Duration: sr.Duration,
		heapSize: rms.TotalAlloc - prevMemStats.TotalAlloc,
		numGCs:   rms.NumGC - prevMemStats.NumGC,
		pauseNS:  rms.PauseTotalNs - prevMemStats.PauseTotalNs,
		Unsolved: unsolvedStr,
		Solution: sr.Puzzle.Solution(),
	}
}
