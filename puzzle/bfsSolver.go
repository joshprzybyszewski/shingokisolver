package puzzle

import (
	"fmt"
)

type bfsQueueItem struct {
	grid *grid
	row  int
	col  int
}

type bfsSolver struct {
	queue []*bfsQueueItem

	numProcessed int
}

func newBFSSolver(
	size int,
	nl []NodeLocation,
) *bfsSolver {
	if len(nl) == 0 {
		return nil
	}

	bestR, bestC, bestVal := -1, -1, int8(-1)
	for _, n := range nl {
		if n.Value > bestVal {
			bestR = n.Row
			bestC = n.Col
			bestVal = n.Value
		}
	}

	return &bfsSolver{
		queue: []*bfsQueueItem{{
			grid: newGrid(size, nl),
			row:  bestR,
			col:  bestC,
		}},
	}
}

func (b *bfsSolver) iterations() int {
	return b.numProcessed
}

func (b *bfsSolver) solve() (*grid, bool) {
	for len(b.queue) > 0 {
		b.numProcessed++
		g, isSolved := b.processQueueItem()
		if isSolved {
			return g, true
		}
		if IncludeProgressLogs && (b.numProcessed < 100 || b.numProcessed%10000 == 0) {
			fmt.Printf("Processed: %d\n%s\n", b.numProcessed, g.String())
			fmt.Scanf("hello there")
		}
	}
	return nil, false
}

func (b *bfsSolver) processQueueItem() (*grid, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]

	if isIncomplete, err := q.grid.IsIncomplete(q.row, q.col); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	b.addQueueItems(q.grid, q.row, q.col)

	return q.grid, false
}

func (b *bfsSolver) addQueueItems(
	g *grid,
	row, col int,
) {
	b.addQueueItem(g, headUp, row, col)
	b.addQueueItem(g, headRight, row, col)
	b.addQueueItem(g, headDown, row, col)
	b.addQueueItem(g, headLeft, row, col)
}

func (b *bfsSolver) addQueueItem(
	g *grid,
	move cardinal,
	r, c int,
) {
	newGrid, err := g.AddEdge(move, r, c)
	if err != nil {
		return
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
		return
	}

	b.queue = append(b.queue, &bfsQueueItem{
		grid: newGrid,
		row:  newRow,
		col:  newCol,
	})
}
