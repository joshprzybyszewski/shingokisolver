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
	nearbyNodes model.NearbyNodes,
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
	nearbyNodes model.NearbyNodes
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

	// TODO come up with a way to cache off this filter so that rules on
	// multiple edges can use the same results. But this may not work well
	// at all in my system...
	minArmByDir, isOnly := model.GetMinArmsByDir(
		an.node.GetFilteredOptions(an.options, ge, an.nearbyNodes),
	)
	if an.index < minArmByDir[an.dir] {
		return model.EdgeExists
	}

	if !isOnly {
		return model.EdgeUnknown
	}

	if minArmByDir[an.dir] == an.index {
		return model.EdgeAvoided
	}

	return model.EdgeUnknown
}
