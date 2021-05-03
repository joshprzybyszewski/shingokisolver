package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (rs *RuleSet) AddAllTwoArmRules(
	node model.Node,
	gn model.GetNoder,
	options []model.TwoArms,
) {

	maxArms := model.GetMaxArmsByDir(options)
	nearbyNodes := model.BuildNearbyNodes(node, gn, maxArms)
	allEPs := make([]model.EdgePair, 0, 8)
	for dir, maxLen := range maxArms {
		start := node.Coord()
		ep := model.NewEdgePair(start, dir)
		for i := int8(0); i <= maxLen; i++ {
			allEPs = append(allEPs, ep)
			rs.Get(ep).addEvaluation(
				newAdvancedNodeEvaluator(
					node,
					dir,
					i,
					options,
					nearbyNodes,
				),
			)
			ep = ep.Next(dir)
		}
	}

	for _, ep := range allEPs {
		rs.Get(ep).addAffected(allEPs...)
	}
}
