package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ evaluator = minArmCheck{}

type minArmCheck struct {
	node model.Node

	myDir   model.Cardinal
	myIndex int
}

func newMinArmCheckEvaluator(
	node model.Node,
	myDir model.Cardinal,
	myIndex int,
) evaluator {

	return minArmCheck{
		node:    node,
		myDir:   myDir,
		myIndex: myIndex,
	}
}

func (mac minArmCheck) evaluate(ge model.GetEdger) model.EdgeState {
	if mac.node.Type() == model.WhiteNode {
		return mac.evaluateWhiteNode(ge)
	}

	// TODO I'm still having trouble with evaluateBlackNode
	// I can skip it for now and everything is fine.
	return model.EdgeUnknown
	return mac.evaluateBlackNode(ge)
}

func (mac minArmCheck) evaluateWhiteNode(ge model.GetEdger) model.EdgeState {
	if !mac.isActiveWhiteNodeEvaluator(ge) {
		// I can't know if I'm supposed to exist or not if my starting
		// edge doesn't exist.
		return model.EdgeUnknown
	}

	maxPossibleExisting := mac.getMaxPossibleInDir(ge, mac.myDir.Opposite())

	if maxPossibleExisting == 0 {
		// This is an error state because of the check at the beginning
		// to see if my starter was started.
		return model.EdgeErrored
	}

	if (mac.myIndex+1)+maxPossibleExisting < int(mac.node.Value()) {
		return model.EdgeExists
	}

	return model.EdgeUnknown
}

func (mac minArmCheck) isAvoidedEdgeBeforeMe(ge model.GetEdger) bool {
	return mac.getMaxPossibleInDir(ge, mac.myDir) < mac.myIndex
}

func (mac minArmCheck) isActiveWhiteNodeEvaluator(ge model.GetEdger) bool {
	if mac.isAvoidedEdgeBeforeMe(ge) {
		return false
	}

	if !ge.IsEdge(model.NewEdgePair(mac.node.Coord(), mac.myDir)) {
		// I can't know if I'm supposed to exist or not if my starting
		// edge doesn't exist.
		return false
	}

	return true
}

func (mac minArmCheck) evaluateBlackNode(ge model.GetEdger) model.EdgeState {
	isActive, maxPerps := mac.isActiveBlackNodeEvaluator(ge)
	if !isActive {
		// I can't know if I'm supposed to exist or not if my starting
		// edge doesn't exist.
		return model.EdgeUnknown
	}

	if maxPerps <= 0 {
		// this node should have at least one arm perpendicular to me!
		return model.EdgeErrored
	}

	if (mac.myIndex+1)+maxPerps < int(mac.node.Value()) {
		return model.EdgeExists
	}

	return model.EdgeUnknown
}

func (mac minArmCheck) isActiveBlackNodeEvaluator(
	ge model.GetEdger,
) (bool, int) {
	if mac.isAvoidedEdgeBeforeMe(ge) {
		return false, -1
	}

	maxPerps := mac.getMaxPossiblePerpendiculars(ge)
	if ge.IsEdge(model.NewEdgePair(mac.node.Coord(), mac.myDir)) {
		return true, maxPerps
	}

	// At this point, I know that my starting edge doesn't exist. That's ok, I may
	// still be able to determine that my opposite side can't
	// fulfill the requirements of this black node.

	// Consider this:
	//                           X
	//                         (   )
	//                           ?
	//                         (   )
	//                           ?
	// (   ) X (   ) ? (   ) ? (b 5) ? (   ) ? (   ) X (   )
	//                           ?
	//                         (   )
	//                           ?
	//                         (   )
	//                           ?
	//                         (   )
	// I need to return false for the HeadUp dir, and true for the HeadDown dir

	// Also consider the following:
	// (   )
	//   ?
	// (   )
	//   ?
	// (   )
	//   ?
	// (b 3)---(   ) X (   )
	//   ?
	// (   )
	//   |
	// (   )
	//   |
	// (   )
	// It's clear to you and me that we must HeadUp in this case, and
	// cannot HeadDown. We need to ensure that the "opposite max" can
	// reflect this.

	myMax := mac.getMaxPossibleInDir(ge, mac.myDir)
	if maxPerps+myMax < int(mac.node.Value()) {
		return false, -1
	}
	maxOpposite := mac.getMaxPossibleInDir(ge, mac.myDir.Opposite())
	if maxPerps+maxOpposite < int(mac.node.Value()) {
		return true, maxPerps
	}

	// Consider this:
	// (b 6)---(   )---(   )---(   ) ? (   )
	//   |
	// (   )
	//   ?
	// (   )
	//   ?
	// (   )
	//   ?
	// (   )
	//   ?
	// (   )
	// I need to return "false, not active" for these edges that are beyond
	// the reach of this arm now.
	minPerps := mac.getMinPossiblePerpendiculars(ge)
	if mac.myIndex >= int(mac.node.Value())-minPerps {
		return false, -1
	}

	return true, maxPerps
}

func (mac minArmCheck) getMaxPossiblePerpendiculars(
	ge model.GetEdger,
) int {
	max := 0
	for _, perp := range mac.myDir.Perpendiculars() {
		v := mac.getMaxPossibleInDir(ge, perp)
		if v > max {
			max = v
		}
	}
	return max
}

func (mac minArmCheck) getMaxPossibleInDir(
	ge model.GetEdger,
	dir model.Cardinal,
) int {
	maxPossibleExisting := 0
	indexStartOfLastExistingChain := -1

	cur := mac.node.Coord()

	for !ge.IsAvoided(model.NewEdgePair(cur, dir)) {
		if ge.IsEdge(model.NewEdgePair(cur, dir)) {
			if indexStartOfLastExistingChain == -1 {
				indexStartOfLastExistingChain = maxPossibleExisting
			}
		} else {
			indexStartOfLastExistingChain = -1
		}

		maxPossibleExisting++
		cur = cur.Translate(dir)

		if maxPossibleExisting == int(mac.node.Value()) {
			break
		}
	}

	if maxPossibleExisting == int(mac.node.Value()) {
		// this means we got all the way out for the node.
		// we need to check if there's an existing chain of
		// edges before we got to here.
		if indexStartOfLastExistingChain > 0 {
			return indexStartOfLastExistingChain
		}
	}

	return maxPossibleExisting
}

func (mac minArmCheck) getMinPossiblePerpendiculars(
	ge model.GetEdger,
) int {
	min := -1
	for _, perp := range mac.myDir.Perpendiculars() {
		v := mac.getMinPossibleInDir(ge, perp)
		if v < min || min == -1 {
			min = v
		}
	}
	return min
}

func (mac minArmCheck) getMinPossibleInDir(
	ge model.GetEdger,
	dir model.Cardinal,
) int {
	minPossible := 0

	cur := mac.node.Coord()

	for ge.IsEdge(model.NewEdgePair(cur, dir)) {
		minPossible++
		cur = cur.Translate(dir)
	}

	return minPossible
}
