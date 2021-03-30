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
func (partial *partialSolutionItem) removeDuplicateLooseEnds() {
	sort.Slice(partial.looseEnds, func(i, j int) bool {
		if partial.looseEnds[i].Row != partial.looseEnds[j].Row {
			return partial.looseEnds[i].Row < partial.looseEnds[j].Row
		}
		return partial.looseEnds[i].Col < partial.looseEnds[j].Col
	})

	for i := 0; i < len(partial.looseEnds)-1; i++ {
		if partial.looseEnds[i] == partial.looseEnds[i+1] {
			partial.looseEnds = append(
				partial.looseEnds[:i],
				partial.looseEnds[i+2:]...)
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
		partial.targeting == nil {
		return
	}
	if !(iterations < 10 || iterations%100 == 0) {
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

type partialSolutionQueue struct {
	items []*partialSolutionItem
}

func newPartialSolutionQueue() *partialSolutionQueue {
	return &partialSolutionQueue{
		items: make([]*partialSolutionItem, 0),
	}
}

func (q *partialSolutionQueue) isEmpty() bool {
	return len(q.items) == 0
}

func (q *partialSolutionQueue) pop() *partialSolutionItem {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *partialSolutionQueue) push(items ...*partialSolutionItem) {
	// TODO for each item, if it already exists in the queue, don't
	// add it a second time? But how is it non-deterministic? How
	// can we generate the same partial solution?
	q.items = append(q.items, items...)
}
