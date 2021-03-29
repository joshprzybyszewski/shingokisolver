package solvers

import (
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type headsAndTailsItem struct {
	puzzle  *puzzle.Puzzle
	head    model.NodeCoord
	tail    model.NodeCoord
	useHead bool
}

type headsAndTailsDFSSolver struct {
	puzzle *puzzle.Puzzle

	numProcessed int
}

func newHeadsAndTailsDFSSolver(
	size int,
	nl []model.NodeLocation,
) *headsAndTailsDFSSolver {
	if len(nl) == 0 {
		return nil
	}

	return &headsAndTailsDFSSolver{
		puzzle: puzzle.NewPuzzle(size, nl),
	}
}

func (d *headsAndTailsDFSSolver) iterations() int {
	return d.numProcessed
}

func (d *headsAndTailsDFSSolver) solve() (*puzzle.Puzzle, bool) {
	bestCoord := d.puzzle.GetCoordForHighestValueNode()
	return d.processQueueItem(&headsAndTailsItem{
		puzzle:  d.puzzle,
		head:    bestCoord,
		tail:    bestCoord,
		useHead: true,
	})
}

func (d *headsAndTailsDFSSolver) processQueueItem(
	q *headsAndTailsItem,
) (*puzzle.Puzzle, bool) {
	if q == nil {
		return nil, false
	}
	nodeToUse := q.tail
	if q.useHead {
		nodeToUse = q.head
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%10000 == 0) {
		fmt.Printf("About to process head(%d, %d) tail(%d, %d): %d\n%s\n",
			q.head.Row, q.head.Col,
			q.tail.Row, q.tail.Col,
			d.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.puzzle.IsIncomplete(nodeToUse); err != nil {
		return q.puzzle, false
	} else if !isIncomplete {
		return q.puzzle, true
	}

	for c, useHead := range map[model.Cardinal]bool{
		model.HeadUp:    q.useHead,
		model.HeadRight: q.useHead,
		model.HeadDown:  q.useHead,
		model.HeadLeft:  q.useHead,
	} {
		qi, err := d.getQueueItem(q.puzzle, c, q, useHead)
		if err != nil {
			continue
			// return q.puzzle, false
		}
		g, isSolved := d.processQueueItem(qi)
		if isSolved {
			return g, true
		}
	}

	// for c, useHead := range map[model.Cardinal]bool{
	// 	model.HeadUp:    !q.useHead,
	// 	model.HeadRight: !q.useHead,
	// 	model.HeadDown:  !q.useHead,
	// 	model.HeadLeft:  !q.useHead,
	// } {
	// 	qi, err := d.getQueueItem(q.puzzle, c, q, useHead)
	// 	if err != nil {
	// 		continue
	// 		// return q.puzzle, false
	// 	}
	// 	g, isSolved := d.processQueueItem(qi)
	// 	if isSolved {
	// 		return g, true
	// 	}
	// }

	return q.puzzle, false
}

func (d *headsAndTailsDFSSolver) getQueueItem(
	g *puzzle.Puzzle,
	move model.Cardinal,
	q *headsAndTailsItem,
	useHead bool,
) (*headsAndTailsItem, error) {

	nodeToUse := q.tail
	if useHead {
		nodeToUse = q.head
	}

	newCoord, NewPuzzle, err := g.AddEdge(move, nodeToUse)
	if err != nil {
		return nil, err
	}

	if NewPuzzle.IsRangeInvalid(
		newCoord.Row-2,
		newCoord.Row+2,
		newCoord.Col-2,
		newCoord.Col+2,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil, errors.New(`invalid local range`)
	}

	newHead := q.head
	newTail := q.tail
	if useHead {
		newHead = newCoord
	} else {
		newTail = newCoord
	}

	return &headsAndTailsItem{
		puzzle:  NewPuzzle,
		head:    newHead,
		tail:    newTail,
		useHead: !useHead,
	}, nil
}
