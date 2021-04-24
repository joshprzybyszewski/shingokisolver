package solvers

import (
	"log"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

func (cs *concurrentSolver) printPuzzleUpdate(
	caller string,
	puzz puzzle.Puzzle,
	targeting model.Target,
) {
	if !includeProgressLogs {
		return
	}

}

func (cs *concurrentSolver) printPayload(
	caller string,
	payload unsolvedPayload,
) {
	if !includeProgressLogs {
		return
	}

	log.Printf("printPuzzleUpdate")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\tisNodesComplete:\t%+v",
		payload.isNodesComplete,
	)

	// log.Printf("\talpha:\n%+v\n", payload.puzz.Alpha())
	// // log.Printf("\tbeta:\n%+v\n", puzz.Beta())

	log.Printf("\tpuzzle:\n%s\n", payload.puzz)
	time.Sleep(100 * time.Millisecond)
}

func printPuzzleUpdate(
	caller string,
	puzz puzzle.Puzzle,
	targeting model.Target,
) {

	if !includeProgressLogs {
		return
	}

	log.Printf("printPuzzleUpdate")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\ttargeting:\t%+v",
		targeting,
	)

	log.Printf("\talpha:\n%+v\n", puzz.Alpha())
	// log.Printf("\tbeta:\n%+v\n", puzz.Beta())

	log.Printf("\tpuzzle:\n%s\n", puzz)
	time.Sleep(100 * time.Millisecond)
	// fmt.Scanf("wait for acknowledgement")
}
