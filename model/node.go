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
		foundAny := false
		nc := myNode.Coord()
		for i := range slice {
			n, ok := gn.GetNode(nc)
			if ok {
				foundAny = true
				slice[i] = &n
			}
			nc = nc.Translate(dir)
		}
		if foundAny {
			otherNodes[dir] = slice
		}
	}

	return otherNodes
}

func (n Node) GetFilteredOptions(
	input []TwoArms,
	ge GetEdger,
	otherNodes map[Cardinal][]*Node,
) []TwoArms {
	filteredOptions := make([]TwoArms, 0, len(input))

	// TODO I have a lot of loops that are doing duplicate checks (isTwoArmsPossible)
	for _, o := range input {
		if !isTwoArmsPossible(n, o, ge) {
			continue
		}
		if isInTheWayOfOtherNodes(n, o, otherNodes) {
			continue
		}
		filteredOptions = append(filteredOptions, o)
	}

	return filteredOptions
}

func isTwoArmsPossible(
	node Node,
	ta TwoArms,
	ge GetEdger,
) bool {

	nc := node.Coord()
	if ge.AnyAvoided(nc, ta.One) ||
		ge.AnyAvoided(nc, ta.Two) ||
		ge.IsEdge(ta.AfterOne(nc)) ||
		ge.IsEdge(ta.AfterTwo(nc)) {
		return false
	}

	for _, arm := range []Arm{ta.One, ta.Two} {
		cur := nc.Translate(arm.Heading)
		for i := int8(1); i < arm.Len; i++ {
			for _, dir := range arm.Heading.Perpendiculars() {
				if ge.IsEdge(NewEdgePair(cur, dir)) {
					return false
				}
			}
			cur = cur.Translate(arm.Heading)
		}
		bothAvoided := true
		for _, dir := range arm.Heading.Perpendiculars() {
			if !ge.IsAvoided(NewEdgePair(cur, dir)) {
				bothAvoided = false
				break
			}
		}
		if bothAvoided {
			return false
		}
	}

	return true
}

func isInTheWayOfOtherNodes(
	myNode Node,
	twoArms TwoArms,
	otherNodes map[Cardinal][]*Node,
) bool {

	isInTheWay := func(otherNodes []*Node, maxLen int, myStraightLineVal int8) bool {
		// TODO this function has a lot of opportunity for improvement
		for i, otherNode := range otherNodes {
			if i == 0 {
				// we can skip over "this" node
				continue
			}
			if i > maxLen {
				return false
			} else if i == maxLen {
				if otherNode == nil {
					return false
				}
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
				return false
			}

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

	a1StraightLineVal := twoArms.One.Len
	a2StraightLineVal := twoArms.Two.Len
	if myNode.Type() == WhiteNode {
		a1StraightLineVal = twoArms.One.Len + twoArms.Two.Len
		a2StraightLineVal = twoArms.One.Len + twoArms.Two.Len
	}

	if isInTheWay(otherNodes[twoArms.One.Heading], int(twoArms.One.Len), a1StraightLineVal) {
		return true
	}
	if isInTheWay(otherNodes[twoArms.Two.Heading], int(twoArms.Two.Len), a2StraightLineVal) {
		return true
	}

	return false
}
