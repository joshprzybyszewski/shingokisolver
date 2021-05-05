package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

const (
	maxProfileDuration = 5 * time.Minute
	maxPuzzlesToSolve  = 250

	cpuFileName = `solverProfile.pprof`
	memFileName = `solverMemory.pprof`
)

var (
	veryHardPuzzlesToProfile = []string{
		// `5,817,105`,
		// `5,996,280`,
		// `9,307,442`,
	}
)

func shouldSkipProfiling(pd model.Definition) bool {
	if pd.NumEdges > state.MaxEdges {
		return true
	}

	if pd.NumEdges <= 20 {
		// only profile large puzzles. I don't care about small
		// ones because I can solve them so fast
		return true
	}
	if pd.Difficulty != model.Hard {
		return true
	}

	for _, hardPID := range veryHardPuzzlesToProfile {
		if strings.Contains(pd.String(), hardPID) {
			// :badpokerface: this puzzle is destroying my machine. I'm skipping
			// it, and that makes me look bad:#
			return true
		}
	}
	return false
}

func runProfiler() {
	allPDs := reader.GetAllPuzzles()
	allSummaries := make([]summary, 0, maxPuzzlesToSolve)

	log.Printf("Starting a pprof...")

	if *shouldRunProfiler {
		log.Printf("\tincluding CPU profile.")
		cpuFile, err := os.Create(cpuFileName)
		if err != nil {
			log.Fatal("could not create CPU profile file: ", err)
		}
		defer cpuFile.Close()

		err = pprof.StartCPUProfile(cpuFile)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *includeMemoryProfile {
		log.Printf("\tincluding memory profile.")
		memFile, err := os.Create(memFileName)
		if err != nil {
			log.Fatal("could not create memory profile file: ", err)
		}
		defer memFile.Close()

		defer pprof.WriteHeapProfile(memFile)

		runtime.GC()
	}

	t0 := time.Now()

	for _, pd := range allPDs {
		if shouldSkipProfiling(pd) {
			continue
		}

		summ, solved := runSolver(pd)
		if solved {
			allSummaries = append(allSummaries, summ)
		}

		if time.Since(t0) > maxProfileDuration || len(allSummaries) >= maxPuzzlesToSolve {
			return
		}
	}
}
