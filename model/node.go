package model

import "fmt"

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

func (n Node) GetFilteredOptions(
	input []TwoArms,
	ge GetEdger,
	gn GetNoder,
) []TwoArms {
	filteredOptions := make([]TwoArms, 0, len(input))

	for _, o := range input {
		if !isTwoArmsPossible(n, o, ge) {
			continue
		}
		// TODO isInTheWayOfOtherNodes could be cache across input
		// and would be faster
		if isInTheWayOfOtherNodes(n, o, gn) {
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
	return !ge.AnyAvoided(nc, ta.One) &&
		!ge.AnyAvoided(nc, ta.Two) &&
		!ge.IsEdge(ta.AfterOne(nc)) &&
		!ge.IsEdge(ta.AfterTwo(nc))
}

func isInTheWayOfOtherNodes(
	node Node,
	ta TwoArms,
	gn GetNoder,
) bool {

	nc := node.Coord()

	a1StraightLineVal := ta.One.Len
	a2StraightLineVal := ta.Two.Len
	if node.Type() == WhiteNode {
		a1StraightLineVal = ta.One.Len + ta.Two.Len
		a2StraightLineVal = ta.One.Len + ta.Two.Len
	}

	for i, a1 := 1, nc; i < int(ta.One.Len); i++ {
		a1 = a1.Translate(ta.One.Heading)
		otherNode, ok := gn.GetNode(a1)
		if !ok {
			continue
		}
		if otherNode.Type() == BlackNode {
			// this arm would pass through this node in a straight line
			// that makes this arm impossible.
			return true
		}
		if otherNode.Value() != a1StraightLineVal {
			// this arm would pass through the other node
			// in a straight line, and the value would not be tenable
			return true
		}
	}
	if otherNode, ok := gn.GetNode(nc.TranslateAlongArm(ta.One)); ok {
		if otherNode.Type() == WhiteNode {
			// this arm would end in a white node. That's not ok because
			// we would need to continue through it
			return true
		}
		if otherNode.Value()-a1StraightLineVal < 1 {
			// this arm meets the other node, and would require going
			// next in a perpendicular path. Since this arm would
			// contribute too much to its value, we can filter it ou.
			return true
		}
	}

	for i, a2 := 1, nc; i < int(ta.Two.Len); i++ {
		a2 = a2.Translate(ta.Two.Heading)
		otherNode, ok := gn.GetNode(a2)
		if !ok {
			continue
		}
		if otherNode.Type() == BlackNode {
			// this arm would pass through this node in a straight line
			// that makes this arm impossible.
			return true
		}
		if otherNode.Value() != a2StraightLineVal {
			// this arm would pass through the other node
			// in a straight line, and the value would not be tenable
			return true
		}
	}

	if otherNode, ok := gn.GetNode(nc.TranslateAlongArm(ta.Two)); ok {
		if otherNode.Type() == WhiteNode {
			// this arm would end in a white node. That's not ok because
			// we would need to continue through it
			return true
		}
		if otherNode.Value()-a2StraightLineVal < 1 {
			// this arm meets the other node, and would require going
			// next in a perpendicular path. Since this arm would
			// contribute too much to its value, we can filter it ou.
			return true
		}
	}

	return false
}
