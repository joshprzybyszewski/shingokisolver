package model

type maxGrid [MAX_EDGES][MAX_EDGES]OutgoingEdges

var _ Grid = (*maxGrid)(nil)

func (mg *maxGrid) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(mg)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (mg *maxGrid) Get(nc NodeCoord) OutgoingEdges {
	return mg[nc.Row][nc.Col]
}

func (mg *maxGrid) Set(nc NodeCoord, efn OutgoingEdges) {
	mg[nc.Row][nc.Col] = efn
}

func (mg *maxGrid) Copy() Grid {
	cpy := maxGrid{}
	for r := range mg {
		for c := range mg[r] {
			cpy[r][c] = mg[r][c]
		}
	}
	return &cpy
}

type grid3x3 [3][3]OutgoingEdges

var _ Grid = (*grid3x3)(nil)

func (g *grid3x3) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid3x3) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid3x3) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid3x3) Copy() Grid {
	cpy := grid3x3{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid6x6 [6][6]OutgoingEdges

var _ Grid = (*grid6x6)(nil)

func (g *grid6x6) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid6x6) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid6x6) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid6x6) Copy() Grid {
	cpy := grid6x6{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid8x8 [8][8]OutgoingEdges

var _ Grid = (*grid8x8)(nil)

func (g *grid8x8) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid8x8) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid8x8) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid8x8) Copy() Grid {
	cpy := grid8x8{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid11x11 [11][11]OutgoingEdges

var _ Grid = (*grid11x11)(nil)

func (g *grid11x11) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid11x11) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid11x11) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid11x11) Copy() Grid {
	cpy := grid11x11{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid16x16 [16][16]OutgoingEdges

var _ Grid = (*grid16x16)(nil)

func (g *grid16x16) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid16x16) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid16x16) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid16x16) Copy() Grid {
	cpy := grid16x16{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid21x21 [21][21]OutgoingEdges

var _ Grid = (*grid21x21)(nil)

func (g *grid21x21) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid21x21) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid21x21) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid21x21) Copy() Grid {
	cpy := grid21x21{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid26x26 [26][26]OutgoingEdges

var _ Grid = (*grid26x26)(nil)

func (g *grid26x26) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	maxIndex := len(g)
	return int(nc.Row) < maxIndex && int(nc.Col) < maxIndex
}

func (g *grid26x26) Get(nc NodeCoord) OutgoingEdges {
	return g[nc.Row][nc.Col]
}

func (g *grid26x26) Set(nc NodeCoord, efn OutgoingEdges) {
	g[nc.Row][nc.Col] = efn
}

func (g *grid26x26) Copy() Grid {
	cpy := grid26x26{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}
