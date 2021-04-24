package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

var _ evaluator = advancedNode{}

func newAdvancedNodeEvaluator(
	node model.Node,
	myDir model.Cardinal,
	myIndex int8,
	options []model.TwoArms,
	nearbyNodes map[model.Cardinal][]*model.Node,
) evaluator {
	return advancedNode{
		node:        node,
		dir:         myDir,
		index:       myIndex,
		options:     options,
		nearbyNodes: nearbyNodes,
	}
}

type advancedNode struct {
	nearbyNodes map[model.Cardinal][]*model.Node
	options     []model.TwoArms

	node  model.Node
	dir   model.Cardinal
	index int8
}

func (an advancedNode) evaluate(ge model.GetEdger) model.EdgeState {
	if an.index > 0 && ge.AnyAvoided(an.node.Coord(), model.Arm{
		Heading: an.dir,
		Len:     an.index,
	}) {
		return model.EdgeUnknown
	}

	// TODO come up with a way to cache off this filter...
	minArmByDir, isOnly := model.GetMinArmsByDir(
		an.node.GetFilteredOptions(an.options, ge, an.nearbyNodes),
	)
	if an.index < minArmByDir[an.dir] {
		return model.EdgeExists
	}

	if !isOnly {
		return model.EdgeUnknown
	}

	if myMin, ok := minArmByDir[an.dir]; !ok {
		return model.EdgeUnknown
	} else if an.index == myMin {
		return model.EdgeAvoided
	}

	return model.EdgeUnknown
}
