package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (d *targetSolver) solveFromLooseEnd(
	input *puzzle.Puzzle,
	start model.NodeCoord,
) *puzzle.Puzzle {

	if input == nil {
		return nil
	}

	otherPuzz, state := d.sendOutDFSPath(
		input.DeepCopy(),
		start,
		model.HeadNowhere,
	)
	switch state {
	case model.NodesComplete:
		if otherPuzz.GetState() == model.Complete {
			return otherPuzz
		}
		return nil
	case model.Incomplete:
		return d.connect(otherPuzz.DeepCopy())
	default:
		return nil
	}
}

func (d *targetSolver) sendOutDFSPath(
	puzz *puzzle.Puzzle,
	fromCoord model.NodeCoord,
	avoid model.Cardinal,
	curPath ...model.EdgePair,
) (*puzzle.Puzzle, model.State) {

	for _, nextHeading := range model.AllCardinals {
		if nextHeading == avoid {
			continue
		}

		ep := model.NewEdgePair(fromCoord, nextHeading)
		switch puzz.GetEdgeState(ep) {
		case model.EdgeUnknown:
			// keep going
		case model.EdgeExists:
			cur := puzz.DeepCopy()

			switch state := cur.AddEdges(append(curPath, ep)...); state {
			case model.Complete, model.Incomplete:
				return cur, state
			}

			continue
		default:
			continue
		}

		if ep.IsIn(curPath...) {
			continue
		}

		otherPuzz, state := d.sendOutDFSPath(
			puzz.DeepCopy(),
			fromCoord.Translate(nextHeading),
			nextHeading.Opposite(),
			append(curPath, ep)...,
		)
		switch state {
		case model.Complete, model.Incomplete:
			return otherPuzz, state
		default:
			return nil, state
		}
	}

	return nil, model.Violation
}
