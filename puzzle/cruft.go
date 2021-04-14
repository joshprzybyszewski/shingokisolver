package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

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
	ge getEdger,
	coord model.NodeCoord,
) int8 {
	var total int8

	for _, dir := range model.AllCardinals {
		ep := NewEdgePair(coord, dir)
		if ge.GetEdge(ep) == model.EdgeExists {
			total++
		}
	}

	return total
}

func (p *Puzzle) GetSumOutgoingStraightLines(
	coord model.NodeCoord,
) int8 {
	return getSumOutgoingStraightLines(p.edges, coord)
}

func getSumOutgoingStraightLines(
	ge getEdger,
	coord model.NodeCoord,
) int8 {
	var total int8

	for _, dir := range model.AllCardinals {
		c := coord
		ep := NewEdgePair(c, dir)
		for ge.GetEdge(ep) == model.EdgeExists {
			total++
			c = c.Translate(dir)
			ep = NewEdgePair(c, dir)
		}
	}

	return total
}
