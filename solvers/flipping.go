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
		switch puzz.GetState(model.InvalidNodeCoord) {
		case model.Complete:
			return puzz
		default:
			return nil
		}
	}

	switch puzzWithEdge := puzz.DeepCopy(); puzzWithEdge.AddEdges(ep) {
	case model.Complete, model.Incomplete:
		if puzzWithEdge.GetState(ep.NodeCoord) == model.Complete {
			return puzzWithEdge
		}

		res := d.flip(puzzWithEdge)
		if res != nil {
			return res
		}
	}

	switch puzzWithoutEdge := puzz.DeepCopy(); puzzWithoutEdge.AvoidEdge(ep) {
	case model.Complete, model.Incomplete:
		if puzzWithoutEdge.GetState(ep.NodeCoord) == model.Complete {
			return puzzWithoutEdge
		}
		return d.flip(puzzWithoutEdge)
	}

	return nil
}
