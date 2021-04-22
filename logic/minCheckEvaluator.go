package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

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
	// return model.EdgeUnknown
	return mac.evaluateBlackNode(ge)
}

func (mac minArmCheck) evaluateWhiteNode(ge model.GetEdger) model.EdgeState {
	if mac.isAvoidedEdgeBeforeOrAtMe(ge) {
		return model.EdgeUnknown
	}

	if !ge.IsEdge(model.NewEdgePair(mac.node.Coord(), mac.myDir)) {
		// I can't know if I'm supposed to exist or not if my starting
		// edge doesn't exist.
		return model.EdgeUnknown
	}

	return mac.evalTheUnseen(ge, mac.myDir.Opposite())

	oppPossible, oppLastChainStartIndex, oppLastChainLen := getNumInDirection(
		ge,
		int(mac.node.Value()),
		mac.node.Coord(),
		mac.myDir.Opposite(),
	)

	if oppLastChainStartIndex == 0 {
		// TODO figure out what this means?
	} else if oppLastChainStartIndex > 0 {
		myLen := getNumExistingInDirection(
			ge,
			mac.node.Coord(),
			mac.myDir,
		)

		if myLen+oppLastChainLen+oppLastChainStartIndex > int(mac.node.Value()) {
			oppPossible = oppLastChainStartIndex - 1
		}
	}

	// TODO clean up this logic?
	if oppPossible < int(mac.node.Value()) {
		if mac.myIndex < int(mac.node.Value())-oppPossible {
			return model.EdgeExists
		}
		return model.EdgeUnknown
	}

	return model.EdgeUnknown
}

func (mac minArmCheck) isAvoidedEdgeBeforeOrAtMe(ge model.GetEdger) bool {
	return ge.AnyAvoided(mac.node.Coord(), model.Arm{
		Heading: mac.myDir,
		Len:     int8(mac.myIndex) + 1,
	})
}

func (mac minArmCheck) evaluateBlackNode(ge model.GetEdger) model.EdgeState {
	if mac.isAvoidedEdgeBeforeOrAtMe(ge) {
		return model.EdgeUnknown
	}

	if ge.IsEdge(model.NewEdgePair(mac.node.Coord(), mac.myDir)) {
		otherSeen := false
		for _, dir := range mac.myDir.Perpendiculars() {
			if ge.IsAvoided(model.NewEdgePair(mac.node.Coord(), dir)) ||
				mac.evalTheUnseen(ge, dir) == model.EdgeExists {
				if otherSeen {
					return model.EdgeExists
				}
				otherSeen = true
			}
		}
		return model.EdgeUnknown
	}
	if mac.myIndex > 0 {
		return model.EdgeUnknown
	}

	numUnavoided := getNumUnavoided(
		ge,
		mac.node.Coord(),
		mac.myDir,
	)
	perp0NumUnavoided := getNumUnavoided(
		ge,
		mac.node.Coord(),
		mac.myDir.Perpendiculars()[0],
	)
	perp1NumUnavoided := getNumUnavoided(
		ge,
		mac.node.Coord(),
		mac.myDir.Perpendiculars()[1],
	)
	maxOtherArm := perp0NumUnavoided
	if perp1NumUnavoided > maxOtherArm {
		maxOtherArm = perp1NumUnavoided
	}
	if maxOtherArm+numUnavoided < int(mac.node.Value()) {
		return model.EdgeAvoided
	}

	oppNumUnavoided := getNumUnavoided(
		ge,
		mac.node.Coord(),
		mac.myDir.Opposite(),
	)
	if maxOtherArm+oppNumUnavoided < int(mac.node.Value()) {
		return model.EdgeExists
	}

	return model.EdgeUnknown
}

func (mac minArmCheck) evalTheUnseen(
	ge model.GetEdger,
	dir model.Cardinal,
) model.EdgeState {
	dirPossible, dirLastChainStartIndex, dirLastChainLen := getNumInDirection(
		ge,
		int(mac.node.Value()),
		mac.node.Coord(),
		dir,
	)

	if dirLastChainStartIndex == 0 {
		// TODO figure out what this means?
		return model.EdgeUnknown
	} else if dirLastChainStartIndex > 0 {
		myLen := getNumExistingInDirection(
			ge,
			mac.node.Coord(),
			mac.myDir,
		)

		if myLen+dirLastChainLen+dirLastChainStartIndex > int(mac.node.Value()) {
			dirPossible = dirLastChainStartIndex - 1
		}
	}

	if dirPossible > 0 && mac.myIndex < int(mac.node.Value())-dirPossible {
		return model.EdgeExists
	}
	return model.EdgeUnknown
}

func getNumExistingInDirection(
	ge model.GetEdger,
	start model.NodeCoord,
	dir model.Cardinal,
) (numExisting int) {

	cur := start

	for ge.IsEdge(model.NewEdgePair(cur, dir)) {
		numExisting++
		cur = cur.Translate(dir)
	}

	return numExisting
}

func getNumInDirection(
	ge model.GetEdger,
	nValue int,
	start model.NodeCoord,
	dir model.Cardinal,
) (numPossibleEdges, existingChainStartIndex, numInLastExistingChain int) {
	existingChainStartIndex = -1

	cur := start

	for !ge.IsAvoided(model.NewEdgePair(cur, dir)) {
		if numPossibleEdges >= nValue && existingChainStartIndex == -1 {
			break
		}
		if ge.IsEdge(model.NewEdgePair(cur, dir)) {
			if existingChainStartIndex == -1 {
				existingChainStartIndex = numPossibleEdges
			}
			numInLastExistingChain++
		} else {
			if numPossibleEdges >= nValue {
				break
			}
			existingChainStartIndex = -1
			numInLastExistingChain = 0
		}

		numPossibleEdges++
		cur = cur.Translate(dir)
	}

	return numPossibleEdges, existingChainStartIndex, numInLastExistingChain
}

func getNumUnavoided(
	ge model.GetEdger,
	start model.NodeCoord,
	dir model.Cardinal,
) (numUnavoided int) {

	cur := start

	for !ge.IsAvoided(model.NewEdgePair(cur, dir)) {
		numUnavoided++
		cur = cur.Translate(dir)
	}

	return numUnavoided
}
