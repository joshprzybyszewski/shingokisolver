// +build debug

package solvers

import (
	// "fmt"
	"log"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (cs *concurrentSolver) printPuzzleUpdate(
	caller string,
	puzz puzzle.Puzzle,
	targeting model.Target,
) {
	log.Printf("cs.printPuzzleUpdate")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\ttarget: %v",
		targeting,
	)
	cs.logMeta()

	log.Printf("\tpuzzle:\n%s\n", puzz)
	time.Sleep(100 * time.Millisecond)
	// fmt.Scanf("wait")
}

func (cs *concurrentSolver) printPayload(
	caller string,
	payload targetPayload,
) {
	log.Printf("cs.printPayload")
	log.Printf("\tcaller: %s",
		caller,
	)
	log.Printf("\ttarget: %v",
		payload.target,
	)
	cs.logMeta()

	log.Printf("\tpuzzle:\n%s\n", payload.puzz)
	time.Sleep(100 * time.Millisecond)
	// fmt.Scanf("wait")
}

func (cs *concurrentSolver) printFlippingPayload(
	caller string,
	payload flippingPayload,
) {
	log.Printf("cs.printFlippingPayload")
	log.Printf("\tnextUnknown:         %s",
		payload.nextUnknown,
	)
	cs.logMeta()
	log.Printf("\tcaller:              %s",
		caller,
	)

	log.Printf("\tpuzzle:\n%s\n", payload.puzz)
	time.Sleep(100 * time.Millisecond)
	// fmt.Scanf("wait")
}

func (cs *concurrentSolver) logMeta() {
	log.Printf("concurrentSolver.logMeta()")
	log.Printf("\tnumImmediates:   %d",
		cs.numImmediates,
	)
	log.Printf("\tnumTargetsAdded:     %d",
		cs.numTargetsAdded,
	)
	log.Printf("\tnumTargetsProcessed: %d",
		cs.numTargetsProcessed,
	)
	log.Printf("\tnumFlipsAdded:       %d",
		cs.numFlipsAdded,
	)
	log.Printf("\tnumFlipsProcessed:   %d",
		cs.numFlipsProcessed,
	)
}

func printPuzzleUpdate(
	caller string,
	puzz puzzle.Puzzle,
	targeting model.Target,
) {

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
	// fmt.Scanf("wait")
	// fmt.Scanf("wait for acknowledgement")
}
