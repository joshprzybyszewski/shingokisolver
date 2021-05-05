package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (su stateUpdater) checkAdvancedRules(
	n model.Node,
) model.State {
	var nm *model.NodeMeta

	for _, m := range su.nms {
		// TODO if I can get this faster than iterating through all of nms, that'd be great...
		if m.Coord() == n.Coord() {
			nm = m
			break
		}
	}

	// if nm == nil {
	// 	panic(fmt.Errorf("dev error: did not find NodeMeta for coord: %s", n))
	// 	return model.Unexpected
	// }

	if nm.IsComplete {
		return model.Complete
	}

	ms := nm.Filter(su.edges)
	if ms != model.Incomplete {
		if ms == model.Complete {
			return model.Incomplete
		}
		return ms
	}

	for _, dir := range model.AllCardinals {
		ms = su.checkAdvancedInDirection(nm, dir)
		if ms != model.Incomplete && ms != model.Duplicate {
			return ms
		}
	}

	ms = nm.CheckComplete(su.edges)
	if ms == model.Complete {
		return model.Incomplete
	}
	return ms
}

func (su stateUpdater) checkAdvancedInDirection(
	nm *model.NodeMeta,
	dir model.Cardinal,
) model.State {
	var ms model.State

	minArm, isOnlyOneLength := getMinArmInDir(
		dir,
		nm.TwoArmOptions,
	)

	arm := model.Arm{
		Heading: dir,
	}
	ep := model.NewEdgePair(nm.Coord(), dir)
	for i := int8(0); i <= minArm; i++ {
		arm.Len = i
		// if i > 0 && su.edges.AnyAvoided(nm.Coord(), arm) {
		// 	// we should not have found an avoided edge before the "min arm len"
		// 	panic(`dev error`)
		// 	return model.Unexpected
		// }

		if i < minArm {
			ms = su.addEdge(ep)
			if ms != model.Incomplete && ms != model.Duplicate {
				return ms
			}
		}

		if isOnlyOneLength && minArm == i {
			ms = su.avoidEdge(ep)
			if ms != model.Incomplete && ms != model.Duplicate {
				return ms
			}
		}

		ep = ep.Next(dir)
	}

	return model.Incomplete
}

func getMinArmInDir(
	dir model.Cardinal,
	ta []model.TwoArms,
) (int8, bool) {
	// if len(ta) == 0 {
	// 	// this is unexpected!
	// 	panic(`dev error`)
	// 	return -1, false
	// }

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
