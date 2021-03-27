package puzzle

import (
	"fmt"
)

type dfsQueueItem struct {
	grid *grid
	row  int
	col  int
}

type dfsSolver struct {
	grid *grid

	numProcessed int
}

func newDFSSolver(
	size int,
	nl []NodeLocation,
) *dfsSolver {
	if len(nl) == 0 {
		return nil
	}

	return &dfsSolver{
		grid: newGrid(size, nl),
	}
}

func (d *dfsSolver) iterations() int {
	return d.numProcessed
}

func (d *dfsSolver) solve() (*grid, bool) {
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
	return d.processQueueItem(&dfsQueueItem{
		grid: d.grid,
		row:  bestR,
		col:  bestC,
	})
}

func (d *dfsSolver) processQueueItem(
	q *dfsQueueItem,
) (*grid, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 0) {
		fmt.Printf("Processing (%d, %d): %d\n%s\n", q.row, q.col, d.numProcessed, q.grid.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.grid.IsIncomplete(q.row, q.col); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	for _, qi := range []*dfsQueueItem{
		d.getQueueItem(q.grid, headUp, q.row, q.col),
		d.getQueueItem(q.grid, headRight, q.row, q.col),
		d.getQueueItem(q.grid, headDown, q.row, q.col),
		d.getQueueItem(q.grid, headLeft, q.row, q.col),
	} {
		g, isSolved := d.processQueueItem(qi)
		if isSolved {
			return g, true
		}
	}

	return q.grid, false
}

func (d *dfsSolver) getQueueItem(
	g *grid,
	move cardinal,
	r, c int,
) *dfsQueueItem {
	newGrid, err := g.AddEdge(move, r, c)
	if err != nil {
		return nil
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

	if newGrid.isRangeInvalid(newRow-2, newRow+2, newCol-2, newCol+2) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil
	}

	return &dfsQueueItem{
		grid: newGrid,
		row:  newRow,
		col:  newCol,
	}
}
