package model

type arrayBackedGrid8 [64]int32

var _ Grid = (*arrayBackedGrid8)(nil)

func (g *arrayBackedGrid8) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int(nc.Row) < 8 && int(nc.Col) < 8
}

func (g *arrayBackedGrid8) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*8 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g[index])
}

func (g *arrayBackedGrid8) Set(nc NodeCoord, oe OutgoingEdges) {
	index := int(nc.Row)*8 + int(nc.Col)

	g[index] = outgoingEdgesToInt32(oe)
}

func (g *arrayBackedGrid8) Copy() Grid {
	cpy := arrayBackedGrid8{}
	for i, v := range g {
		cpy[i] = v
	}
	return &cpy
}

type arrayBackedGrid11 [121]int32

var _ Grid = (*arrayBackedGrid11)(nil)

func (g *arrayBackedGrid11) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int(nc.Row) < 11 && int(nc.Col) < 11
}

func (g *arrayBackedGrid11) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*11 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g[index])
}

func (g *arrayBackedGrid11) Set(nc NodeCoord, oe OutgoingEdges) {
	index := int(nc.Row)*11 + int(nc.Col)

	g[index] = outgoingEdgesToInt32(oe)
}

func (g *arrayBackedGrid11) Copy() Grid {
	cpy := arrayBackedGrid11{}
	for i, v := range g {
		cpy[i] = v
	}
	return &cpy
}

type maxSliceBackedArray []int32

var _ Grid = (maxSliceBackedArray)(nil)

func newMaxSliceBackedArray() maxSliceBackedArray {
	return make(maxSliceBackedArray, int(MAX_EDGES)*int(MAX_EDGES))
}

func (g maxSliceBackedArray) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int(nc.Row) < MAX_EDGES && int(nc.Col) < MAX_EDGES
}

func (g maxSliceBackedArray) Get(nc NodeCoord) OutgoingEdges {
	// int(nc.Row)*MAX_EDGES = int(nc.Row) << 5
	index := int(nc.Row)<<5 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g[index])
}

func (g maxSliceBackedArray) Set(nc NodeCoord, oe OutgoingEdges) {
	index := int(nc.Row)<<5 + int(nc.Col)

	g[index] = outgoingEdgesToInt32(oe)
}

func (g maxSliceBackedArray) Copy() Grid {
	cpy := make(maxSliceBackedArray, len(g))
	copy(cpy, g)
	return cpy
}
