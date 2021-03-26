package puzzlegrid

import (
	"errors"
	"fmt"
	"log"
)

type queueItem struct {
	grid *grid
	row  int
	col  int
}

type BfsSolver struct {
	queue []*queueItem
}

func NewSolver(
	size int,
	nl []NodeLocation,
) *BfsSolver {
	if len(nl) == 0 {
		return nil
	}

	return &BfsSolver{
		queue: []*queueItem{{
			grid: newGrid(size, nl),
			row:  nl[0].Row,
			col:  nl[0].Col,
		}},
	}
}

func (b *BfsSolver) Solve() error {
	g, isSolved := b.solve()
	if !isSolved {
		fmt.Println(`could not solve`)
		return errors.New(`unsolvable`)
	}
	// TODO verify if this is what we want?
	fmt.Println(g.String())
	return nil
}

func (b *BfsSolver) solve() (*grid, bool) {
	numProcessed := 0
	for len(b.queue) > 0 {
		g, isSolved := b.processQueueItem()
		if isSolved {
			return g, true
		}
		numProcessed++
		if numProcessed%1600 == 0 {
			fmt.Printf("Processed: %d\n%s\n", numProcessed, g.String())
			fmt.Scanf("hello there")
		}
	}
	return nil, false
}

func (b *BfsSolver) processQueueItem() (*grid, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]
	log.Printf("q.grid:\n%s\n", q.grid.String())

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
