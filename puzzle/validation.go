package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func SetNodesComplete(p *Puzzle) {
	p.areNodesComplete = true
}

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
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(w.cur, dir)
			if !p.edges.IsDefined(ep) {
				nextUnknown = ep
				break
			}
		}
		return nextUnknown, model.Incomplete
	}

	for _, n := range p.nodes {
		if !seenNodes.IsCoordSeen(n.Coord()) {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.InvalidEdgePair, model.Violation
		}
	}

	// it's a loop that has all of the nodes completed!
	return model.InvalidEdgePair, model.Complete
}

func (p Puzzle) getRandomCoord() model.NodeCoord {
	for _, n := range p.nodes {
		return n.Coord()
	}

	if p.rules == nil {
		panic(`dev error: p.rules == nil`)
	}

	panic(`dev error: getRandomCoord couldn't find anything`)
}
