package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ evaluator = withinArm{}

type withinArm struct {
	node model.Node

	oppArm   model.Arm
	afterOpp model.EdgePair

	myArm   model.Arm
	afterMe model.EdgePair

	isFirstEdge bool
}

func newWithinArmEvaluator(
	node model.Node,
	ta model.TwoArms,
	after1, after2 model.EdgePair,
	me model.EdgePair,
) evaluator {
	oppArm := ta.Two
	afterOpp := after2
	myArm := ta.One
	afterMe := after1

	switch oppArm.Heading {
	case me.Cardinal, me.Cardinal.Opposite():
		oppArm = ta.One
		afterOpp = after1
		myArm = ta.Two
		afterMe = after2
	}

	return withinArm{
		node:        node,
		oppArm:      oppArm,
		afterOpp:    afterOpp,
		myArm:       myArm,
		afterMe:     afterMe,
		isFirstEdge: model.NewEdgePair(node.Coord(), myArm.Heading) == me,
	}
}

func (wa withinArm) evaluate(ge model.GetEdger) model.EdgeState {
	if !ge.IsAvoided(wa.afterOpp) {
		return model.EdgeUnknown
	}
	nc := wa.node.Coord()

	if !ge.AllExist(nc, wa.oppArm) {
		return model.EdgeUnknown
	}

	// at this point, all of the opposite arm exists, with the edge at the
	// end being avoided.

	// TODO I think this evaluator could be improved if it kept track of
	// what index it is in myArm

	if wa.node.Type() == model.BlackNode {
		if ge.IsEdge(wa.afterMe) {
			if wa.isFirstEdge {
				// this means that the first edge (closest) to the defined
				// node should be avoided because the whole arm cannot be
				// completed as desired.
				return model.EdgeAvoided
			}

			return model.EdgeUnknown
		}

		anyAvoided := ge.AnyAvoided(nc, wa.myArm)
		if anyAvoided {
			if wa.isFirstEdge {
				return model.EdgeAvoided
			}
			return model.EdgeUnknown
		}

		// for black nodes, we need to check the arm that could be on the opposite side
		oppArm := wa.myArm
		oppArm.Heading = wa.myArm.Heading.Opposite()
		anyAvoided = ge.AnyAvoided(nc, oppArm)
		if !anyAvoided {
			// we need to know that at _least_ one of the opposite arms is avoided
			// otherwise, we can't claim to know that this one works
			return model.EdgeUnknown
		}
	}

	return model.EdgeExists
}
