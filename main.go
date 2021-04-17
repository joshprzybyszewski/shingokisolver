package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

var (
	addPprof            = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	includeProgressLogs = flag.Bool(`includeProcessLogs`, false, `set to see each solver's progress logs`)
	runCompetitive      = flag.Bool(`competitive`, false, `set to true to get a puzzle from the internet and submit a response`)
)

func main() {
	flag.Parse()
	if *includeProgressLogs {
		puzzle.AddProgressLogs()
		solvers.AddProgressLogs()
	}

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

	if *runCompetitive {
		compete.Run()
		return
	}

	t0 := time.Now()

	for _, pd := range reader.GetAllPuzzles() {
		if !strings.Contains(pd.String(), `5,434,778`) {
			continue
		}
		go runSolver(pd)
		time.Sleep(30 * time.Second)

		if *addPprof && (time.Since(t0) > 10*time.Second ||
			pd.NumEdges > 50) {
			return
		}
	}
}

func runSolver(
	pd reader.PuzzleDef,
) {

	log.Printf("Starting to solve %q...\n", pd.String())
	s := solvers.NewSolver(
		pd.NumEdges,
		pd.Nodes,
	)

	sr, err := s.Solve()
	if err != nil {
		p := puzzle.NewPuzzle(
			pd.NumEdges,
			pd.Nodes,
		)
		log.Printf("Could not solve! %v: %s\n%s\n\n\n", err, sr, p.String())
	} else {
		log.Printf("Solved: %s\n\n\n", sr)
	}
}
