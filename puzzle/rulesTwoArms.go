package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

// TODO clean up how we add all of the two arm rules.
func (rs *ruleSet) addAllTwoArmRules(
	node model.Node,
	numEdges int,
) {

	// TODO consider ways that we can make this set of rules
	// faster/smarter/better/stronger
	// I think we can by dynamically assessin the the "current option" o
	// is still viable or not. If not, then skip the
	// corresponding rule for the node.
	options := model.BuildTwoArmOptions(node, numEdges)

	for _, o := range options {
		afterArm1 := model.NewEdgePair(
			node.Coord().TranslateN(o.One.Heading, int(o.One.Len)),
			o.One.Heading,
		)

		afterArm2 := model.NewEdgePair(
			node.Coord().TranslateN(o.Two.Heading, int(o.Two.Len)),
			o.Two.Heading,
		)

		rs.addExtendedRulesForAvoidedArm(
			node.Coord(),
			o,
			afterArm1,
		)
		rs.addExtendedRulesForAvoidedArm(
			node.Coord(),
			o,
			afterArm2,
		)

		rs.addExtendedRulesForExistingArm(
			node,
			o.One,
			afterArm1, afterArm2,
			o.Two,
		)
		rs.addExtendedRulesForExistingArm(
			node,
			o.Two,
			afterArm2, afterArm1,
			o.One,
		)

		// now _ensure_ that every arm and the off-arm are affecting each other!
		arm1Edges := getArmEdgesAndEnd(node.Coord(), o.One)
		arm2Edges := getArmEdgesAndEnd(node.Coord(), o.Two)
		allEdges := make([]model.EdgePair, 0, len(arm1Edges)+len(arm2Edges)+2)
		allEdges = append(allEdges, arm1Edges...)
		allEdges = append(allEdges, arm2Edges...)
		allEdges = append(allEdges, afterArm1, afterArm2)
		for _, e := range allEdges {
			rs.getRules(e).addAffected(allEdges...)
		}
	}
}

func getArmEdgesAndEnd(
	start model.NodeCoord,
	arm model.Arm,
) []model.EdgePair {
	armEdges := make([]model.EdgePair, 0, arm.Len)
	arm1End := start
	for i := int8(0); i < arm.Len; i++ {
		armEdges = append(armEdges, model.NewEdgePair(arm1End, arm.Heading))
		arm1End = arm1End.Translate(arm.Heading)
	}
	return armEdges
}

func (rs *ruleSet) addExtendedRulesForAvoidedArm(
	nc model.NodeCoord,
	ta model.TwoArms,
	thenAvoid model.EdgePair,
) {
	rs.getRules(thenAvoid).addEvaluations(func(ge model.GetEdger) model.EdgeState {
		if ge.AllExist(nc, ta.One) && ge.AllExist(nc, ta.Two) {
			return model.EdgeAvoided
		}

		return model.EdgeUnknown
	})

}

func (rs *ruleSet) addExtendedRulesForExistingArm(
	node model.Node,
	otherArm model.Arm,
	needToAvoid, needsToNotExist model.EdgePair,
	myArm model.Arm,
) {
	for i, edge := range getArmEdgesAndEnd(node.Coord(), myArm) {
		rs.getRules(edge).addEvaluations(
			getRuleForOppositeArm(
				node,
				otherArm,
				needToAvoid,
				needsToNotExist,
				myArm,
				i == 0,
			),
		)
	}
}

func getRuleForOppositeArm(
	node model.Node,
	otherArm model.Arm,
	needToAvoid model.EdgePair,
	needsToNotExist model.EdgePair,
	myArm model.Arm,
	firstNode bool,
) func(ge model.GetEdger) model.EdgeState {
	return func(ge model.GetEdger) model.EdgeState {
		if !ge.AllExist(node.Coord(), otherArm) {
			return model.EdgeUnknown
		}

		if !ge.IsAvoided(needToAvoid) {
			return model.EdgeUnknown
		}

		// at this point, all of the opposite arm exists, with the edge at the
		// end being avoided.

		if ge.IsEdge(needsToNotExist) {
			if firstNode {
				// this means that the first edge (closest) to the defined
				// node should be avoided because the whole arm cannot be
				// completed as desired.
				return model.EdgeAvoided
			}

			return model.EdgeUnknown
		}

		anyExist, anyAvoided := ge.Any(node.Coord(), myArm)
		if anyAvoided {
			if firstNode {
				return model.EdgeAvoided
			}
			return model.EdgeUnknown
		}

		if !anyExist {
			return model.EdgeUnknown
		}

		if node.Type() == model.BlackNode {
			// for black nodes, we need to check the arm that could be on the opposite side
			oppArm := myArm
			oppArm.Heading = myArm.Heading.Opposite()
			_, anyAvoided = ge.Any(node.Coord(), oppArm)
			if !anyAvoided {
				// we need to know that at _least_ one of the opposite arms is avoided
				// otherwise, we can't claim to know that this one works
				return model.EdgeUnknown
			}
		}

		return model.EdgeExists
	}
}
