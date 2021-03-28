package puzzle

type gridNodeGetter interface {
	get(nodeCoord) edgesFromNode
}

type gridNodeSetter interface {
	set(nodeCoord, edgesFromNode)
}

type gridCopyer interface {
	copy() gridNoder
}

type gridNoder interface {
	gridNodeGetter
	gridNodeSetter
	gridCopyer
}

func newGridNoder(numEdges int) gridNoder {
	switch numEdges {
	case 1, 2:
		return &grid3x3{}
	case 3, 4, 5:
		return &grid6x6{}
	case 6, 7:
		return &grid8x8{}
	case 8, 9, 10:
		return &grid11x11{}
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

const MAX_EDGES = 32

type maxGrid [MAX_EDGES][MAX_EDGES]edgesFromNode

var _ gridNoder = (*maxGrid)(nil)

func (mg *maxGrid) get(nc nodeCoord) edgesFromNode {
	return mg[nc.row][nc.col]
}

func (mg *maxGrid) set(nc nodeCoord, efn edgesFromNode) {
	mg[nc.row][nc.col] = efn
}

func (mg *maxGrid) copy() gridNoder {
	cpy := maxGrid{}
	for r := range mg {
		for c := range mg[r] {
			cpy[r][c] = mg[r][c]
		}
	}
	return &cpy
}

type grid3x3 [3][3]edgesFromNode

var _ gridNoder = (*grid3x3)(nil)

func (g *grid3x3) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid3x3) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid3x3) copy() gridNoder {
	cpy := grid3x3{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid6x6 [6][6]edgesFromNode

var _ gridNoder = (*grid6x6)(nil)

func (g *grid6x6) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid6x6) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid6x6) copy() gridNoder {
	cpy := grid6x6{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid8x8 [8][8]edgesFromNode

var _ gridNoder = (*grid8x8)(nil)

func (g *grid8x8) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid8x8) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid8x8) copy() gridNoder {
	cpy := grid8x8{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid11x11 [11][11]edgesFromNode

var _ gridNoder = (*grid11x11)(nil)

func (g *grid11x11) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid11x11) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid11x11) copy() gridNoder {
	cpy := grid11x11{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid16x16 [16][16]edgesFromNode

var _ gridNoder = (*grid16x16)(nil)

func (g *grid16x16) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid16x16) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid16x16) copy() gridNoder {
	cpy := grid16x16{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid21x21 [21][21]edgesFromNode

var _ gridNoder = (*grid21x21)(nil)

func (g *grid21x21) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid21x21) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid21x21) copy() gridNoder {
	cpy := grid21x21{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}

type grid26x26 [26][26]edgesFromNode

var _ gridNoder = (*grid26x26)(nil)

func (g *grid26x26) get(nc nodeCoord) edgesFromNode {
	return g[nc.row][nc.col]
}

func (g *grid26x26) set(nc nodeCoord, efn edgesFromNode) {
	g[nc.row][nc.col] = efn
}

func (g *grid26x26) copy() gridNoder {
	cpy := grid26x26{}
	for r := range g {
		for c := range g[r] {
			cpy[r][c] = g[r][c]
		}
	}
	return &cpy
}
