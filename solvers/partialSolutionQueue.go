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
	for i, le := range partial.looseEnds {
		looseEndsCpy[i] = le
	}
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

func printPartialSolution(
	caller string,
	partial *partialSolutionItem,
	iterations int,
) {

	if !includeProgressLogs {
		return
	}
	if partial.puzzle == nil ||
		partial.targeting != nil {
		return
	}
	if !(iterations < 10 || iterations%10000 == 0) {
		return
	}

	log.Printf("printPartialSolution from %s (%d iterations): (targeting %v) %v",
		caller,
		iterations,
		partial.targeting,
		partial.looseEnds,
	)
	log.Printf(":\n%s\n", partial.puzzle.String())
	fmt.Scanf("hello there")
}
