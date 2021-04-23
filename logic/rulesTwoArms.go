package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (rs *RuleSet) AddAllTwoArmRules(
	node model.Node,
	gn model.GetNoder,
	options []model.TwoArms,
) {

	for _, ta := range options {
		afterArm1 := model.NewEdgePair(
			node.Coord().TranslateAlongArm(ta.One),
			ta.One.Heading,
		)

		afterArm2 := model.NewEdgePair(
			node.Coord().TranslateAlongArm(ta.Two),
			ta.Two.Heading,
		)

		// ensure that the after-arm1 will be avoided when appropriate
		rs.Get(afterArm1).addEvaluation(
			newAfterArmEvaluator(node.Coord(), ta),
		)
		// ensure that the after-arm2 will be avoided when appropriate
		rs.Get(afterArm2).addEvaluation(
			newAfterArmEvaluator(node.Coord(), ta),
		)
	}

	maxArms := model.GetMaxArmsByDir(options)
	nearbyNodes := model.BuildNearbyNodes(node, options, gn)
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

func getTwoArmsWithDirection(
	allOptions []model.TwoArms,
	dir model.Cardinal,
) []model.TwoArms {
	filtered := make([]model.TwoArms, 0, len(allOptions))

	for _, ta := range allOptions {
		if ta.One.Heading == dir || ta.Two.Heading == dir {
			filtered = append(filtered, ta)
		}
	}

	return filtered
}
