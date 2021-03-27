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
