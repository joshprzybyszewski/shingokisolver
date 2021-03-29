package solvers

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type bfsQueueItem struct {
	puzzle *puzzle.Puzzle
	coord  model.NodeCoord
}

type bfsSolver struct {
	queue []*bfsQueueItem

	numProcessed int
}

func newBFSSolver(
	size int,
	nl []model.NodeLocation,
) *bfsSolver {
	if len(nl) == 0 {
		return nil
	}

	p := puzzle.NewPuzzle(size, nl)
	bestCoord := p.GetCoordForHighestValueNode()

	return &bfsSolver{
		queue: []*bfsQueueItem{{
			puzzle: p,
			coord:  bestCoord,
		}},
	}
}

func (b *bfsSolver) iterations() int {
	return b.numProcessed
}

func (b *bfsSolver) solve() (*puzzle.Puzzle, bool) {
	for len(b.queue) > 0 {
		b.numProcessed++
		g, isSolved := b.processQueueItem()
		if isSolved {
			return g, true
		}
	}
	return nil, false
}

func (b *bfsSolver) processQueueItem() (*puzzle.Puzzle, bool) {
	q := b.queue[0]
	b.queue = b.queue[1:]

	if includeProgressLogs && (b.numProcessed < 100 || b.numProcessed%10000 == 0) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.Row, q.coord.Col, b.numProcessed, q.puzzle.String())
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
	g *puzzle.Puzzle,
	coord model.NodeCoord,
) {
	b.addQueueItem(g, model.HeadUp, coord)
	b.addQueueItem(g, model.HeadRight, coord)
	b.addQueueItem(g, model.HeadDown, coord)
	b.addQueueItem(g, model.HeadLeft, coord)
}

func (b *bfsSolver) addQueueItem(
	g *puzzle.Puzzle,
	move model.Cardinal,
	coord model.NodeCoord,
) {
	newCoord, np, err := g.AddEdge(move, coord)
	if err != nil {
		return
	}

	if np.IsRangeInvalid(
		newCoord.Row-2,
		newCoord.Row+2,
		newCoord.Col-2,
		newCoord.Col+2,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return
	}

	b.queue = append(b.queue, &bfsQueueItem{
		puzzle: np,
		coord:  newCoord,
	})
}
