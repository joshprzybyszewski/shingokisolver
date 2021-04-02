package solvers

import (
	"fmt"
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type twoArms struct {
	arm1Len int8
	arm2Len int8

	arm1Heading model.Cardinal
	arm2Heading model.Cardinal
}

func buildTwoArmOptions(n model.Node) []twoArms {
	var options []twoArms

	for i, feeler1 := range model.AllCardinals {
		for _, feeler2 := range model.AllCardinals[i+1:] {
			if n.IsInvalidMotions(feeler1, feeler2) {
				continue
			}

			for arm1 := int8(1); arm1 <= n.Value()/2; arm1++ {
				options = append(options, twoArms{
					arm1Heading: feeler1,
					arm2Heading: feeler2,
					arm1Len:     arm1,
					arm2Len:     n.Value() - arm1,
				})
			}
		}
	}

	return options
}

// eliminates loose ends that don't actually exist
// leaves the looseEnds slice in the order that it had previously
func getLooseEndsWithoutDuplicates(looseEnds []model.NodeCoord) []model.NodeCoord {

	numExisting := make(map[model.NodeCoord]int, len(looseEnds))
	for _, le := range looseEnds {
		numExisting[le] += 1
	}

	looseEndsDeduped := make([]model.NodeCoord, 0, len(looseEnds))
	for _, le := range looseEnds {
		if existing := numExisting[le]; existing%2 == 0 {
			looseEndsDeduped = append(looseEndsDeduped, le)
		}
	}

	return looseEndsDeduped
}

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
	targeting target,
	looseEnds []model.NodeCoord,
	iterations int,
) {

	if !includeProgressLogs {
		return
	}
	shouldSkip := false
	if _, ok := seen[targeting.coord]; ok {
		shouldSkip = false
		seen[targeting.coord] = struct{}{}
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
