package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func SetNodesComplete(p *Puzzle) {
	p.loop = newAllSegmentsFromNodesComplete(p.nodes, &p.edges)
}

func (p Puzzle) GetState() model.State {
	_, s := p.GetStateOfLoop()
	return s
}

func (p Puzzle) GetStateOfLoop() (model.EdgePair, model.State) {
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

	if p.loop != nil {
		if p.loop.IsLoop() {
			if p.loop.NumNodesInLoop() != len(p.nodes) {
				return model.InvalidEdgePair, model.Violation
			}
			return model.InvalidEdgePair, model.Complete
		}

		return p.loop.GetUnknownEdge(&p.edges)
	}

	if coord == model.InvalidNodeCoord {
		coord = p.getRandomCoord()
	}

	w := newWalker(&p.edges, coord)
	lastNC, seenNodes, isLoop := w.walk()
	if !isLoop {
		nextUnknown := model.InvalidEdgePair
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(lastNC, dir)
			if !p.edges.IsDefined(ep) {
				nextUnknown = ep
				break
			}
		}
		return nextUnknown, model.Incomplete
	}

	for _, nm := range p.nodes {
		if !seenNodes.IsCoordSeen(nm.node.Coord()) {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.InvalidEdgePair, model.Violation
		}
	}

	// it's a loop that has all of the nodes completed!
	return model.InvalidEdgePair, model.Complete
}

func (p Puzzle) getRandomCoord() model.NodeCoord {
	return getRandomCoord(p.nodes)
}

func getRandomCoord(metas []*nodeMeta) model.NodeCoord {
	for _, nm := range metas {
		return nm.node.Coord()
	}

	panic(`dev error: getRandomCoord couldn't find anything`)
}
