package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) GetNodeState(
	nc model.NodeCoord,
) model.State {
	if p == nil {
		return model.Incomplete
	}

	n, ok := p.getNode(nc)
	if !ok {
		return model.Incomplete
	}

	nOut, isMax := p.GetSumOutgoingStraightLines(n.Coord())
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

func (p *Puzzle) GetSumOutgoingStraightLines(
	coord model.NodeCoord,
) (int8, bool) {
	var total int8
	numAvoids := 0

	for _, dir := range model.AllCardinals {
		c := coord

		ep := model.NewEdgePair(c, dir)
		for p.edges.IsEdge(ep) {
			total++
			c = c.Translate(dir)
			ep = model.NewEdgePair(c, dir)
		}

		if c != coord && p.edges.IsAvoided(ep) {
			numAvoids++
		}
	}

	return total, numAvoids >= 2
}
