package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

// TODO clean this up and move it
func (p *Puzzle) GetLooseEnd() (model.NodeCoord, model.State) {
	for r := model.RowIndex(0); r <= model.RowIndex(p.numEdges); r++ {
		for c := model.ColIndex(0); c <= model.ColIndex(p.numEdges); c++ {
			nc := model.NewCoord(r, c)
			oe, ok := p.GetOutgoingEdgesFrom(nc)
			if !ok {
				return model.InvalidNodeCoord, model.Violation
			}

			switch numEdges := oe.GetNumOutgoingDirections(); {
			case numEdges > 2:
				return model.InvalidNodeCoord, model.Violation
			case numEdges == 1:
				return nc, model.Incomplete
			}
		}
	}
	return model.InvalidNodeCoord, model.NodesComplete
}

// TODO clean this up (maybe remove it?) and move it
func (p *Puzzle) GetOutgoingEdgesFrom(
	coord model.NodeCoord,
) (model.OutgoingEdges, bool) {
	if coord.Row < 0 ||
		coord.Col < 0 ||
		uint8(coord.Row) > p.numEdges ||
		uint8(coord.Col) > p.numEdges {
		return model.OutgoingEdges{}, false
	}

	var above, below, left, right int8

	c := coord
	ep, err := standardizeInput(c, model.HeadUp)
	for err == nil && p.edges.GetEdge(ep) == model.EdgeExists {
		above++
		c = c.Translate(model.HeadUp)
		ep, err = standardizeInput(c, model.HeadUp)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadDown)
	for err == nil && p.edges.GetEdge(ep) == model.EdgeExists {
		below++
		c = c.Translate(model.HeadDown)
		ep, err = standardizeInput(c, model.HeadDown)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadLeft)
	for err == nil && p.edges.GetEdge(ep) == model.EdgeExists {
		left++
		c = c.Translate(model.HeadLeft)
		ep, err = standardizeInput(c, model.HeadLeft)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadRight)
	for err == nil && p.edges.GetEdge(ep) == model.EdgeExists {
		right++
		c = c.Translate(model.HeadRight)
		ep, err = standardizeInput(c, model.HeadRight)
	}

	return model.NewOutgoingEdges(above, below, left, right), true
}

func getNumStraightLineOutgoingEdges(
	ge getEdger,
	coord model.NodeCoord,
) int8 {
	var total int8

	c := coord
	ep, err := standardizeInput(c, model.HeadUp)
	for err == nil && ge.GetEdge(ep) == model.EdgeExists {
		total++
		c = c.Translate(model.HeadUp)
		ep, err = standardizeInput(c, model.HeadUp)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadDown)
	for err == nil && ge.GetEdge(ep) == model.EdgeExists {
		total++
		c = c.Translate(model.HeadDown)
		ep, err = standardizeInput(c, model.HeadDown)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadLeft)
	for err == nil && ge.GetEdge(ep) == model.EdgeExists {
		total++
		c = c.Translate(model.HeadLeft)
		ep, err = standardizeInput(c, model.HeadLeft)
	}

	c = coord
	ep, err = standardizeInput(c, model.HeadRight)
	for err == nil && ge.GetEdge(ep) == model.EdgeExists {
		total++
		c = c.Translate(model.HeadRight)
		ep, err = standardizeInput(c, model.HeadRight)
	}

	return total
}
