package model

type arrayBackedGrid8 struct {
	n    int8
	grid [64]int32
}

var _ Grid = (*arrayBackedGrid8)(nil)

func (g *arrayBackedGrid8) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int8(nc.Row) < g.n && int8(nc.Col) < g.n
}

func (g *arrayBackedGrid8) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*int(g.n) + int(nc.Col)
	return newOutgoingEdgesFromInt32(g.grid[index])
}

func (g *arrayBackedGrid8) applyUpdates(updates []gridUpdate) {
	for _, u := range updates {
		nc := u.coord
		index := int(nc.Row)*int(g.n) + int(nc.Col)

		g.grid[index] = outgoingEdgesToInt32(u.newVal)
	}
}

func (g *arrayBackedGrid8) Copy() Grid {
	cpy := arrayBackedGrid8{
		n: g.n,
	}
	for i, v := range g.grid {
		cpy.grid[i] = v
	}
	return &cpy
}

type arrayBackedGrid11 struct {
	n    int8
	grid [121]int32
}

var _ Grid = (*arrayBackedGrid11)(nil)

func (g *arrayBackedGrid11) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int8(nc.Row) < g.n && int8(nc.Col) < g.n
}

func (g *arrayBackedGrid11) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)*int(g.n) + int(nc.Col)
	return newOutgoingEdgesFromInt32(g.grid[index])
}

func (g *arrayBackedGrid11) applyUpdates(updates []gridUpdate) {
	for _, u := range updates {
		nc := u.coord
		index := int(nc.Row)*int(g.n) + int(nc.Col)

		g.grid[index] = outgoingEdgesToInt32(u.newVal)
	}
}

func (g *arrayBackedGrid11) Copy() Grid {
	cpy := arrayBackedGrid11{
		n: g.n,
	}
	for i, v := range g.grid {
		cpy.grid[i] = v
	}
	return &cpy
}

type maxSliceBackedArray struct {
	n    int8
	grid []int32
}

var _ Grid = (*maxSliceBackedArray)(nil)

func (g maxSliceBackedArray) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return int8(nc.Row) < g.n && int8(nc.Col) < g.n
}

func (g maxSliceBackedArray) Get(nc NodeCoord) OutgoingEdges {
	index := int(nc.Row)<<5 + int(nc.Col)
	return newOutgoingEdgesFromInt32(g.grid[index])
}

func (g maxSliceBackedArray) applyUpdates(updates []gridUpdate) {
	for _, u := range updates {
		nc := u.coord
		index := int(nc.Row)*int(g.n) + int(nc.Col)

		g.grid[index] = outgoingEdgesToInt32(u.newVal)
	}
}

func (g maxSliceBackedArray) Copy() Grid {
	cpy := maxSliceBackedArray{
		n: g.n,
	}
	cpy.grid = make([]int32, len(g.grid))
	copy(cpy.grid, g.grid)
	return cpy
}
