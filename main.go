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

	puzzles, err := reader.CachedWebsitePuzzles()
	if err != nil {
		log.Printf("CachedWebsitePuzzles err: %+v\n", err)
		return
	}

	puzzles = append(puzzles, reader.DefaultPuzzles()...)

	for _, st := range solvers.AllSolvers {
		for _, pd := range puzzles {
			runSolver(st, pd)

			if *addPprof && (time.Since(t0) > 10*time.Second ||
				pd.NumEdges > 50) {
				return
			}
		}
	}
}

func runSolver(
	st solvers.SolverType,
	pd reader.PuzzleDef,
) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("caught panic: %+v", r)
	// 	}
	// }()
	if !strings.Contains(pd.String(), `6,483,955`) {
		// return
	}
	// if strings.Contains(pd.String(), `2,589,287`) || pd.NumEdges > 20 {
	// 	return
	// }

	if st != solvers.TargetSolverType {
		return
	}

	log.Printf("Starting to solve %q with %s...\n", pd.String(), st)
	s := solvers.NewSolver(
		pd.NumEdges,
		pd.Nodes,
		st,
	)

	sr, err := s.Solve()
	if err != nil {
		p := puzzle.NewPuzzle(
			pd.NumEdges,
			pd.Nodes,
		)
		log.Printf("%s could not solve! %v: %s\n%s\n\n\n", st, err, sr, p.String())
	} else {
		log.Printf("%s solved: %s\n\n\n", st, sr)
	}
}
