package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

// GetNodeState returns the state of the Node at the given Coord.
// Warning: this method could be expensive!
// If the node is satisfied, then we will walk the outgoing
// edges from it to see if there exists a loop that shouldn't exist.
func (p Puzzle) GetNodeState(
	nc model.NodeCoord,
) model.State {
	n, ok := p.GetNode(nc)
	if !ok {
		// why are you asking about this?
		return model.Unexpected
	}

	ns := getNodeState(n, &p.edges)
	if ns != model.Complete {
		return ns
	}

	// We know the node is complete at this point.
	// Let's pre-emptively check for a loop that shouldn't exist
	_, s := p.getStateOfLoop(n.Coord())
	if s != model.Complete && s != model.Incomplete {
		return model.Violation
	}

	return model.Complete
}

// getStateOfNodes checks all of the nodes' state quickly
// Returns `model.NodesComplete` if all nodes are satisfied.
// Does not check for loops.
func (p Puzzle) getStateOfNodes() model.State {
	for _, n := range p.nodes {
		switch s := getNodeState(n, &p.edges); s {
		case model.Complete:
		default:
			return s
		}
	}

	return model.NodesComplete
}

func getNodeState(
	n model.Node,
	ge model.GetEdger,
) model.State {
	nOut, isMax := getSumOutgoingStraightLines(n.Coord(), ge)
	switch {
	case nOut > n.Value():
		return model.Violation
	case n.Value() == nOut:
		return model.Complete
	case isMax:
		return model.Violation
	default:
		return model.Incomplete
	}
}

func getSumOutgoingStraightLines(
	coord model.NodeCoord,
	ge model.GetEdger,
) (int8, bool) {
	var total int8
	numAvoids := 0

	for _, dir := range model.AllCardinals {
		var myTotal int8

		ep := model.NewEdgePair(coord, dir)
		for ge.IsEdge(ep) {
			ep = ep.Next(dir)
			myTotal++
		}

		if myTotal > 0 && ge.IsAvoided(ep) {
			numAvoids++
		}

		total += myTotal
	}

	return total, numAvoids == 4
}