package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) IsCompleteNode(
	nc model.NodeCoord,
) bool {
	if p == nil {
		return false
	}

	n, ok := p.nodes[nc]
	if !ok {
		return false
	}

	return n.Value() == p.GetSumOutgoingStraightLines(n.Coord())
}
