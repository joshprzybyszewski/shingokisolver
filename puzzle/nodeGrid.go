package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ model.GetNoder = (nodeGrid)(nil)

type nodeGrid [][]*model.Node

func newNodeGrid(nl []model.Node) nodeGrid {

	maxRowIndex := int8(-1)
	for _, n := range nl {
		if int8(n.Coord().Row) > maxRowIndex {
			maxRowIndex = int8(n.Coord().Row)
		}
	}

	ng := make(nodeGrid, maxRowIndex+1)
	for r := int8(0); r < maxRowIndex+1; r++ {
		maxColIndex := int8(-1)
		for _, n := range nl {
			if n.Coord().Row != model.RowIndex(r) {
				continue
			}
			if int8(n.Coord().Col) > maxColIndex {
				maxColIndex = int8(n.Coord().Col)
			}
		}

		row := make([]*model.Node, maxColIndex+1)

		for _, n := range nl {
			if n.Coord().Row != model.RowIndex(r) {
				continue
			}
			n := n
			row[n.Coord().Col] = &n
		}

		ng[r] = row
	}

	return ng
}

func (ng nodeGrid) GetNode(nc model.NodeCoord) (model.Node, bool) {
	if ng != nil && nc.Row >= 0 && nc.Row < model.RowIndex(len(ng)) {
		nRow := ng[nc.Row]
		if nc.Col >= 0 && nc.Col < model.ColIndex(len(nRow)) {
			if n := nRow[nc.Col]; n != nil {
				return *n, true
			}
		}
	}

	return model.Node{}, false
}
