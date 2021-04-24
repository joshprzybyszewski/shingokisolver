package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type nodeList []model.Node

func (nl nodeList) GetNode(nc model.NodeCoord) (model.Node, bool) {
	// TODO replace nodeList with something better?
	for _, n := range nl {
		if n.Coord() == nc {
			return n, true
		}
	}
	return model.Node{}, false
}

func (nl nodeList) toNodeGrid() nodeGrid {
	maxRowIndex := -1
	maxColIndex := -1
	for _, n := range nl {
		if int(n.Coord().Row) > maxRowIndex {
			maxRowIndex = int(n.Coord().Row)
		}
		if int(n.Coord().Col) > maxColIndex {
			maxColIndex = int(n.Coord().Col)
		}
	}

	ng := make(nodeGrid, maxRowIndex+1)
	for r := 0; r < len(ng); r++ {
		row := make([]*model.Node, maxColIndex+1)
		ng[r] = row
	}

	for _, n := range nl {
		n := n
		ng[n.Coord().Row][n.Coord().Col] = &n
	}

	return ng
}

type nodeGrid [][]*model.Node

func (ng nodeGrid) GetNode(nc model.NodeCoord) (model.Node, bool) {
	if ng != nil && int(nc.Row) < len(ng) {
		if nRow := ng[nc.Row]; nRow != nil && int(nc.Col) < len(nRow) {
			if n := nRow[nc.Col]; n != nil {
				return *n, true
			}
		}
	}

	return model.Node{}, false
}
