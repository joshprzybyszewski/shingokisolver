package puzzle

import (
	"errors"
	"fmt"
)

type headsAndTailsItem struct {
	grid    *grid
	head    nodeCoord
	tail    nodeCoord
	useHead bool
}

type headsAndTailsDFSSolver struct {
	grid *grid

	numProcessed int
}

func newHeadsAndTailsDFSSolver(
	size int,
	nl []NodeLocation,
) *headsAndTailsDFSSolver {
	if len(nl) == 0 {
		return nil
	}

	return &headsAndTailsDFSSolver{
		grid: newGrid(size, nl),
	}
}

func (d *headsAndTailsDFSSolver) iterations() int {
	return d.numProcessed
}

func (d *headsAndTailsDFSSolver) solve() (*grid, bool) {
	bestR, bestC, bestVal := rowIndex(-1), colIndex(-1), int8(-1)
	for nc, n := range d.grid.nodes {
		if n.val > bestVal {
			bestR = nc.row
			bestC = nc.col
			bestVal = n.val
		}
	}
	return d.processQueueItem(&headsAndTailsItem{
		grid: d.grid,
		head: nodeCoord{
			row: bestR,
			col: bestC,
		},
		tail: nodeCoord{
			row: bestR,
			col: bestC,
		},
		useHead: true,
	})
}

func (d *headsAndTailsDFSSolver) processQueueItem(
	q *headsAndTailsItem,
) (*grid, bool) {
	if q == nil {
		return nil, false
	}
	nodeToUse := q.tail
	if q.useHead {
		nodeToUse = q.head
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%10000 == 0) {
		fmt.Printf("About to process head(%d, %d) tail(%d, %d): %d\n%s\n",
			q.head.row, q.head.col,
			q.tail.row, q.tail.col,
			d.numProcessed, q.grid.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.grid.IsIncomplete(nodeToUse); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	for c, useHead := range map[cardinal]bool{
		headUp:    q.useHead,
		headRight: q.useHead,
		headDown:  q.useHead,
		headLeft:  q.useHead,
	} {
		qi, err := d.getQueueItem(q.grid, c, q, useHead)
		if err != nil {
			continue
			// return q.grid, false
		}
		g, isSolved := d.processQueueItem(qi)
		if isSolved {
			return g, true
		}
	}

	for c, useHead := range map[cardinal]bool{
		headUp:    !q.useHead,
		headRight: !q.useHead,
		headDown:  !q.useHead,
		headLeft:  !q.useHead,
	} {
		qi, err := d.getQueueItem(q.grid, c, q, useHead)
		if err != nil {
			continue
			// return q.grid, false
		}
		g, isSolved := d.processQueueItem(qi)
		if isSolved {
			return g, true
		}
	}

	return q.grid, false
}

func (d *headsAndTailsDFSSolver) getQueueItem(
	g *grid,
	move cardinal,
	q *headsAndTailsItem,
	useHead bool,
) (*headsAndTailsItem, error) {

	nodeToUse := q.tail
	if useHead {
		nodeToUse = q.head
	}

	newCoord, newGrid, err := g.AddEdge(move, nodeToUse)
	if err != nil {
		return nil, err
	}

	if newGrid.isRangeInvalidWithBoundsCheck(
		newCoord.row-2,
		newCoord.row+2,
		newCoord.col-2,
		newCoord.col+2,
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
		grid:    newGrid,
		head:    newHead,
		tail:    newTail,
		useHead: !useHead,
	}, nil
}
