package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) GetNodeState(
	nc model.NodeCoord,
) model.State {
	if p == nil {
		return model.Incomplete
	}

	n, ok := p.nodes[nc]
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
