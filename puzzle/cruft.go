package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) GetUnknownEdge() (model.EdgePair, bool) {
	for r := 0; r <= p.NumEdges(); r++ {
		for c := 0; c <= p.NumEdges(); c++ {
			nc := model.NewCoordFromInts(r, c)

			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}

			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}
		}
	}
	return model.EdgePair{}, false
}

func (p *Puzzle) GetLooseEnd() (model.NodeCoord, model.State) {
	for r := model.RowIndex(0); r <= model.RowIndex(p.numEdges); r++ {
		for c := model.ColIndex(0); c <= model.ColIndex(p.numEdges); c++ {
			nc := model.NewCoord(r, c)

			switch numEdges := getNumOutgoingDirections(p.edges, nc); {
			case numEdges > 2:
				return model.InvalidNodeCoord, model.Violation
			case numEdges == 1:
				return nc, model.Incomplete
			}
		}
	}
	return model.InvalidNodeCoord, model.NodesComplete
}

func (p *Puzzle) HasTwoOutgoingEdges(
	coord model.NodeCoord,
) bool {
	nOut := getNumOutgoingDirections(p.edges, coord)
	p.printMsg(
		"HasTwoOutgoingEdges(%+v) = %d",
		coord,
		nOut,
	)
	return nOut == 2
}

func getNumOutgoingDirections(
	ge model.GetEdger,
	coord model.NodeCoord,
) int8 {
	var total int8

	for _, dir := range model.AllCardinals {
		ep := model.NewEdgePair(coord, dir)
		if ge.GetEdge(ep) == model.EdgeExists {
			total++
		}
	}

	return total
}

func (p *Puzzle) GetSumOutgoingStraightLines(
	coord model.NodeCoord,
) (int8, bool) {
	return getSumOutgoingStraightLines(p.edges, coord)
}

func getSumOutgoingStraightLines(
	ge model.GetEdger,
	coord model.NodeCoord,
) (int8, bool) {
	var total int8
	numAvoids := 0

	for _, dir := range model.AllCardinals {
		c := coord
		ep := model.NewEdgePair(c, dir)
		for ge.GetEdge(ep) == model.EdgeExists {
			total++
			c = c.Translate(dir)
			ep = model.NewEdgePair(c, dir)
		}
		if c != coord {
			switch ge.GetEdge(ep) {
			case model.EdgeAvoided, model.EdgeOutOfBounds:
				numAvoids++
			}
		}
	}

	return total, numAvoids >= 2
}
