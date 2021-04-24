package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type RuleSet struct {
	rows [][]*Rules
	cols [][]*Rules
}

func New(
	ge model.GetEdger,
	numEdges int,
	nodes []model.Node,
) *RuleSet {
	rs := &RuleSet{
		rows: make([][]*Rules, numEdges+1),
		cols: make([][]*Rules, numEdges),
	}

	for r := 0; r <= numEdges; r++ {
		rs.rows[r] = make([]*Rules, numEdges)
		if r < numEdges {
			rs.cols[r] = make([]*Rules, numEdges+1)
		}
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)
			if c < numEdges {
				ep := model.NewEdgePair(nc, model.HeadRight)
				rs.rows[r][c] = newRules(
					ge,
					ep,
				)
			}

			if r < numEdges {
				ep := model.NewEdgePair(nc, model.HeadDown)
				rs.cols[r][c] = newRules(
					ge,
					ep,
				)
			}
		}
	}

	for _, n := range nodes {
		nc := n.Coord()
		for _, dir := range model.AllCardinals {
			rs.Get(model.NewEdgePair(nc, dir)).addSimpleNodeRules(n, dir)
		}
	}

	return rs
}

func (rs *RuleSet) Get(
	ep model.EdgePair,
) *Rules {
	// copied from isInBounds:#
	if ep.Row < 0 || ep.Col < 0 {
		// negative coords are bad
		return nil
	}

	switch ep.Cardinal {
	case model.HeadRight:
		if int(ep.Row) < len(rs.rows) && int(ep.Col) < len(rs.rows[0]) {
			return rs.rows[ep.Row][ep.Col]
		}
	case model.HeadDown:
		if int(ep.Row) < len(rs.cols) && int(ep.Col) < len(rs.cols[0]) {
			return rs.cols[ep.Row][ep.Col]
		}
	default:
		// unexpected input
	}
	return nil
}
