package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) GetState() model.State {
	if p == nil {
		return model.Violation
	}

	nodeState := p.getStateOfNodes()
	switch nodeState {
	case model.Incomplete, model.Complete:
		// keep going through checks...
	default:
		return nodeState
	}

	var coord model.NodeCoord
	for nc := range p.nodes {
		// just need a random starting node for the walker
		coord = nc
		break
	}

	w := newWalker(p.edges, coord)
	seenNodes, walkerState := w.walk()
	switch walkerState {
	case model.Complete:
		// keep going through checks...
	default:
		return walkerState
	}

	for nc := range p.nodes {
		if _, ok := seenNodes[nc]; !ok {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.Violation
		}
	}

	return nodeState
}

func (p *Puzzle) getStateOfNodes() model.State {
	// it's cheaper for us to just iterate all of the nodes
	// and check for their validity than it is to check every
	// (r, c) or filtering out to only be in the range
	for nc := range p.nodes {
		switch s := p.GetNodeState(nc); s {
		case model.Complete:

		default:
			return s
		}
	}

	return model.Complete
}
