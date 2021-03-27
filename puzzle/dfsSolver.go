package puzzle

import (
	"fmt"
)

type dfsSolverStep struct {
	puzzle *puzzle
	coord  nodeCoord
}

type dfsSolver struct {
	puzzle *puzzle

	numProcessed int
}

func newDFSSolver(
	size int,
	nl []NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &dfsSolver{
		puzzle: newGrid(size, nl),
	}
}

func (d *dfsSolver) iterations() int {
	return d.numProcessed
}

func (d *dfsSolver) solve() (*puzzle, bool) {
	var bestCoord nodeCoord
	bestVal := int8(-1)
	for nc, n := range d.puzzle.nodes {
		if n.val > bestVal {
			bestCoord = nc
			bestVal = n.val
		}
	}
	return d.takeNextStepIntoDepth(&dfsSolverStep{
		puzzle: d.puzzle,
		coord:  bestCoord,
	})
}

func (d *dfsSolver) takeNextStepIntoDepth(
	q *dfsSolverStep,
) (*puzzle, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 0) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.row, q.coord.col, d.numProcessed, q.puzzle.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.puzzle.IsIncomplete(q.coord); err != nil {
		return q.puzzle, false
	} else if !isIncomplete {
		return q.puzzle, true
	}

	for _, qi := range []*dfsSolverStep{
		d.getNextStep(q.puzzle, headUp, q.coord),
		d.getNextStep(q.puzzle, headRight, q.coord),
		d.getNextStep(q.puzzle, headDown, q.coord),
		d.getNextStep(q.puzzle, headLeft, q.coord),
	} {
		g, isSolved := d.takeNextStepIntoDepth(qi)
		if isSolved {
			return g, true
		}
	}

	return q.puzzle, false
}

func (d *dfsSolver) getNextStep(
	g *puzzle,
	move cardinal,
	nc nodeCoord,
) *dfsSolverStep {
	newCoord, newGrid, err := g.AddEdge(move, nc)
	if err != nil {
		return nil
	}

	if newGrid.isRangeInvalidWithBoundsCheck(
		newCoord.row-2,
		newCoord.row+2,
		newCoord.col-2,
		newCoord.col+2,
	) {
		// this is a sanity check to reduce the amount of calc we need to do
		return nil
	}

	return &dfsSolverStep{
		puzzle: newGrid,
		coord:  newCoord,
	}
}
