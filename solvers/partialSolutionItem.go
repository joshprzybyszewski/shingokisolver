package solvers

import (
	"fmt"
	"log"
	"sort"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type partialSolutionItem struct {
	puzzle *puzzle.Puzzle

	targeting *target

	looseEnds []model.NodeCoord
}

// eliminates loose ends that don't actually exist
// leaves the looseEnds slice in the order that it had previously
func (partial *partialSolutionItem) removeDuplicateLooseEnds() {
	looseEndsCpy := make([]model.NodeCoord, len(partial.looseEnds))
	copy(looseEndsCpy, partial.looseEnds)
	sort.Slice(looseEndsCpy, func(i, j int) bool {
		if looseEndsCpy[i].Row != looseEndsCpy[j].Row {
			return looseEndsCpy[i].Row < looseEndsCpy[j].Row
		}
		return looseEndsCpy[i].Col < looseEndsCpy[j].Col
	})

	shouldRemove := make(map[model.NodeCoord]struct{})
	for i := 0; i < len(looseEndsCpy)-1; i++ {
		if looseEndsCpy[i] == looseEndsCpy[i+1] {
			shouldRemove[looseEndsCpy[i]] = struct{}{}
		}
	}

	for i := 0; i < len(partial.looseEnds); i++ {
		if _, ok := shouldRemove[partial.looseEnds[i]]; ok {
			partial.looseEnds = append(
				partial.looseEnds[:i],
				partial.looseEnds[i+1:]...,
			)
			i--
		}
	}
}

var (
	seen = map[model.NodeCoord]struct{}{}
)

func printPartialSolution(
	caller string,
	partial *partialSolutionItem,
	iterations int,
) {

	if !includeProgressLogs {
		return
	}
	shouldSkip := false
	if partial.puzzle == nil {
		shouldSkip = true
	}
	if partial.targeting != nil {
		shouldSkip = true
		if _, ok := seen[partial.targeting.coord]; ok {
			shouldSkip = false
			seen[partial.targeting.coord] = struct{}{}
		}
	}
	if iterations < 10 || iterations%10000 == 0 {
		shouldSkip = false
	}
	if shouldSkip {
		return
	}

	log.Printf("printPartialSolution from %s (%d iterations): (targeting %v) looseEnds: %v",
		caller,
		iterations,
		partial.targeting,
		partial.looseEnds,
	)
	log.Printf(":\n%s\n", partial.puzzle)
	fmt.Scanf("wait for acknowledgement")
}
