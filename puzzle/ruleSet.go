package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type ruleSet struct {
	rows [][]*rules
	cols [][]*rules
}

func newRuleSet(
	numEdges int,
	nodes map[model.NodeCoord]model.Node,
) *ruleSet {
	rs := &ruleSet{
		rows: make([][]*rules, numEdges+1),
		cols: make([][]*rules, numEdges),
	}

	for r := 0; r <= numEdges; r++ {
		rs.rows[r] = make([]*rules, numEdges)
		if r < numEdges {
			rs.cols[r] = make([]*rules, numEdges+1)
		}
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)
			if c < numEdges {
				ep := model.NewEdgePair(nc, model.HeadRight)
				rs.rows[r][c] = newRules(
					ep,
					numEdges,
				)
			}

			if r < numEdges {
				ep := model.NewEdgePair(nc, model.HeadDown)
				rs.cols[r][c] = newRules(
					ep,
					numEdges,
				)
			}
		}
	}

	for nc, n := range nodes {
		for _, dir := range model.AllCardinals {
			rs.getRules(model.NewEdgePair(nc, dir)).addRulesForNode(n, dir)
		}

		rs.addAllTwoArmRules(n, numEdges)
	}

	return rs
}

func (rs *ruleSet) getRules(
	ep model.EdgePair,
) *rules {
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
