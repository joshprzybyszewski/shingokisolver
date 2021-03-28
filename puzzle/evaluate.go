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
	efn edgesFromNode,
) bool {
	n, ok := p.nodes[nc]
	if !ok || n.nType == noNode {
		// no node == not invalid
		return false
	}

	// check that the node type rules are not broken
	if n.nType.isInvalidEdges(efn) {
		return false
	}

	// check that the num of straight line edges < n.val
	return n.val == efn.totalEdges()
}
