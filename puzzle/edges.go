package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) GetUnknownEdge() (model.EdgePair, bool) {

	// start from a point that is guaranteed to have an edge.
	// so we choose a node!
	nc := p.getRandomCoord()

	// now let's walk to the end of the line
	w := newWalker(&p.edges, nc)
	nc, isLoop := w.walkToTheEndOfThePath()
	if isLoop {
		// This is an error case. We made a loop, but we weren't expecting to.
		return model.EdgePair{}, false
	}

	// now from this end of the path, choose a random edge
	// off of it that is unknown.
	for dir := range model.AllCardinalsMap {
		ep := model.NewEdgePair(nc, dir)
		if !p.edges.IsDefined(ep) {
			return ep, true
		}
	}

	// we walked to the end of the path and did not find an unknown edge.
	// this is an error case.
	return model.EdgePair{}, false
}

func (p Puzzle) IsEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	return p.edges.GetEdge(model.NewEdgePair(nc, move)) == model.EdgeExists
}
