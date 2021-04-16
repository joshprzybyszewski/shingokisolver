package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type walker interface {
	walk() (map[model.NodeCoord]struct{}, model.State)
}

type simpleWalker struct {
	provider   model.GetEdger
	start      model.NodeCoord
	cur        model.NodeCoord
	seen       map[model.NodeCoord]struct{}
	foundStart bool
}

func newWalker(
	ge model.GetEdger,
	start model.NodeCoord,
) walker {
	return &simpleWalker{
		provider: ge,
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

	// TODO I think I can remove this switch because the rules
	// will freak out before we get to this point
	// that means I can stop returning a state from this walker too
	switch nOutgoing := getNumOutgoingDirections(sw.provider, sw.cur); {
	case nOutgoing > 2:
		return model.HeadNowhere, model.Violation
	case nOutgoing < 2:
		return model.HeadNowhere, model.Incomplete
	}

	for dir := range model.AllCardinalsMap {
		if avoid == dir {
			continue
		}

		if sw.provider.GetEdge(model.NewEdgePair(sw.cur, dir)) != model.EdgeExists {
			continue
		}

		sw.markNodesAsSeen(sw.cur)
		sw.cur = sw.cur.Translate(dir)
		return dir.Opposite(), model.Incomplete
	}

	return model.HeadNowhere, model.Incomplete
}

func getNumOutgoingDirections(
	ge model.GetEdger,
	coord model.NodeCoord,
) int8 {
	var total int8

	for _, dir := range model.AllCardinals {
		ep := model.NewEdgePair(coord, dir)
		if ge.GetEdge(ep) == model.EdgeExists {
			total++
		}
	}

	return total
}

func (sw *simpleWalker) markNodesAsSeen(
	nc model.NodeCoord,
) {
	if _, ok := sw.seen[nc]; ok {
		// this means we've seen the starting node before
		sw.foundStart = true
	}
	sw.seen[nc] = struct{}{}
}
