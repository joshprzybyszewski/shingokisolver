package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type getEdgesFromNoder interface {
	GetOutgoingEdgesFrom(model.NodeCoord) (model.OutgoingEdges, bool)
}

type walker interface {
	walk() (map[model.NodeCoord]struct{}, bool)
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

func (sw *simpleWalker) walk() (map[model.NodeCoord]struct{}, bool) {
	move := sw.walkToNextPoint(model.HeadNowhere)
	if move == model.HeadNowhere {
		// our path all the way around was incomplete
		return nil, false
	}

	for sw.cur.Row != sw.start.Row || sw.cur.Col != sw.start.Col {
		move = sw.walkToNextPoint(move)
		if move == model.HeadNowhere || sw.foundStart {
			// if we can't go anywhere, then we'll break out of the loop
			// because this means we don't have a loop.
			break
		}
	}

	if sw.foundStart || sw.cur.Row == sw.start.Row && sw.cur.Col == sw.start.Col {
		return sw.seen, true
	}

	return nil, false
}

func (sw *simpleWalker) walkToNextPoint(
	avoid model.Cardinal,
) model.Cardinal {

	efn, ok := sw.provider.GetOutgoingEdgesFrom(sw.cur)
	if !ok {
		return model.HeadNowhere
	}

	if efn.IsAbove() && avoid != model.HeadUp {
		nextRow := sw.cur.Row - model.RowIndex(efn.Above())
		sw.markNodesAsSeen(nextRow, sw.cur.Row, sw.cur.Col, sw.cur.Col)
		sw.cur.Row = nextRow
		return model.HeadDown
	}

	if efn.IsLeft() && avoid != model.HeadLeft {
		nextCol := sw.cur.Col - model.ColIndex(efn.Left())
		sw.markNodesAsSeen(sw.cur.Row, sw.cur.Row, nextCol, sw.cur.Col)
		sw.cur.Col = nextCol
		return model.HeadRight
	}

	if efn.IsBelow() && avoid != model.HeadDown {
		nextRow := sw.cur.Row + model.RowIndex(efn.Below())
		sw.markNodesAsSeen(sw.cur.Row, nextRow, sw.cur.Col, sw.cur.Col)
		sw.cur.Row = nextRow
		return model.HeadUp
	}

	if efn.IsRight() && avoid != model.HeadRight {
		nextCol := sw.cur.Col + model.ColIndex(efn.Right())
		sw.markNodesAsSeen(sw.cur.Row, sw.cur.Row, sw.cur.Col, nextCol)
		sw.cur.Col = nextCol
		return model.HeadLeft
	}

	return model.HeadNowhere
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
