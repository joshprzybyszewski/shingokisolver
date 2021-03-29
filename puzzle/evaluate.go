package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func isInvalidNode(
	p *Puzzle,
	nc model.NodeCoord,
	oe model.OutgoingEdges,
) bool {
	n, ok := p.nodes[nc]
	return ok && n.IsInvalid(oe)
}

func IsCompleteNode(
	p *Puzzle,
	nc model.NodeCoord,
) bool {
	oe, ok := p.getOutgoingEdgesFrom(nc)
	if !ok {
		// the coordinate must be out of bounds
		return false
	}

	n, ok := p.nodes[nc]
	return ok && n.IsComplete(oe)
}
