package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (rs *ruleSet) addAllTwoArmRules(
	node model.Node,
	numEdges int,
) {

	options := model.BuildTwoArmOptions(node, numEdges)

	for _, o := range options {
		arm1Edges, afterArm1 := getArmEdgesAndEnd(node.Coord(), o.One)
		arm2Edges, afterArm2 := getArmEdgesAndEnd(node.Coord(), o.Two)

		allExistingArms := make([]EdgePair, 0, len(arm1Edges)+len(arm2Edges))
		allExistingArms = append(allExistingArms, arm1Edges...)
		allExistingArms = append(allExistingArms, arm2Edges...)

		rs.addExtendedRulesForAvoidedArm(allExistingArms, afterArm1)
		rs.addExtendedRulesForAvoidedArm(allExistingArms, afterArm2)

		rs.addExtendedRulesForExistingArm(arm1Edges, afterArm1, afterArm2, arm2Edges)
		rs.addExtendedRulesForExistingArm(arm2Edges, afterArm2, afterArm1, arm1Edges)
	}
}

func getArmEdgesAndEnd(
	start model.NodeCoord,
	arm model.Arm,
) ([]EdgePair, EdgePair) {
	armEdges := make([]EdgePair, 0, arm.Len)
	arm1End := start
	for i := int8(0); i < arm.Len; i++ {
		armEdges = append(armEdges, NewEdgePair(arm1End, arm.Heading))
		arm1End = arm1End.Translate(arm.Heading)
	}
	return armEdges, NewEdgePair(arm1End, arm.Heading)
}

func (rs *ruleSet) addExtendedRulesForAvoidedArm(
	needToExist []EdgePair,
	thenAvoid EdgePair,
) {
	avoidExt := rs.rulesByEdges[thenAvoid]
	if avoidExt == nil {
		// nothing to do!
		return
	}

	avoidExt.addEvaluations(func(ge getEdger) model.EdgeState {
		printDebugMsg("running check for avoided end-of-arm")

		for _, ep := range needToExist {
			if ge.GetEdge(ep) != model.EdgeExists {
				// one of the edges that needs to exist doesn't
				// therefore we can't know for sure
				return model.EdgeUnknown
			}
		}
		return model.EdgeAvoided
	})

	for _, ep := range needToExist {
		rs.rulesByEdges[ep].addAffected(thenAvoid)
	}
}

func (rs *ruleSet) addExtendedRulesForExistingArm(
	firstArm []EdgePair,
	endOfFirstArm EdgePair,
	endOfSecondArm EdgePair,
	secondArm []EdgePair,
) {
	rs.rulesByEdges[endOfFirstArm].addAffected(secondArm...)
	for _, ep := range firstArm {
		rs.rulesByEdges[ep].addAffected(secondArm...)
	}

	for _, ep := range secondArm {
		rs.rulesByEdges[ep].addAffected(firstArm...)
		rs.rulesByEdges[ep].addAffected(endOfFirstArm)
	}

	for i, edge := range secondArm {
		rs.rulesByEdges[edge].addEvaluations(getRuleForOppositeArm(
			firstArm,
			endOfFirstArm,
			endOfSecondArm,
			secondArm,
			i == 0,
		))
	}
}

func getRuleForOppositeArm(
	needToExist []EdgePair,
	needToAvoid EdgePair,
	needsToNotExist EdgePair,
	thenExists []EdgePair,
	firstNode bool,
) func(ge getEdger) model.EdgeState {
	return func(ge getEdger) model.EdgeState {
		printDebugMsg(
			"running check for opposite arm\n\toppArm: %+v\n\tendOfArm: %s\n\tmyArm: %+v\n\tisFirstNode: %v\n ",
			needToExist,
			needToAvoid,
			thenExists,
			firstNode,
		)

		switch s := ge.GetEdge(needToAvoid); s {
		case model.EdgeAvoided, model.EdgeOutOfBounds:
			// the one we need to avoid _is_ in fact avoided
		default:
			printDebugMsg(
				"edge needed to avoid (%s) was not avoided: %s",
				needToAvoid,
				s,
			)
			return model.EdgeUnknown
		}

		for _, ep := range needToExist {
			if s := ge.GetEdge(ep); s != model.EdgeExists {
				printDebugMsg(
					"edge needed to exist (%s) was not existing: %s",
					ep,
					s,
				)
				// one of the edges that needs to exist doesn't. say "we don't know"
				return model.EdgeUnknown
			}
		}

		switch s := ge.GetEdge(needsToNotExist); s {
		case model.EdgeExists:
			if firstNode {
				printDebugMsg(
					"FIRST NODE: edge that should have NOT existed (%s) was existing",
					needsToNotExist,
				)
				// this means that the first edge (closest) to the defined
				// node should be avoided because the whole arm cannot be
				// completed as desired.
				return model.EdgeAvoided
			}

			printDebugMsg(
				"edge that should have NOT existed (%s) was existing",
				needsToNotExist,
			)
			return model.EdgeUnknown
		}

		// all of the opposite arm exists, with the edge at the
		// end being avoided.

		shouldExist := false
		for _, thenExist := range thenExists {
			switch s := ge.GetEdge(thenExist); s {
			case model.EdgeExists:
				// and at least one node in this arm
				// exists, so we can say this entire arm exists.
				shouldExist = true
			case model.EdgeAvoided, model.EdgeOutOfBounds:
				// Surprise! one of the edges in this arm is
				// avoided (or, unexpectedly, out of bounds)
				if firstNode {
					printDebugMsg(
						"FIRST NODE: edge that also should have existed (%s) was not existing: %s",
						thenExist,
						s,
					)
					// this means that the first edge (closest) to the defined
					// node should be avoided because the whole arm cannot be
					// completed as desired.
					return model.EdgeAvoided
				}

				printDebugMsg(
					"edge that also should have existed (%s) was not existing: %s",
					thenExist,
					s,
				)
				// since this isn't the first edge, we don't have enough
				// info, so we can't claim knowledge
				return model.EdgeUnknown
			}
		}

		if shouldExist {
			return model.EdgeExists
		}

		printDebugMsg(
			"none of the shouldExists existed: %+v",
			thenExists,
		)

		return model.EdgeUnknown
	}
}
