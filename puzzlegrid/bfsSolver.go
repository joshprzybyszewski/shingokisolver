package puzzlegrid

import (
	"fmt"
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

func (b *BfsSolver) Solve() {
	g, isSolved := b.solve()
	if isSolved {
		// TODO verify if this is what we want?
		fmt.Println(g.String())
	} else {
		fmt.Println(`could not solve`)
	}
}

func (b *BfsSolver) solve() (*grid, bool) {
	for len(b.queue) > 0 {
		g, isSolved := b.processQueueItem()
		if isSolved {
			return g, true
		}
	}
	return nil, false
}

func (b *BfsSolver) processQueueItem() (*grid, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]
	// TODO determine this
	// log.Printf("q.grid:\n%s\n", q.grid.String())
	// fmt.Scanf("hello there")

	if isIncomplete, err := q.grid.IsIncomplete(); err != nil {
		return nil, false
	} else if !isIncomplete {
		return q.grid, true
	}

	b.addQueueItems(q.grid, q.row, q.col)

	return nil, false
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
	var dir edgeDirection
	var eIndex, start int

	switch move {
	case headUp:
		dir = colDir
		start = r - 1
		eIndex = c
	case headDown:
		dir = colDir
		start = r
		eIndex = c
	case headLeft:
		dir = rowDir
		start = c - 1
		eIndex = r
	case headRight:
		dir = rowDir
		start = c
		eIndex = r
	}
	newGrid, err := g.AddEdge(dir, eIndex, start)

	if err == nil {
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

		b.queue = append(b.queue, &queueItem{
			grid: newGrid,
			row:  newRow,
			col:  newCol,
		})
	}
}
