package main

import (
	"flag"

	"github.com/joshprzybyszewski/shingokisolver/compete"
)

var (
	shouldUseConcurrency = flag.Bool(`concurrency`, false, `set to true to enable concurrency in the solver`)

	shouldRunProfiler    = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	includeMemoryProfile = flag.Bool(`includeMemProf`, false, `set if you'd like to include a memory profile output`)

	shouldWriteResults = flag.Bool(`results`, false, `set to update the results in the READMEs`)
	shouldRunCompete   = flag.Bool(`competitive`, false, `set to true to compete online!`)
)

func main() {
	flag.Parse()

	if *shouldRunProfiler || *includeMemoryProfile {
		runProfiler()
		return
	}

	if *shouldRunCompete {
		compete.Run()
		return
	}

	runStandardSolver()

	/* fetch the daily/weekly/monthly special
	compete.PopulateCache(30, model.Easy, 5)
	compete.PopulateCache(35, model.Easy, 5)
	compete.PopulateCache(40, model.Easy, 5)
	*/
}
