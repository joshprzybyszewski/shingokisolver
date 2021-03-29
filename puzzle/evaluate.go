package puzzle

func isInvalidNode(
	p *puzzle,
	nc nodeCoord,
	efn edgesFromNode,
) bool {
	n, ok := p.nodes[nc]
	if !ok || n.nType == noNode {
		// no node == not invalid
		return false
	}

	// check that the node type rules are not broken
	if n.nType.isInvalidEdges(efn) {
		return true
	}

	// check that the num of straight line edges does not exceed the node n.val
	return efn.totalEdges() > n.val
}

func isCompleteNode(
	p *puzzle,
	nc nodeCoord,
) bool {
	oe, ok := p.getOutgoingEdgesFrom(nc)
	if !ok {
		// the coordinate must be out of bounds
		return false
	}

	n, ok := p.nodes[nc]
	if !ok || n.nType == noNode {
		// no node == not invalid
		return false
	}

	// check that the node type rules are not broken
	if n.nType.isInvalidEdges(oe) {
		return false
	}

	// this node needs two outgoing edges and for the sum of the straight lines
	// to be equal to its value
	return oe.getNumOutgoingDirections() == 2 && n.val == oe.totalEdges()
}
