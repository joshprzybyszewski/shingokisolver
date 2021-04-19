package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) GetState(
	coord model.NodeCoord,
) model.State {

	nodeState := p.getStateOfNodes()
	switch nodeState {
	case model.Incomplete, model.Complete:
		// keep going through checks...
	default:
		return nodeState
	}

	if coord == model.InvalidNodeCoord {
		for _, n := range p.nodes {
			// just need a random starting node for the walker
			coord = n.Coord()
			break
		}
	}

	w := newWalker(p.edges, coord)
	seenNodes, isLoop := w.walk()
	if !isLoop {
		return model.Incomplete
	}

	for _, n := range p.nodes {
		if _, ok := seenNodes[n.Coord()]; !ok {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.Violation
		}
	}

	return nodeState
}

func (p Puzzle) getStateOfNodes() model.State {
	// it's cheaper for us to just iterate all of the nodes
	// and check for their validity than it is to check every
	// (r, c) or filtering out to only be in the range
	for _, n := range p.nodes {
		switch s := p.GetNodeState(n.Coord()); s {
		case model.Complete:

		default:
			return s
		}
	}

	return model.Complete
}
