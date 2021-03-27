package puzzlegrid

import (
	"errors"
	"fmt"
	"time"
)

var (
	showProcess = false
)

type queueItem struct {
	grid *grid
	row  int
	col  int
}

type BfsSolver struct {
	queue []*queueItem

	numProcessed int
}

func NewBFSSolver(
	size int,
	nl []NodeLocation,
) *BfsSolver {
	if len(nl) == 0 {
		return nil
	}

	bestR, bestC, bestVal := -1, -1, -1
	for _, n := range nl {
		if n.Value > bestVal {
			bestR = n.Row
			bestC = n.Col
			bestVal = n.Value
		}
	}

	return &BfsSolver{
		queue: []*queueItem{{
			grid: newGrid(size, nl),
			row:  bestR,
			col:  bestC,
		}},
	}
}

func (b *BfsSolver) Solve() error {
	defer func(t0 time.Time) {
		fmt.Printf("BfsSolver.Solve processed %d paths in %s\n", b.numProcessed, time.Since(t0).String())
	}(time.Now())
	g, isSolved := b.solve()
	if !isSolved {
		fmt.Println(`could not solve`)
		return errors.New(`unsolvable`)
	}
	fmt.Println(g.String())
	return nil
}

func (b *BfsSolver) solve() (*grid, bool) {
	for len(b.queue) > 0 {
		b.numProcessed++
		g, isSolved := b.processQueueItem()
		if isSolved {
			fmt.Printf("Processed: %d\n", b.numProcessed)
			return g, true
		}
		if showProcess && (b.numProcessed < 100 || b.numProcessed%100 == 0) {
			fmt.Printf("Processed: %d\n%s\n", b.numProcessed, g.String())
			fmt.Scanf("hello there")
		}
	}
	return nil, false
}

func (b *BfsSolver) processQueueItem() (*grid, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]

	if isIncomplete, err := q.grid.IsIncomplete(); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	b.addQueueItems(q.grid, q.row, q.col)

	return q.grid, false
}

func (b *BfsSolver) addQueueItems(
	g *grid,
	row, col int,
) {
	b.addQueueItem(g, headUp, row, col)
	b.addQueueItem(g, headRight, row, col)
	b.addQueueItem(g, headDown, row, col)
	b.addQueueItem(g, headLeft, row, col)
}

func (b *BfsSolver) addQueueItem(
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

	if newGrid.isRangeInvalid(newRow-2, newRow+2, newCol-2, newCol+2) {
		// this is a sanity check to reduce the amount of calc we need to do
		return
	}

	b.queue = append(b.queue, &queueItem{
		grid: newGrid,
		row:  newRow,
		col:  newCol,
	})
}
