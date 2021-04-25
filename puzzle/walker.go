package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type simpleWalker struct {
	provider model.GetEdger

	seen     state.CoordSeener
	skipSeen bool

	start model.NodeCoord
	cur   model.NodeCoord
}

func newWalker(
	ge model.GetEdger,
	start model.NodeCoord,
) *simpleWalker {
	return &simpleWalker{
		provider: ge,
		start:    start,
		cur:      start,
		seen:     state.NewCoordSeen(ge.NumEdges()),
	}
}

func (sw *simpleWalker) walkToTheEndOfThePath() (model.NodeCoord, bool) {
	sw.skipSeen = true
	_, isLoop := sw.walk()
	return sw.cur, isLoop
}

func (sw *simpleWalker) walk() (state.CoordSeener, bool) {
	move := sw.walkToNextPoint(model.HeadNowhere)
	if move == model.HeadNowhere {
		// our path all the way around was incomplete
		return nil, false
	}

	for sw.cur != sw.start {
		move = sw.walkToNextPoint(move)

		if move == model.HeadNowhere {
			// if we can't go anywhere, then we'll break out of the loop
			// because this means we don't have a loop.
			return nil, false
		}
	}

	return sw.seen, true
}

func (sw *simpleWalker) walkToNextPoint(
	avoid model.Cardinal,
) model.Cardinal {

	for _, dir := range model.AllCardinals {
		if avoid == dir {
			continue
		}

		if !sw.provider.IsEdge(model.NewEdgePair(sw.cur, dir)) {
			continue
		}

		if !sw.skipSeen {
			sw.seen.Mark(sw.cur)
		}

		sw.cur = sw.cur.Translate(dir)
		return dir.Opposite()
	}

	return model.HeadNowhere
}
