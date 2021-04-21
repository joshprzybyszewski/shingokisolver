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

		allArmEdges := ta.GetAllEdges(node.Coord())

		allArmEdgesWithAfters := make([]model.EdgePair, len(allArmEdges), len(allArmEdges)+2)
		copy(allArmEdgesWithAfters, allArmEdges)
		allArmEdgesWithAfters = append(allArmEdgesWithAfters, afterArm1, afterArm2)

		// ensure that the after-arm1 will be avoided when appropriate
		rs.Get(afterArm1).addEvaluation(
			newAfterArmEvaluator(node.Coord(), ta),
		)
		rs.Get(afterArm1).addAffected(allArmEdgesWithAfters...)

		// ensure that the after-arm2 will be avoided when appropriate
		rs.Get(afterArm2).addEvaluation(
			newAfterArmEvaluator(node.Coord(), ta),
		)
		rs.Get(afterArm2).addAffected(allArmEdgesWithAfters...)

		for _, edge := range allArmEdges {
			rs.Get(edge).addEvaluation(
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

	for _, dir := range model.AllCardinals {
		firstEdge := model.NewEdgePair(
			node.Coord(),
			dir,
		)
		rs.Get(firstEdge).addEvaluation(
			newTwoArmFulfiller(
				node,
				getTwoArmsWithDirection(options, dir),
			),
		)
	}

	if node.Value() > 2 {
		for _, dir := range model.AllCardinals {
			cur := node.Coord()
			for edgeIndex := 0; edgeIndex < int(node.Value()); edgeIndex++ {
				ep := model.NewEdgePair(
					cur,
					dir,
				)
				rs.Get(ep).addEvaluation(
					newMinArmCheckEvaluator(
						node,
						dir,
						edgeIndex,
					),
				)
				cur = cur.Translate(dir)
			}
		}
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
