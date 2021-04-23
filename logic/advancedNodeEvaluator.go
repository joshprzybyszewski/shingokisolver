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
	node  model.Node
	dir   model.Cardinal
	index int8

	options     []model.TwoArms
	nearbyNodes map[model.Cardinal][]*model.Node
}

func (an advancedNode) evaluate(ge model.GetEdger) model.EdgeState {
	if ge.AnyAvoided(an.node.Coord(), model.Arm{
		Heading: an.dir,
		Len:     an.index,
	}) {
		return model.EdgeUnknown
	}

	filteredTAs := an.node.GetFilteredOptions(an.options, ge, an.nearbyNodes)
	minArmByDir, isOnly := model.GetMinArmsByDir(filteredTAs)
	if an.index < minArmByDir[an.dir] {
		return model.EdgeExists
	}

	if isOnly && an.index == minArmByDir[an.dir] {
		return model.EdgeAvoided
	}

	return model.EdgeUnknown
}
