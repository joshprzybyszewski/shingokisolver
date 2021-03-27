package puzzle

import (
	"fmt"
)

type dfsSolverStep struct {
	grid  *grid
	coord nodeCoord
}

type dfsSolver struct {
	grid *grid

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
		grid: newGrid(size, nl),
	}
}

func (d *dfsSolver) iterations() int {
	return d.numProcessed
}

func (d *dfsSolver) solve() (*grid, bool) {
	var bestCoord nodeCoord
	bestVal := int8(-1)
	for nc, n := range d.grid.nodes {
		if n.val > bestVal {
			bestCoord = nc
			bestVal = n.val
		}
	}
	return d.takeNextStepIntoDepth(&dfsSolverStep{
		grid:  d.grid,
		coord: bestCoord,
	})
}

func (d *dfsSolver) takeNextStepIntoDepth(
	q *dfsSolverStep,
) (*grid, bool) {
	if q == nil {
		return nil, false
	}

	d.numProcessed++
	if IncludeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 0) {
		fmt.Printf("About to process (%d, %d): %d\n%s\n", q.coord.row, q.coord.col, d.numProcessed, q.grid.String())
		fmt.Scanf("hello there")
	}

	if isIncomplete, err := q.grid.IsIncomplete(q.coord); err != nil {
		return q.grid, false
	} else if !isIncomplete {
		return q.grid, true
	}

	for _, qi := range []*dfsSolverStep{
		d.getNextStep(q.grid, headUp, q.coord),
		d.getNextStep(q.grid, headRight, q.coord),
		d.getNextStep(q.grid, headDown, q.coord),
		d.getNextStep(q.grid, headLeft, q.coord),
	} {
		g, isSolved := d.takeNextStepIntoDepth(qi)
		if isSolved {
			return g, true
		}
	}

	return q.grid, false
}

func (d *dfsSolver) getNextStep(
	g *grid,
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
		grid:  newGrid,
		coord: newCoord,
	}
}
