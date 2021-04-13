package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type getEdger interface {
	GetEdge(
		ep edgePair,
	) model.EdgeState
}

type ruleSet struct {
	rulesByEdges map[edgePair]*rules
}

func newRuleSet(
	numEdges int,
	nodes map[model.NodeCoord]model.Node,
) *ruleSet {
	rs := ruleSet{
		rulesByEdges: make(map[edgePair]*rules, (numEdges-1)*numEdges*2),
	}

	for r := 0; r <= numEdges; r++ {
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)
			if c < numEdges {
				ep := edgePair{
					coord: nc,
					dir:   model.HeadRight,
				}

				rs.rulesByEdges[ep] = newRules(
					ep,
					numEdges,
				)
			}

			if r < numEdges {
				ep := edgePair{
					coord: nc,
					dir:   model.HeadDown,
				}

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
				ep, err := standardizeInput(nc, dir)
				if err != nil {
					continue
				}

				rs.rulesByEdges[ep].addRulesForNode(n, dir)
			}
		}
	}

	return &rs
}

func (rs *ruleSet) getRules(
	ep edgePair,
) *rules {
	return rs.rulesByEdges[ep]
}
