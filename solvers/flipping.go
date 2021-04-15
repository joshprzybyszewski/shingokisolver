package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) flip(
	puzz *puzzle.Puzzle,
) *puzzle.Puzzle {

	ep, ok := puzz.GetUnknownEdge()
	if !ok {
		switch puzz.GetState() {
		case model.Complete:
			return puzz
		default:
			return nil
		}
	}

	puzzCpy := puzz.DeepCopy()

	switch puzz.AddEdges(ep) {
	case model.Complete, model.Incomplete:
		switch puzzCpy.GetState() {
		case model.Complete:
			return puzzCpy
		}
	}

	switch s := puzz.AvoidEdge(ep); s {
	case model.Complete, model.Incomplete:
		if puzz.GetState() == model.Complete {
			return puzz
		}
	}

	return d.flip(puzz)
}
