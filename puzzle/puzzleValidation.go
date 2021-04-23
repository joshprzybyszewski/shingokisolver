package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) GetState(
	coord model.NodeCoord,
) model.State {

	nodeState := p.getStateOfNodes()
	switch nodeState {
	case model.Incomplete, model.NodesComplete:
		// keep going through checks...
	default:
		return nodeState
	}

	if coord == model.InvalidNodeCoord {
		coord = p.getRandomCoord()
	}

	w := newWalker(&p.edges, coord)
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

	if nodeState == model.NodesComplete {
		// it's a loop that has all of the nodes completed!
		return model.Complete
	}

	// it's a loop that didn't complete all of the nodes!
	return model.Violation
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

	return model.NodesComplete
}

func (p Puzzle) getRandomCoord() model.NodeCoord {
	for _, n := range p.nodes {
		return n.Coord()
	}

	for _, nwo := range p.twoArmOptions {
		return nwo.Coord()
	}

	if p.rules == nil {
		panic(`dev error: p.rules == nil`)
	}

	panic(`dev error: getRandomCoord couldn't find anything`)
}
