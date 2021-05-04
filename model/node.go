package model

import (
	"fmt"
)

var (
	InvalidNode = Node{
		coord: InvalidNodeCoord,
	}
)

type Node struct {
	coord NodeCoord
	nType NodeType
	val   int8
}

func NewNode(coord NodeCoord, isWhite bool, v int8) Node {
	nt := BlackNode
	if isWhite {
		nt = WhiteNode
	}
	return Node{
		coord: coord,
		nType: nt,
		val:   v,
	}
}

func (n Node) String() string {
	return fmt.Sprintf("Node{%s%2d @ %+v}", n.nType, n.val, n.coord)
}

func (n Node) PrettyString() string {
	if n.Type() == WhiteNode {
		return fmt.Sprintf("(w%2d)", n.val)
	}
	return fmt.Sprintf("(b%2d)", n.val)
}

func (n Node) Coord() NodeCoord {
	return n.coord
}

func (n Node) Type() NodeType {
	return n.nType
}

func (n Node) Value() int8 {
	return n.val
}

func (n Node) IsInvalidMotions(c1, c2 Cardinal) bool {
	return n.Type().isInvalidMotions(c1, c2)
}

func (n Node) GetFilteredOptions(
	input []TwoArms,
	ge GetEdger,
	nearby NearbyNodes,
) []TwoArms {
	filteredOptions := make([]TwoArms, 0, len(input))

	hInfos := buildHeadingInfos(n, ge)

	for _, o := range input {
		if !n.isTwoArmsPossible(o, hInfos, ge) {
			continue
		}
		if n.isInTheWayOfOtherNodes(o, nearby) {
			continue
		}
		filteredOptions = append(filteredOptions, o)
	}

	return filteredOptions
}

func (n Node) isTwoArmsPossible(
	ta TwoArms,
	his []headingInfo,
	ge GetEdger,
) bool {

	for _, arm := range []Arm{ta.One, ta.Two} {
		if !his[arm.Heading].isValidArm(arm) {
			return false
		}
	}

	nc := n.Coord()
	if ge.AnyAvoided(nc, ta.One) ||
		ge.AnyAvoided(nc, ta.Two) ||
		ge.IsEdge(ta.AfterOne(nc)) ||
		ge.IsEdge(ta.AfterTwo(nc)) {
		return false
	}

	return true
}

func (n Node) isInTheWayOfOtherNodes(
	twoArms TwoArms,
	nearby NearbyNodes,
) bool {

	sumLens := twoArms.One.Len + twoArms.Two.Len

	for _, arm := range []Arm{twoArms.One, twoArms.Two} {
		slVal := arm.Len
		if n.Type() == WhiteNode {
			slVal = sumLens
		}

		if isInTheWay(nearby.Get(arm.Heading), arm.Len, slVal) {
			return true
		}
	}

	return false
}
