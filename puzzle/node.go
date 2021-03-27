package puzzle

type edgesFromNode struct {
	above int8
	below int8
	left  int8
	right int8
}

func (efn edgesFromNode) totalEdges() int8 {
	return efn.above + efn.below + efn.left + efn.right
}

func (efn edgesFromNode) isabove() bool {
	return efn.above != 0
}

func (efn edgesFromNode) isbelow() bool {
	return efn.below != 0
}

func (efn edgesFromNode) isleft() bool {
	return efn.left != 0
}

func (efn edgesFromNode) isright() bool {
	return efn.right != 0
}

func (efn edgesFromNode) getNumOutgoingDirections() int8 {
	var numBranches int8

	if efn.isabove() {
		numBranches++
	}
	if efn.isbelow() {
		numBranches++
	}
	if efn.isleft() {
		numBranches++
	}
	if efn.isright() {
		numBranches++
	}

	return numBranches
}

type nodeType uint8

const (
	// no constraints
	noNode nodeType = 0
	// must be passed through in a straight line
	whiteNode nodeType = 1
	// must be turned upon
	blackNode nodeType = 2
)

func (nt nodeType) isInvalidEdges(efn edgesFromNode) bool {
	switch nt {
	case whiteNode:
		// white nodes need to be straight. therefore, they're
		// invalid if they have opposing directions set
		return (efn.isabove() || efn.isbelow()) && (efn.isleft() || efn.isright())
	case blackNode:
		// black nodes need to be bent. therefore, they're
		// invalid if they have a straight line in them
		return (efn.isabove() && efn.isbelow()) || (efn.isleft() && efn.isright())
	default:
		return false
	}
}

type node struct {
	nType nodeType
	val   int8
}

func (n node) copy() node {
	return n
}
