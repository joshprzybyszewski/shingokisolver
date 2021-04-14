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

			extEvals := getExtendedEvalsForNode(n, numEdges)
			for ep, evals := range extEvals {
				rs.rulesByEdges[ep].addExtendedEval(n, evals)
			}
		}
	}

	return &rs
}

func getExtendedEvalsForNode(
	node model.Node,
	numEdges int,
) map[EdgePair]extendedRules {

	options := model.BuildTwoArmOptions(node, numEdges)

	res := make(map[EdgePair]extendedRules, len(options))

	for _, o := range options {
		arm1Edges, offArm1 := getArmEdgesAndEnd(node.Coord(), o.One)
		arm2Edges, offArm2 := getArmEdgesAndEnd(node.Coord(), o.Two)

		allArmEdges := make([]EdgePair, 0, len(arm1Edges)+len(arm2Edges))
		allArmEdges = append(allArmEdges, arm1Edges...)
		allArmEdges = append(allArmEdges, arm2Edges...)

		eval := func(ge getEdger) model.EdgeState {
			for _, ep := range allArmEdges {
				if ge.GetEdge(ep) != model.EdgeExists {
					return model.EdgeUnknown
				}
			}
			return model.EdgeAvoided
		}

		for _, armEP := range []EdgePair{offArm1, offArm2} {
			ext, ok := res[armEP]
			if !ok {
				res[armEP] = ext
			}
			ext.evals = append(ext.evals, eval)
		}

		for _, ep := range allArmEdges {
			ext, ok := res[ep]
			if !ok {
				res[ep] = ext
			}
			ext.couldAffect = append(ext.couldAffect, offArm1, offArm2)
		}
	}

	return res
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

func (rs *ruleSet) getRules(
	ep EdgePair,
) *rules {
	return rs.rulesByEdges[ep]
}
