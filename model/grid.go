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

	switch numEdges {
	case 1, 2, 3, 4, 5, 6, 7:
		return &arrayBackedGrid8{
			n: int8(numEdges) + 1,
		}
	case 8, 9, 10:
		return &arrayBackedGrid11{
			n: int8(numEdges) + 1,
		}
	default:
		return &maxSliceBackedArray{
			n:    int8(numEdges) + 1,
			grid: make([]int32, (numEdges+1)*(numEdges+1)),
		}
	}
}
