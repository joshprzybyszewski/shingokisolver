package puzzlegrid

import (
	"errors"
	"fmt"
	"time"
)

type DfsSolver struct {
	grid *grid

	numProcessed int
}

func NewDFSSolver(
	size int,
	nl []NodeLocation,
) *DfsSolver {
	if len(nl) == 0 {
		return nil
	}

	return &DfsSolver{
		grid: newGrid(size, nl),
	}
}

func (d *DfsSolver) Solve() error {
	defer func(t0 time.Time) {
		fmt.Printf("DfsSolver.Solve processed %d paths in %s\n", d.numProcessed, time.Since(t0).String())
	}(time.Now())

	g, isSolved := d.solve()
	if !isSolved {
		fmt.Println(`could not solve`)
		return errors.New(`unsolvable`)
	}
	fmt.Println(g.String())
	return nil
}

func (d *DfsSolver) solve() (*grid, bool) {
	bestR, bestC, bestVal := -1, -1, -1
	for r, cMap := range d.grid.nodes {
		for c, n := range cMap {
			if n.val > bestVal {
				bestR = int(r)
				bestC = int(c)
				bestVal = n.val
			}
		}
	}
	return d.processQueueItem(&queueItem{
		grid: d.grid,
		row:  bestR,
		col:  bestC,
	})
}

func (d *DfsSolver) processQueueItem(
	q *queueItem,
) (*grid, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if showProcess && (d.numProcessed < 100 || d.numProcessed%10 == 0) {
		fmt.Printf("Processing (%d, %d): %d\n%s\n", q.row, q.col, d.numProcessed, q.grid.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.grid.IsIncomplete(); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	for _, qi := range []*queueItem{
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

func (d *DfsSolver) getQueueItem(
	g *grid,
	move cardinal,
	r, c int,
) *queueItem {
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

	return &queueItem{
		grid: newGrid,
		row:  newRow,
		col:  newCol,
	}
}
