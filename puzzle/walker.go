package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type getEdgesFromNoder interface {
	GetOutgoingEdgesFrom(model.NodeCoord) (model.OutgoingEdges, bool)
}

type walker interface {
	walk() (map[model.NodeCoord]struct{}, model.State)
}

type simpleWalker struct {
	provider   getEdgesFromNoder
	start      model.NodeCoord
	cur        model.NodeCoord
	seen       map[model.NodeCoord]struct{}
	foundStart bool
}

func newWalker(
	getEFNer getEdgesFromNoder,
	start model.NodeCoord,
) walker {
	return &simpleWalker{
		provider: getEFNer,
		start:    start,
		cur:      start,
		seen:     map[model.NodeCoord]struct{}{},
	}
}

func (sw *simpleWalker) walk() (map[model.NodeCoord]struct{}, model.State) {
	move, state := sw.walkToNextPoint(model.HeadNowhere)
	if move == model.HeadNowhere || state != model.Incomplete {
		// our path all the way around was incomplete
		return nil, state
	}

	for sw.cur.Row != sw.start.Row || sw.cur.Col != sw.start.Col {
		move, state = sw.walkToNextPoint(move)

		switch state {
		case model.Complete, model.Incomplete:
		default:
			return nil, state
		}

		if move == model.HeadNowhere || sw.foundStart {
			// if we can't go anywhere, then we'll break out of the loop
			// because this means we don't have a loop.
			break
		}
	}

	if sw.foundStart || sw.cur.Row == sw.start.Row && sw.cur.Col == sw.start.Col {
		return sw.seen, model.Complete
	}

	return nil, model.Incomplete
}

func (sw *simpleWalker) walkToNextPoint(
	avoid model.Cardinal,
) (model.Cardinal, model.State) {

	oe, ok := sw.provider.GetOutgoingEdgesFrom(sw.cur)
	if !ok {
		return model.HeadNowhere, model.Unexpected
	}

	switch nOutgoing := oe.GetNumOutgoingDirections(); {
	case nOutgoing > 2:
		return model.HeadNowhere, model.Violation
	case nOutgoing < 2:
		return model.HeadNowhere, model.Incomplete
	}

	if oe.IsAbove() && avoid != model.HeadUp {
		nextRow := sw.cur.Row - model.RowIndex(oe.Above())
		sw.markNodesAsSeen(nextRow, sw.cur.Row, sw.cur.Col, sw.cur.Col)
		sw.cur.Row = nextRow
		return model.HeadDown, model.Incomplete
	}

	if oe.IsLeft() && avoid != model.HeadLeft {
		nextCol := sw.cur.Col - model.ColIndex(oe.Left())
		sw.markNodesAsSeen(sw.cur.Row, sw.cur.Row, nextCol, sw.cur.Col)
		sw.cur.Col = nextCol
		return model.HeadRight, model.Incomplete
	}

	if oe.IsBelow() && avoid != model.HeadDown {
		nextRow := sw.cur.Row + model.RowIndex(oe.Below())
		sw.markNodesAsSeen(sw.cur.Row, nextRow, sw.cur.Col, sw.cur.Col)
		sw.cur.Row = nextRow
		return model.HeadUp, model.Incomplete
	}

	if oe.IsRight() && avoid != model.HeadRight {
		nextCol := sw.cur.Col + model.ColIndex(oe.Right())
		sw.markNodesAsSeen(sw.cur.Row, sw.cur.Row, sw.cur.Col, nextCol)
		sw.cur.Col = nextCol
		return model.HeadLeft, model.Incomplete
	}

	return model.HeadNowhere, model.Unexpected
}

func (sw *simpleWalker) markNodesAsSeen(
	minR, maxR model.RowIndex,
	minC, maxC model.ColIndex,
) {
	for r := minR; r <= maxR; r++ {
		for c := minC; c <= maxC; c++ {
			nc := model.NewCoord(r, c)
			if nc == sw.start {
				if _, ok := sw.seen[nc]; ok {
					// this means we've seen the starting node before
					sw.foundStart = true
				}
			}
			sw.seen[nc] = struct{}{}
		}
	}
}
