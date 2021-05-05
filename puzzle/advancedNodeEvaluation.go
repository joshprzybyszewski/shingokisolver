package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

func checkAdvancedRules(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
	nm *model.NodeMeta,
) model.State {
	if nm == nil {
		panic(`dev error!`)
		return model.Unexpected
	}

	ms := nm.Filter(edges)
	if ms != model.Incomplete {
		if ms == model.Complete {
			return model.Incomplete
		}
		return ms
	}

	for _, dir := range model.AllCardinals {
		ms = checkAdvancedInDirection(edges, rq, rules, nm, dir)
		if ms != model.Incomplete && ms != model.Duplicate {
			return ms
		}
	}

	ms = nm.CheckComplete(edges)
	if ms == model.Complete {
		return model.Incomplete
	}
	return ms
}

func checkAdvancedInDirection(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
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
		if i > 0 && edges.AnyAvoided(nm.Coord(), arm) {
			// we should not have found an avoided edge before the "min arm len"
			panic(`dev error`)
			return model.Unexpected
		}

		if i < minArm {
			ms = addEdge(
				edges,
				ep,
				rq,
				rules,
			)
			if ms != model.Incomplete && ms != model.Duplicate {
				return ms
			}
		}

		if isOnlyOneLength && minArm == i {
			ms = avoidEdge(
				edges,
				ep,
				rq,
				rules,
			)
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
	if len(ta) == 0 {
		// this is unexpected!
		panic(`dev error`)
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
