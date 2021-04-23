package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p Puzzle) GetNodeState(
	nc model.NodeCoord,
) model.State {
	n, ok := p.GetNode(nc)
	if !ok {
		return model.Incomplete
	}

	return getNodeState(n, &p.edges)
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
