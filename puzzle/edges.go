package puzzle

// should match rowIndex and colIndex
type edgeIndex int8

// This means that the max puzzle size we can support is 32 edges
type edges int32

const MAX_EDGES = 32

func newEdges() edges {
	return 0
}

func (e edges) isEdge(start edgeIndex) bool {
	if start < 0 || start >= MAX_EDGES {
		// sanity check. we could probably remove for speed
		return false
	}
	return e&(1<<start) != 0
}

func (e edges) addEdge(start edgeIndex) edges {
	if start < 0 || start >= MAX_EDGES {
		return e
	}
	return e | (1 << start)
}
