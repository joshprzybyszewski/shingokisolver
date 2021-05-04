package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

const (
	maxProfileDuration = 5 * time.Second
	maxPuzzlesToSolve  = 50

	cpuFileName = `solverProfile.pprof`
	memFileName = `solverMemory.pprof`
)

func runProfiler() {
	log.Printf("Starting a pprof...")
	f, err := os.Create(cpuFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	var memFile *os.File
	if *includeMemoryProfile {
		memFile, err = os.Create(memFileName)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
	}

	t0 := time.Now()
	numSolved := 0

	for _, pd := range reader.GetAllPuzzles() {
		if pd.NumEdges <= 20 {
			// only profile large puzzles. I don't care about small
			// ones because I can solve them so fast
			continue
		}
		if !strings.Contains(pd.String(), `5,817,105`) {
			continue
		}

		if *includeMemoryProfile {
			runtime.GC()
		}

		go runSolver(pd)
		time.Sleep(time.Minute)
		numSolved++

		if *includeMemoryProfile {
			if err := pprof.WriteHeapProfile(memFile); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			return
		}

		if time.Since(t0) > maxProfileDuration || numSolved >= maxPuzzlesToSolve {
			return
		}
	}
}
