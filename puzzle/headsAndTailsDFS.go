package puzzle

import (
	"errors"
	"fmt"
)

type headsAndTailsNode struct {
	row int
	col int
}

type headsAndTailsItem struct {
	grid    *grid
	head    headsAndTailsNode
	tail    headsAndTailsNode
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
	bestR, bestC, bestVal := -1, -1, int8(-1)
	for r, cMap := range d.grid.nodes {
		for c, n := range cMap {
			if n.val > bestVal {
				bestR = int(r)
				bestC = int(c)
				bestVal = n.val
			}
		}
	}
	return d.processQueueItem(&headsAndTailsItem{
		grid: d.grid,
		head: headsAndTailsNode{
			row: bestR,
			col: bestC,
		},
		tail: headsAndTailsNode{
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
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 0) {
		fmt.Printf("Processing (%d, %d): %d\n%s\n", nodeToUse.row, nodeToUse.col, d.numProcessed, q.grid.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.grid.IsIncomplete(nodeToUse.row, nodeToUse.col); err != nil {
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
	r := nodeToUse.row
	c := nodeToUse.col

	newGrid, err := g.AddEdge(move, r, c)
	if err != nil {
		return nil, err
	}

	var newRow, newCol int
	switch move {
	case headUp:
		newRow = r - 1
		newCol = c
	case headDown:
		newRow = r + 1
		newCol = c
	case headLeft:
		newRow = r
		newCol = c - 1
	case headRight:
		newRow = r
		newCol = c + 1
	}

	if newGrid.isRangeInvalidWithBoundsCheck(newRow-2, newRow+2, newCol-2, newCol+2) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil, errors.New(`invalid local range`)
	}

	newHead := q.head
	newTail := q.tail
	if useHead {
		newHead.row = newRow
		newHead.col = newCol
	} else {
		newTail.row = newRow
		newTail.col = newCol
	}

	return &headsAndTailsItem{
		grid:    newGrid,
		head:    newHead,
		tail:    newTail,
		useHead: !useHead,
	}, nil
}
