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
	minArm, isOnlyOneLength := getMinArmInDir(
		an.dir,
		an.node.GetFilteredOptions(an.options, ge, an.nearbyNodes),
	)
	if an.index < minArm {
		return model.EdgeExists
	}

	if isOnlyOneLength && minArm == an.index {
		return model.EdgeAvoided
	}

	return model.EdgeUnknown
}

func getMinArmInDir(
	dir model.Cardinal,
	ta []model.TwoArms,
) (int8, bool) {
	if len(ta) == 0 {
		// this is unexpected!
		// TODO panic
		return -1, false
	}

	var min int8
	if ta[0].One.Heading == dir {
		min = ta[0].One.Len
	} else if ta[0].Two.Heading == dir {
		min = ta[0].Two.Len
	} else {
		// the first TwoArms option doesn't go in the right direction.
		return -1, false
	}

	isOnlyOneLength := true

	for i := 1; i < len(ta); i++ {
		switch dir {
		case ta[i].One.Heading:
			if min != ta[i].One.Len {
				isOnlyOneLength = false
			}
			if min > ta[i].One.Len {
				min = ta[i].One.Len
			}
		case ta[i].Two.Heading:
			if min != ta[i].Two.Len {
				isOnlyOneLength = false
			}
			if min > ta[i].Two.Len {
				min = ta[i].Two.Len
			}
		default:
			// this TwoArms option doesn't go in the right direction.
			return -1, false
		}
	}

	return min, isOnlyOneLength
}
