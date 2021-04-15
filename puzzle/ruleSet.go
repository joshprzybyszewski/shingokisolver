package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type getEdger interface {
	GetEdge(
		ep EdgePair,
	) model.EdgeState
}

type ruleSet struct {
	rulesByEdges map[EdgePair]*rules
}

func newRuleSet(
	numEdges int,
	nodes map[model.NodeCoord]model.Node,
) *ruleSet {
	rs := ruleSet{
		rulesByEdges: make(map[EdgePair]*rules, (numEdges-1)*numEdges*2),
	}

	for r := 0; r <= numEdges; r++ {
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)
			if c < numEdges {
				ep := NewEdgePair(nc, model.HeadRight)
				rs.rulesByEdges[ep] = newRules(
					ep,
					numEdges,
				)
			}

			if r < numEdges {
				ep := NewEdgePair(nc, model.HeadDown)
				rs.rulesByEdges[ep] = newRules(
					ep,
					numEdges,
				)
			}

			n, ok := nodes[nc]
			if !ok {
				continue
			}

			for _, dir := range model.AllCardinals {
				rs.rulesByEdges[NewEdgePair(nc, dir)].addRulesForNode(n, dir)
			}

			rs.addAllTwoArmRules(n, numEdges)
		}
	}

	return &rs
}

func (rs *ruleSet) getRules(
	ep EdgePair,
) *rules {
	return rs.rulesByEdges[ep]
}
