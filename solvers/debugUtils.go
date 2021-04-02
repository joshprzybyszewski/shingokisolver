package solvers

import (
	"fmt"
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

var (
	seen = map[model.NodeCoord]struct{}{}
)

func printAllTargetsHit(
	caller string,
	puzzle *puzzle.Puzzle,
	looseEnds []model.NodeCoord,
	iterations int,
) {
	if !includeProgressLogs {
		return
	}
	shouldSkip := true
	if iterations < 10 || iterations%10000 == 0 {
		shouldSkip = false
	}
	if shouldSkip {
		return
	}

	log.Printf("printAllTargetsHit")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\titerations:\t%d",
		iterations,
	)
	log.Printf("\tlooseEnds:\t%+v",
		looseEnds,
	)
	log.Printf("\tpuzzle:\n%s\n", puzzle)
	fmt.Scanf("wait for acknowledgement")
}

func printPuzzleUpdate(
	caller string,
	depth int,
	puzzle *puzzle.Puzzle,
	targeting model.Target,
	looseEnds []model.NodeCoord,
	iterations int,
) {

	if !includeProgressLogs {
		return
	}
	shouldSkip := false
	if _, ok := seen[targeting.Coord]; ok {
		shouldSkip = false
		seen[targeting.Coord] = struct{}{}
	}
	if iterations < 10 || iterations%10000 == 0 {
		shouldSkip = false
	}
	if shouldSkip {
		return
	}

	log.Printf("printPuzzleUpdate")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\tdepth:   \t%d",
		depth,
	)
	log.Printf("\titerations:\t%d",
		iterations,
	)
	log.Printf("\ttargeting:\t%+v",
		targeting,
	)
	log.Printf("\tlooseEnds:\t%+v",
		looseEnds,
	)
	log.Printf("\tpuzzle:\n%s\n", puzzle)
	fmt.Scanf("wait for acknowledgement")
}
