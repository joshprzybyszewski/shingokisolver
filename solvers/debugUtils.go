package solvers

import (
	"fmt"
	"log"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

var (
	allIterationsUnder = 10000
	iterationsModulo   = 10000
)

func include(
	puzz *puzzle.Puzzle,
) bool {
	line2 := `(w 5)         (b 4) 2-> <- 1(XXX) 1-> <- 2(XXX)         (b 4) 1-> <- 1(XXX)`
	line4 := `(XXX)         (b 3) 1-> <- 1(b 2)         (XXX)         (XXX)         (XXX)`
	lUpl5 := `  5                          1`

	return strings.Contains(puzz.DebugString(), line2) &&
		strings.Contains(puzz.DebugString(), line4) &&
		strings.Contains(puzz.DebugString(), lUpl5)

}

var (
	seen = map[model.NodeCoord]struct{}{}
)

func printAllTargetsHit(
	caller string,
	puzz *puzzle.Puzzle,
	iterations int,
) {
	if !includeProgressLogs {
		return
	}
	shouldSkip := true
	if iterations < allIterationsUnder ||
		iterations%iterationsModulo == 0 {
		shouldSkip = false
	}
	if include(puzz) {
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

	log.Printf("\talpha:\n%+v\n", puzz.Alpha())
	// log.Printf("\tbeta:\n%+v\n", puzz.Beta())

	log.Printf("\tpuzzle:\n%s\n", puzz)
	fmt.Scanf("wait for acknowledgement")
}

var (
	hasSeen = false
)

func printPuzzleUpdate(
	caller string,
	depth int,
	puzz *puzzle.Puzzle,
	targeting model.Target,
	iterations int,
) {

	if !includeProgressLogs {
		return
	}
	shouldSkip := true
	if _, ok := seen[targeting.Coord]; ok {
		shouldSkip = false
		seen[targeting.Coord] = struct{}{}
	}
	if iterations < allIterationsUnder ||
		iterations%iterationsModulo == 0 {
		shouldSkip = false
	}
	if include(puzz) {
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

	log.Printf("\talpha:\n%+v\n", puzz.Alpha())
	// log.Printf("\tbeta:\n%+v\n", puzz.Beta())

	log.Printf("\tpuzzle:\n%s\n", puzz)
	fmt.Scanf("wait for acknowledgement")
}

func copyAndRemove(orig []model.NodeCoord, exclude model.NodeCoord) []model.NodeCoord {
	cpy := make([]model.NodeCoord, 0, len(orig)-1)
	for _, le := range orig {
		if le != exclude {
			cpy = append(cpy, le)
		}
	}
	return cpy
}
