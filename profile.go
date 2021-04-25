package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

const (
	maxProfileDuration = 5 * time.Minute
	maxPuzzlesToSolve  = 50
)

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

	t0 := time.Now()
	numSolved := 0

	for _, pd := range reader.GetAllPuzzles() {
		if pd.NumEdges <= 20 {
			// only profile large puzzles. I don't care about small
			// ones because I can solve them so fast
			continue
		}
		// if !strings.Contains(pd.String(), `90,104`) {
		// 	continue
		// }

		runSolver(pd)
		numSolved++

		if time.Since(t0) > maxProfileDuration || numSolved >= maxPuzzlesToSolve {
			return
		}
	}
}
