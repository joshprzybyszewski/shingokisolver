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
	maxLensByDir := GetMaxArmsByDir(tas)
	otherNodes := make(map[Cardinal][]*Node, len(maxLensByDir))

	for dir, maxLen := range maxLensByDir {
		slice := make([]*Node, maxLen)
		foundAny := false
		nc := myNode.Coord()
		for i := range otherNodes[dir] {
			n, ok := gn.GetNode(nc)
			if ok {
				foundAny = true
				otherNodes[dir][i] = &n
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
	return !ge.AnyAvoided(nc, ta.One) &&
		!ge.AnyAvoided(nc, ta.Two) &&
		!ge.IsEdge(ta.AfterOne(nc)) &&
		!ge.IsEdge(ta.AfterTwo(nc))
}

func isInTheWayOfOtherNodes(
	myNode Node,
	ta TwoArms,
	otherNodes map[Cardinal][]*Node,
) bool {

	isInTheWay := func(otherNodes []*Node, maxLen int, myStraightLineVal int8) bool {
		for i, otherNode := range otherNodes {
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

	a1StraightLineVal := ta.One.Len
	a2StraightLineVal := ta.Two.Len
	if myNode.Type() == WhiteNode {
		a1StraightLineVal = ta.One.Len + ta.Two.Len
		a2StraightLineVal = ta.One.Len + ta.Two.Len
	}

	if isInTheWay(otherNodes[ta.One.Heading], int(ta.One.Len), a1StraightLineVal) {
		return true
	}
	if isInTheWay(otherNodes[ta.Two.Heading], int(ta.Two.Len), a2StraightLineVal) {
		return true
	}

	return false
}
