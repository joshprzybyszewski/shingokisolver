package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ evaluator = twoArmFulfiller{}

type twoArmFulfiller struct {
	allTAs []model.TwoArms
	node   model.Node
}

func newTwoArmFulfiller(
	node model.Node,
	allTAs []model.TwoArms,
) evaluator {

	return twoArmFulfiller{
		allTAs: allTAs,
		node:   node,
	}
}

func (e twoArmFulfiller) evaluate(ge model.GetEdger) model.EdgeState {
	for _, ta := range e.allTAs {
		if !(ge.AnyAvoided(e.node.Coord(), ta.One) ||
			ge.AnyAvoided(e.node.Coord(), ta.Two)) {
			return model.EdgeUnknown
		}
	}

	return model.EdgeAvoided
}
