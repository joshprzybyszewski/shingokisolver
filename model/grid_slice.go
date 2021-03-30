package model

type gridSlice8 [64]int32

var _ Grid = (*gridSlice8)(nil)

func (g *gridSlice8) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int(nc.Row) < 8 && int(nc.Col) < 8
}

func (g *gridSlice8) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*8 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g[index])
}

func (g *gridSlice8) Set(nc NodeCoord, oe OutgoingEdges) {
	index := int(nc.Row)*8 + int(nc.Col)

	g[index] = outgoingEdgesToInt32(oe)
}

func (g *gridSlice8) Copy() Grid {
	cpy := gridSlice8{}
	for i, v := range g {
		cpy[i] = v
	}
	return &cpy
}

type gridSlice11 [121]int32

var _ Grid = (*gridSlice11)(nil)

func (g *gridSlice11) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int(nc.Row) < 11 && int(nc.Col) < 11
}

func (g *gridSlice11) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*11 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g[index])
}

func (g *gridSlice11) Set(nc NodeCoord, oe OutgoingEdges) {
	index := int(nc.Row)*11 + int(nc.Col)

	g[index] = outgoingEdgesToInt32(oe)
}

func (g *gridSlice11) Copy() Grid {
	cpy := gridSlice11{}
	for i, v := range g {
		cpy[i] = v
	}
	return &cpy
}
