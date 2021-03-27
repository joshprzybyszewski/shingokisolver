package puzzle

type edgesFromNode struct {
	above      int8
	below      int8
	left       int8
	right      int8
	totalEdges int8

	isabove bool
	isbelow bool
	isleft  bool
	isright bool

	isPopulated bool
}

func newEdgesFromNode(
	above, below, left, right int8,
) edgesFromNode {
	return edgesFromNode{
		totalEdges:  above + below + left + right,
		above:       above,
		below:       below,
		left:        left,
		right:       right,
		isabove:     above != 0,
		isbelow:     below != 0,
		isleft:      left != 0,
		isright:     right != 0,
		isPopulated: true,
	}
}

func (efn edgesFromNode) getNumOutgoingDirections() int8 {
	var numBranches int8

	if efn.isabove {
		numBranches++
	}
	if efn.isbelow {
		numBranches++
	}
	if efn.isleft {
		numBranches++
	}
	if efn.isright {
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
		return (efn.isabove || efn.isbelow) && (efn.isleft || efn.isright)
	case blackNode:
		// black nodes need to be bent. therefore, they're
		// invalid if they have a straight line in them
		return (efn.isabove && efn.isbelow) || (efn.isleft && efn.isright)
	default:
		return false
	}
}

type node struct {
	nType nodeType
	val   int8

	seen bool
}

func (n *node) copy() *node {
	if n == nil {
		return nil
	}
	return &node{
		nType: n.nType,
		val:   n.val,
	}
}
