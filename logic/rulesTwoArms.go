package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (rs *RuleSet) AddAllTwoArmRules(
	node model.Node,
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

		allArmEdges := getAllEdges(node.Coord(), ta)

		allArmEdgesWithAfters := make([]model.EdgePair, len(allArmEdges), len(allArmEdges)+2)
		copy(allArmEdgesWithAfters, allArmEdges)
		allArmEdgesWithAfters = append(allArmEdgesWithAfters, afterArm1, afterArm2)

		// ensure that the after-arm1 will be avoided when appropriate
		rs.Get(afterArm1).addEvaluations(
			newAfterArmEvaluator(node.Coord(), ta),
		)
		rs.Get(afterArm1).addAffected(allArmEdgesWithAfters...)

		// ensure that the after-arm2 will be avoided when appropriate
		rs.Get(afterArm2).addEvaluations(
			newAfterArmEvaluator(node.Coord(), ta),
		)
		rs.Get(afterArm2).addAffected(allArmEdgesWithAfters...)

		for _, edge := range allArmEdges {
			rs.Get(edge).addEvaluations(
				newWithinArmEvaluator(
					node,
					ta,
					afterArm1, afterArm2,
					edge,
				),
			)

			// now _ensure_ that every arm and the off-arm are affecting each other!
			rs.Get(edge).addAffected(allArmEdgesWithAfters...)
		}
	}
}

func getAllEdges(
	start model.NodeCoord,
	ta model.TwoArms,
) []model.EdgePair {
	armEdges := make([]model.EdgePair, 0, ta.One.Len+ta.Two.Len)

	curEnd := start
	for i := int8(0); i < ta.One.Len; i++ {
		armEdges = append(armEdges, model.NewEdgePair(curEnd, ta.One.Heading))
		curEnd = curEnd.Translate(ta.One.Heading)
	}

	curEnd = start
	for i := int8(0); i < ta.Two.Len; i++ {
		armEdges = append(armEdges, model.NewEdgePair(curEnd, ta.Two.Heading))
		curEnd = curEnd.Translate(ta.Two.Heading)
	}

	return armEdges
}
