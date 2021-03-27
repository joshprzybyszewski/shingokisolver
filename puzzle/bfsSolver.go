package puzzle

import (
	"fmt"
)

type bfsQueueItem struct {
	puzzle *puzzle
	coord  nodeCoord
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

	g := newGrid(size, nl)
	var bestCoord nodeCoord
	bestVal := int8(-1)
	for nc, n := range g.nodes {
		if n.val > bestVal {
			bestCoord = nc
			bestVal = n.val
		}
	}

	return &bfsSolver{
		queue: []*bfsQueueItem{{
			puzzle: g,
			coord:  bestCoord,
		}},
	}
}

func (b *bfsSolver) iterations() int {
	return b.numProcessed
}

func (b *bfsSolver) solve() (*puzzle, bool) {
	for len(b.queue) > 0 {
		b.numProcessed++
		g, isSolved := b.processQueueItem()
		if isSolved {
			return g, true
		}
	}
	return nil, false
}

func (b *bfsSolver) processQueueItem() (*puzzle, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]

	if IncludeProgressLogs && (b.numProcessed < 100 || b.numProcessed%10000 == 0) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.row, q.coord.col, b.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.puzzle.IsIncomplete(q.coord); err != nil {
		return q.puzzle, false
	} else if !isIncomplete {
		return q.puzzle, true
	}

	b.addQueueItems(q.puzzle, q.coord)

	return q.puzzle, false
}

func (b *bfsSolver) addQueueItems(
	g *puzzle,
	coord nodeCoord,
) {
	b.addQueueItem(g, headUp, coord)
	b.addQueueItem(g, headRight, coord)
	b.addQueueItem(g, headDown, coord)
	b.addQueueItem(g, headLeft, coord)
}

func (b *bfsSolver) addQueueItem(
	g *puzzle,
	move cardinal,
	coord nodeCoord,
) {
	newCoord, newGrid, err := g.AddEdge(move, coord)
	if err != nil {
		return
	}

	if newGrid.isRangeInvalidWithBoundsCheck(
		newCoord.row-2,
		newCoord.row+2,
		newCoord.col-2,
		newCoord.col+2,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return
	}

	b.queue = append(b.queue, &bfsQueueItem{
		puzzle: newGrid,
		coord:  newCoord,
	})
}
