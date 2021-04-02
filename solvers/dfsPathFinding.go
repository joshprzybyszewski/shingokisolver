package solvers

import (
	"fmt"
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) solveFromLooseEnd(
	input *puzzle.Puzzle,
	start model.NodeCoord,
) (*puzzle.Puzzle, model.State) {

	printAllTargetsHit(fmt.Sprintf(`solveFromLooseEnd(%+v)`, start), input, d.iterations())
	if input == nil {
		return nil, model.Unexpected
	}

	p, state := d.dfsOutFrom(
		input.DeepCopy(),
		start,
	)
	return p, state
}

func (d *targetSolver) dfsOutFrom(
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
) (*puzzle.Puzzle, model.State) {
	if puzz == nil {
		return nil, model.Unexpected
	}

	switch s := puzz.GetState(); s {
	case model.Complete:
		return puzz, s
	case model.Incomplete:
		// continue in the func
	default:
		return nil, s
	}

	d.numProcessed++
	if includeProgressLogs && (d.numProcessed < 100 || d.numProcessed%1000 == 500) {
		log.Printf("dfsOutFrom about to process (%+v): %d\n%s\n",
			fromCoord,
			d.numProcessed,
			puzz.String(),
		)
		fmt.Scanf("hello there")
	}

	for _, nextHeading := range model.AllCardinals {
		nextPuzz := puzz.DeepCopy()

		nextCoord, state := nextPuzz.AddEdge(nextHeading, fromCoord)
		switch state {
		case model.Violation, model.Unexpected, model.Duplicate:
			continue
		}

		if nextPuzz.NumLooseEnds() != puzz.NumLooseEnds() {
			// iterate down
			retPuzz := d.connect(nextPuzz)
			if retPuzz != nil {
				return retPuzz, model.Complete
			}
		} else {
			retPuzz, s := d.dfsOutFrom(
				nextPuzz,
				nextCoord,
			)
			switch s {
			case model.Complete:
				return retPuzz, s
			}
		}
	}

	return nil, model.Incomplete
}
