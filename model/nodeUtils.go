package model

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

func isInTheWay(otherNodes []*Node, maxLen int8, myStraightLineVal int8) bool {
	// if we end on a node, then we need to check it's value against the straight line logic
	if maxLen < int8(len(otherNodes)) {
		otherNode := otherNodes[maxLen]
		if otherNode != nil {
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
