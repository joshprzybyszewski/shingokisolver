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
	allEPs := make([]model.EdgePair, 0, 16)
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
		for _, ep2 := range allEPs {
			perps := ep2.Perpendiculars()
			ep2End := ep2.NodeCoord.Translate(ep2.Cardinal)

			rs.Get(model.NewEdgePair(ep2.NodeCoord, perps[0])).addAffected(ep)
			rs.Get(model.NewEdgePair(ep2.NodeCoord, perps[1])).addAffected(ep)
			rs.Get(model.NewEdgePair(ep2End, perps[0])).addAffected(ep)
			rs.Get(model.NewEdgePair(ep2End, perps[1])).addAffected(ep)
		}
	}
}
