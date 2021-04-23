package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) GetState() model.State {
	_, s := p.GetStateOfLoop(model.InvalidNodeCoord)
	return s
}

func (p Puzzle) GetStateOfLoop(
	coord model.NodeCoord,
) (model.EdgePair, model.State) {
	nodeState := p.getStateOfNodes()
	switch nodeState {
	case model.Incomplete, model.NodesComplete:
		// check the state of the loop before returning...
	default:
		return model.InvalidEdgePair, nodeState
	}

	lastUnknown, s := p.getStateOfLoop(model.InvalidNodeCoord)
	if s == model.Complete && nodeState == model.Incomplete {
		return model.InvalidEdgePair, model.Violation
	}
	return lastUnknown, s
}

func (p Puzzle) getStateOfLoop(
	coord model.NodeCoord,
) (model.EdgePair, model.State) {

	if coord == model.InvalidNodeCoord {
		coord = p.getRandomCoord()
	}

	w := newWalker(&p.edges, coord)
	seenNodes, isLoop := w.walk()
	if !isLoop {
		nextUnknown := model.InvalidEdgePair
		for dir := range model.AllCardinalsMap {
			ep := model.NewEdgePair(w.cur, dir)
			if !p.edges.IsDefined(ep) {
				nextUnknown = ep
			}
		}
		return nextUnknown, model.Incomplete
	}

	for _, n := range p.nodes {
		if _, ok := seenNodes[n.Coord()]; !ok {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.InvalidEdgePair, model.Violation
		}
	}

	// it's a loop that has all of the nodes completed!
	return model.InvalidEdgePair, model.Complete
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
