package model

// MAX_EDGES = 2 ^ 32 is out constraint because of the uint8
// that we use in a number of placs
const MAX_EDGES = 32

type GridNodeBoundser interface {
	IsInBounds(NodeCoord) bool
}

type GridNodeGetter interface {
	Get(NodeCoord) OutgoingEdges
}

type GridNodeSetter interface {
	Set(NodeCoord, OutgoingEdges)
}

type GridSetterAndGetter interface {
	GridNodeBoundser
	GridNodeGetter
	GridNodeSetter
}

type GridCopyer interface {
	Copy() Grid
}

type Grid interface {
	GridNodeBoundser
	GridNodeGetter
	GridCopyer

	applyUpdates([]gridUpdate)
}

func NewGrid(numEdges int) Grid {
	if numEdges <= 0 || numEdges > MAX_EDGES {
		return nil
	}

	// I'm seeing array backed grids performing better
	// for copys and lookups than arrays of arrays locally.
	// However, once we get to 10x10 puzzles, I need a three
	// dimensional grid structure that will keep pointers around
	// so I don't have to copy _everything_ on copy
	// return newQuadTree(numEdges)

	switch numEdges {
	case 1, 2:
		return &grid3x3{}
	case 3, 4, 5:
		return &grid6x6{}
	case 6:
		return &grid8x8{}
	case 7:
		return &arrayBackedGrid8{}
	case 8, 9:
		return &grid11x11{}
	case 10:
		return &arrayBackedGrid11{}
	case 11, 12, 13, 14, 15:
		return &grid16x16{}
	case 16, 17, 18, 19, 20:
		return &grid21x21{}
	case 21, 22, 23, 24, 25:
		return &grid26x26{}
	default:
		return &maxGrid{}
	}
}
