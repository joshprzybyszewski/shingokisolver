package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
)

var (
	addPprof            = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	includeProgressLogs = flag.Bool(`includeProcessLogs`, false, `set to see each solver's progress logs`)
)

func main() {
	flag.Parse()
	puzzle.IncludeProgressLogs = *includeProgressLogs

	if *addPprof {
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
	}

	t0 := time.Now()

	for _, pd := range reader.DefaultPuzzles() {
		for _, st := range puzzle.AllSolvers {
			fmt.Printf("Starting to solve with %s...\n", st)
			s := puzzle.NewSolver(
				pd.NumEdges,
				pd.Nodes,
				st,
			)

			sr, err := s.Solve()
			if err != nil {
				fmt.Printf("%s could not solve! %v: %s\n\n\n", st, err, sr)
				continue
			}

			fmt.Printf("%s solved: %s\n\n\n", st, sr)
			if time.Since(t0) > 30*time.Second {
				break
			}
		}
	}
}
