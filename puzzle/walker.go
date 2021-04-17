package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type simpleWalker struct {
	provider model.GetEdger
	start    model.NodeCoord
	cur      model.NodeCoord

	skipSeen bool
	seen     map[model.NodeCoord]struct{}
}

func newWalker(
	ge model.GetEdger,
	start model.NodeCoord,
) *simpleWalker {
	return &simpleWalker{
		provider: ge,
		start:    start,
		cur:      start,
		seen:     make(map[model.NodeCoord]struct{}, 16),
	}
}
func (sw *simpleWalker) walkToTheEndOfThePath() (model.NodeCoord, bool) {
	sw.skipSeen = true
	_, isLoop := sw.walk()
	return sw.cur, isLoop
}

func (sw *simpleWalker) walk() (map[model.NodeCoord]struct{}, bool) {
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

	for dir := range model.AllCardinalsMap {
		if avoid == dir {
			continue
		}

		if !sw.provider.IsEdge(model.NewEdgePair(sw.cur, dir)) {
			continue
		}

		if !sw.skipSeen {
			// avoid allocs if skipSeen is true
			sw.seen[sw.cur] = struct{}{}
		}

		sw.cur = sw.cur.Translate(dir)
		return dir.Opposite()
	}

	return model.HeadNowhere
}
