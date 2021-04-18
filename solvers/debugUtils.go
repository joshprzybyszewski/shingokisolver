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

var (
	seen = map[model.NodeCoord]struct{}{}
)

func printPuzzleUpdate(
	caller string,
	depth int,
	puzz *puzzle.Puzzle,
	targeting *model.Target,
) {

	if !includeProgressLogs {
		return
	}

	log.Printf("printPuzzleUpdate")
	log.Printf("\tcaller:  \t%s",
		caller,
	)
	log.Printf("\tdepth:   \t%d",
		depth,
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

func copyAndRemove(orig []model.NodeCoord, exclude model.NodeCoord) []model.NodeCoord {
	cpy := make([]model.NodeCoord, 0, len(orig)-1)
	for _, le := range orig {
		if le != exclude {
			cpy = append(cpy, le)
		}
	}
	return cpy
}
