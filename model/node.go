package model

import (
	"fmt"
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

func (n Node) Copy() Node {
	return n
}

func (n Node) IsInvalidMotions(c1, c2 Cardinal) bool {
	return n.Type().isInvalidMotions(c1, c2)
}

func BuildNearbyNodes(
	myNode Node,
	tas []TwoArms,
	gn GetNoder,
) map[Cardinal][]*Node {
	// TODO see if I can't speed this up
	maxLensByDir := GetMaxArmsByDir(tas)
	otherNodes := make(map[Cardinal][]*Node, len(maxLensByDir))

	for dir, maxLen := range maxLensByDir {
		slice := make([]*Node, maxLen)
		nc := myNode.Coord().Translate(dir)
		lastNodeIndex := 0
		for i := 1; i < len(slice); i++ {
			n, ok := gn.GetNode(nc)
			if ok {
				lastNodeIndex = i
				slice[i] = &n
			}
			nc = nc.Translate(dir)
		}
		otherNodes[dir] = slice[:lastNodeIndex+1]
	}

	return otherNodes
}

func (n Node) GetFilteredOptions(
	input []TwoArms,
	ge GetEdger,
	otherNodes map[Cardinal][]*Node,
) []TwoArms {
	filteredOptions := make([]TwoArms, 0, len(input))

	hInfos := buildHeadingInfos(n, ge)

	for _, o := range input {
		if !isTwoArmsPossible(n, o, hInfos, ge) {
			continue
		}
		if isInTheWayOfOtherNodes(n, o, otherNodes) {
			continue
		}
		filteredOptions = append(filteredOptions, o)
	}

	return filteredOptions
}

func buildHeadingInfos(
	node Node,
	ge GetEdger,
) map[Cardinal]headingInfo {
	res := make(map[Cardinal]headingInfo, 4)

	nc := node.Coord()
	maxLen := node.Value() - 1

	for _, heading := range AllCardinals {
		hInfo := headingInfo{
			armLensDoubleAvoids: make(map[int8]struct{}, 2),
		}
		cur := nc.Translate(heading)
		for i := int8(1); i <= maxLen; i++ {
			bothAvoided := true
			for _, dir := range heading.Perpendiculars() {
				if ge.IsEdge(NewEdgePair(cur, dir)) {
					hInfo.maxArmLenUntilIncomingEdge = i
					break
				}
				if !ge.IsAvoided(NewEdgePair(cur, dir)) {
					bothAvoided = false
					break
				}
			}

			if hInfo.maxArmLenUntilIncomingEdge > 0 {
				break
			}
			if bothAvoided {
				hInfo.armLensDoubleAvoids[i] = struct{}{}
			}

			cur = cur.Translate(heading)
		}

		res[heading] = hInfo
	}

	return res
}

type headingInfo struct {
	maxArmLenUntilIncomingEdge int8
	armLensDoubleAvoids        map[int8]struct{}
}

func (hi headingInfo) isValidArm(arm Arm) bool {
	if hi.maxArmLenUntilIncomingEdge > 0 &&
		hi.maxArmLenUntilIncomingEdge < arm.Len {
		return false
	}

	if _, ok := hi.armLensDoubleAvoids[arm.Len]; ok {
		return false
	}

	return true
}

func isTwoArmsPossible(
	node Node,
	ta TwoArms,
	his map[Cardinal]headingInfo,
	ge GetEdger,
) bool {

	for _, arm := range []Arm{ta.One, ta.Two} {
		if !his[arm.Heading].isValidArm(arm) {
			return false
		}
	}

	nc := node.Coord()
	if ge.AnyAvoided(nc, ta.One) ||
		ge.AnyAvoided(nc, ta.Two) ||
		ge.IsEdge(ta.AfterOne(nc)) ||
		ge.IsEdge(ta.AfterTwo(nc)) {
		return false
	}

	return true
}

func isInTheWayOfOtherNodes(
	myNode Node,
	twoArms TwoArms,
	otherNodes map[Cardinal][]*Node,
) bool {

	a1StraightLineVal := twoArms.One.Len
	a2StraightLineVal := twoArms.Two.Len
	if myNode.Type() == WhiteNode {
		a1StraightLineVal = twoArms.One.Len + twoArms.Two.Len
		a2StraightLineVal = twoArms.One.Len + twoArms.Two.Len
	}

	if isInTheWay(otherNodes[twoArms.One.Heading], twoArms.One.Len, a1StraightLineVal) {
		return true
	}
	if isInTheWay(otherNodes[twoArms.Two.Heading], twoArms.Two.Len, a2StraightLineVal) {
		return true
	}

	return false
}

func isInTheWay(otherNodes []*Node, maxLen int8, myStraightLineVal int8) bool {

	// if we end on a node, then we have different logic
	if maxLen < int8(len(otherNodes)) {
		if otherNode := otherNodes[maxLen]; otherNode != nil {
			if otherNode.Type() == WhiteNode {
				// this arm would end in a white node. That's not ok because
				// we would need to continue through it
				return true
			}
			if otherNode.Value()-myStraightLineVal < 1 {
				// this arm meets the other node, and would require going
				// next in a perpendicular path. Since this arm would
				// contribute too much to its value, we can filter it ou.
				return true
			}
		}
	}

	// if we pass through a node, then we need to validate that's ok
	for i := int8(1); i < maxLen && i < int8(len(otherNodes)); i++ {
		otherNode := otherNodes[i]

		if otherNode == nil {
			continue
		}

		if otherNode.Type() == BlackNode {
			// this arm would pass through this node in a straight line
			// that makes this arm impossible.
			return true
		}
		if otherNode.Value() != myStraightLineVal {
			// this arm would pass through the other node
			// in a straight line, and the value would not be tenable
			return true
		}
	}

	return false
}
