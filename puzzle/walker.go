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
}

func newWalker(
	ge model.GetEdger,
	start model.NodeCoord,
) *simpleWalker {
	return &simpleWalker{
		provider: ge,
		start:    start,
		seen:     state.NewCoordSeen(ge.NumEdges()),
	}
}

func (sw *simpleWalker) walkToTheEndOfThePath() (model.NodeCoord, bool) {
	sw.skipSeen = true
	cur, _, isLoop := sw.walkWithInfo(model.HeadNowhere)
	return cur, isLoop
}

func (sw *simpleWalker) walk() (model.NodeCoord, state.CoordSeener, bool) {
	lastNC, _, isLoop := sw.walkWithInfo(model.HeadNowhere)
	if isLoop {
		return model.InvalidNodeCoord, sw.seen, true
	}
	return lastNC, nil, false
}

func (sw *simpleWalker) walkWithInfo(
	initialMove model.Cardinal,
) (
	model.NodeCoord, model.Cardinal, bool,
) {

	cur := sw.start
	lastMove := initialMove

	move := sw.walkToNextPoint(cur, lastMove)
	if move == model.HeadNowhere {
		// our path all the way around was incomplete
		return cur, model.HeadNowhere, false
	}

	cur = cur.Translate(move)
	lastMove = move

	for cur != sw.start {
		move = sw.walkToNextPoint(cur, move.Opposite())

		if move == model.HeadNowhere {
			// if we can't go anywhere, then we'll break out of the loop
			// because this means the path has a loose end.
			return cur, lastMove, false
		}
		cur = cur.Translate(move)
		lastMove = move
	}

	return cur, lastMove, true
}

func (sw *simpleWalker) walkToNextPoint(
	from model.NodeCoord,
	avoid model.Cardinal,
) model.Cardinal {

	going := getNextEdge(sw.provider, from, avoid)

	if going == model.HeadNowhere {
		return model.HeadNowhere
	}

	if !sw.skipSeen {
		sw.seen.Mark(from)
	}

	return going
}

func getNextEdge(
	ge model.GetEdger,
	nc model.NodeCoord,
	avoid model.Cardinal,
) model.Cardinal {

	for _, dir := range model.AllCardinals {
		if avoid == dir {
			continue
		}

		if ge.IsEdge(model.NewEdgePair(nc, dir)) {
			return dir
		}
	}

	return model.HeadNowhere
}
