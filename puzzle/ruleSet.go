package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type ruleSet struct {
	rulesByEdges map[model.EdgePair]*rules
}

func newRuleSet(
	numEdges int,
	nodes map[model.NodeCoord]model.Node,
) *ruleSet {
	rs := &ruleSet{
		rulesByEdges: make(map[model.EdgePair]*rules, (numEdges-1)*numEdges*2),
	}

	for r := 0; r <= numEdges; r++ {
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)
			if c < numEdges {
				ep := model.NewEdgePair(nc, model.HeadRight)
				rs.rulesByEdges[ep] = newRules(
					ep,
					numEdges,
				)
			}

			if r < numEdges {
				ep := model.NewEdgePair(nc, model.HeadDown)
				rs.rulesByEdges[ep] = newRules(
					ep,
					numEdges,
				)
			}
		}
	}

	for nc, n := range nodes {
		for _, dir := range model.AllCardinals {
			rs.rulesByEdges[model.NewEdgePair(nc, dir)].addRulesForNode(n, dir)
		}

		rs.addAllTwoArmRules(n, numEdges)
	}

	return rs
}

func (rs *ruleSet) getRules(
	ep model.EdgePair,
) *rules {
	return rs.rulesByEdges[ep]
}
