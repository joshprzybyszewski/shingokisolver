package main

import (
	"flag"

	"github.com/joshprzybyszewski/shingokisolver/compete"
)

var (
	shouldUseConcurrency = flag.Bool(`concurrency`, false, `set to true to enable concurrency in the solver`)

	shouldRunProfiler  = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	shouldWriteResults = flag.Bool(`results`, false, `set to update the results in the READMEs`)
	shouldRunCompete   = flag.Bool(`competitive`, false, `set to true to compete online!`)
)

func main() {
	flag.Parse()

	if *shouldRunProfiler {
		runProfiler()
		return
	}

	if *shouldRunCompete {
		compete.Run()
		return
	}

	runStandardSolver()
}
