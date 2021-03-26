package puzzlegrid

import (
	"errors"
	"fmt"
	"strings"
)

type NodeLocation struct {
	Row     int
	Col     int
	IsWhite bool
	Value   int
}

type cardinal int

const (
	headNowhere cardinal = 0
	headUp      cardinal = 1
	headLeft    cardinal = 2
	headDown    cardinal = 3
	headRight   cardinal = 4
)

type row int
type col int

// This means that the max puzzle size we can support is 32 edges
type edges int32

const MAX_EDGES = 32

func newEdges() edges {
	return 0
}

func (e edges) isEdge(start int) bool {
	if start < 0 || start >= MAX_EDGES {
		// sanity check. we could probably remove for speed
		return false
	}
	return e&(1<<start) != 0
}

func (e edges) addEdge(start int) edges {
	if start < 0 || start >= MAX_EDGES {
		return e
	}
	return e | (1 << start)
}

type edgeDirection uint8

const (
	rowDir edgeDirection = 0
	colDir edgeDirection = 1
)

type edgesFromNode struct {
	above      int
	below      int
	left       int
	right      int
	totalEdges int

	isabove bool
	isbelow bool
	isleft  bool
	isright bool
}

func (efn edgesFromNode) hasBranch() bool {
	numBranches := 0

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

	return numBranches > 2
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

func (nt nodeType) isInvalidEdges(efn *edgesFromNode) bool {
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
	val   int

	straightLineLens int // a cached value of the current number of edges
	seen             bool
}

func (n node) copy() node {

	return node{
		nType: n.nType,
		val:   n.val,
	}
}

type grid struct {
	rows  []edges
	cols  []edges
	nodes map[row]map[col]node

	cachedEFNs [][]*edgesFromNode
}

func newGrid(
	size int,
	nodeLocations []NodeLocation,
) *grid {
	if size > MAX_EDGES {
		return nil
	}

	nodes := map[row]map[col]node{}
	for _, nl := range nodeLocations {
		if _, ok := nodes[row(nl.Row)]; !ok {
			nodes[row(nl.Row)] = make(map[col]node)
		}
		nt := blackNode
		if nl.IsWhite {
			nt = whiteNode
		}
		nodes[row(nl.Row)][col(nl.Col)] = node{
			nType: nt,
			val:   nl.Value,
		}
	}
	rows := make([]edges, size)
	for i := range rows {
		rows[i] = newEdges()
	}
	cols := make([]edges, size)
	for i := range cols {
		cols[i] = newEdges()
	}
	efnCache := make([][]*edgesFromNode, size)
	for i := range efnCache {
		efnCache[i] = make([]*edgesFromNode, size)
	}
	return &grid{
		rows:       rows,
		cols:       cols,
		nodes:      nodes,
		cachedEFNs: efnCache,
	}
}

func (g *grid) Copy() *grid {
	if g == nil {
		return nil
	}

	nodes := map[row]map[col]node{}
	for r, cMap := range g.nodes {
		nodes[r] = make(map[col]node, len(cMap))
		for c, n := range cMap {
			nodes[r][c] = n.copy()
		}
	}

	rows := make([]edges, len(g.rows))
	for i := range rows {
		rows[i] = g.rows[i]
	}

	cols := make([]edges, len(g.cols))
	for i := range cols {
		cols[i] = g.cols[i]
	}

	size := len(g.cachedEFNs)
	efnCache := make([][]*edgesFromNode, size)
	for i := range efnCache {
		efnCache[i] = make([]*edgesFromNode, size)
	}

	return &grid{
		rows:       rows,
		cols:       cols,
		nodes:      nodes,
		cachedEFNs: efnCache,
	}
}

func (g *grid) IsEdge(
	dir edgeDirection,
	eIndex, start int,
) bool {
	if eIndex < 0 {
		return false
	}

	var edges []edges

	switch dir {
	case rowDir:
		edges = g.rows
	case colDir:
		edges = g.cols
	}

	return eIndex < len(edges) && edges[eIndex].isEdge(start)
}

func (g *grid) AddEdge(
	dir edgeDirection,
	eIndex, start int,
) (*grid, error) {
	if start < 0 {
		return nil, errors.New(`negative start`)
	}
	if eIndex < 0 {
		return nil, errors.New(`negative eIndex`)
	}
	if g.IsEdge(dir, eIndex, start) {
		return nil, errors.New(`already edge`)
	}

	gCpy := g.Copy()

	switch dir {
	case rowDir:
		if eIndex >= len(gCpy.rows) {
			return nil, errors.New(`eIndex too large`)
		}
		gCpy.rows[eIndex] = gCpy.rows[eIndex].addEdge(start)
	case colDir:
		if eIndex >= len(gCpy.cols) {
			return nil, errors.New(`eIndex too large`)
		}
		gCpy.cols[eIndex] = gCpy.cols[eIndex].addEdge(start)
	}

	return gCpy, nil
}

func (g *grid) IsInvalid() bool {
	return g.isRangeInvalid(0, len(g.rows), 0, len(g.cols))
}

func (g *grid) isRangeInvalid(
	startR, stopR int,
	startC, stopC int,
) bool {
	if startR < 0 {
		startR = 0
	}
	if maxR := len(g.rows); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := len(g.cols); stopC > maxC {
		stopC = maxC
	}

	for r := startR; r < stopR; r++ {
		for c := startC; c < stopC; c++ {
			// check that this point doesn't branch
			efn := g.getEdgesFromNode(r, c)
			if efn.hasBranch() {
				// this is a branch.
				// therefore, this grid is invalid
				return true
			}

			// if this point is a node, check for if it's invalid
			if g.isInvalidNode(r, c, efn) {
				return true
			}
		}
	}

	return false
}

func (g *grid) getEdgesFromNode(
	r, c int,
) *edgesFromNode {
	if r < 0 || c < 0 {
		return nil
	}

	if g.cachedEFNs[r][c] != nil {
		return g.cachedEFNs[r][c]
	}

	above := g.getLenStraightPath(colDir, headUp, r, c)
	below := g.getLenStraightPath(colDir, headDown, r, c)
	left := g.getLenStraightPath(rowDir, headLeft, r, c)
	right := g.getLenStraightPath(rowDir, headRight, r, c)
	efn := edgesFromNode{
		totalEdges: above + below + left + right,
		above:      above,
		below:      below,
		left:       left,
		right:      right,
		isabove:    above > 0,
		isbelow:    below > 0,
		isleft:     left > 0,
		isright:    right > 0,
	}

	g.cachedEFNs[r][c] = &efn

	return &efn
}

func (g *grid) getLenStraightPath(
	dir edgeDirection,
	move cardinal,
	r, c int,
) int {
	var eIndex, start int

	switch dir {
	case rowDir:
		eIndex = r
		start = c
	case colDir:
		eIndex = c
		start = r
	}

	switch move {
	case headUp, headLeft:
		start--
	}

	numStraight := 0
	for i := start; g.IsEdge(dir, eIndex, i); {
		numStraight++

		switch move {
		case headUp, headLeft:
			i--
		default:
			i++
		}
	}
	return numStraight
}

func (g *grid) isInvalidNode(
	r, c int,
	efn *edgesFromNode,
) bool {
	n, ok := g.nodes[row(r)][col(c)]
	if !ok || n.nType == noNode {
		// no node == not invalid
		return false
	}

	// check that the node type rules are not broken
	if n.nType.isInvalidEdges(efn) {
		return true
	}

	// check that the num of straight line edges < n.val
	n.straightLineLens = efn.totalEdges
	return n.val < n.straightLineLens
}

func (g *grid) IsIncomplete() (bool, error) {
	if g.IsInvalid() {
		return true, errors.New(`invalid grid`)
	}

	firstR := -1
	firstC := -1
	for r := 0; r < len(g.rows) && firstR < 0; r++ {
		for c := 0; c < len(g.cols) && firstC < 0; c++ {
			hasRow := g.IsEdge(rowDir, r, c)
			hasCol := g.IsEdge(colDir, r, c)
			if hasRow || hasCol {
				if !hasRow || !hasCol {
					// don't need to walk the whole path if we see
					// from the start that it's not going to complete.
					return true, nil
				}
				// found firstEdge
				firstR = r
				firstC = c
			}
		}
	}
	// from the firstEdge, make your way around the grid until we get back
	// to the start. if we can't complete, then it's incomplete. as we see nodes,
	// mark them as seen
	nextR, nextC, move := g.walkToNextPoint(firstR, firstC, headUp)
	for move != headNowhere && nextR != firstR && nextC != firstC {
		nextR, nextC, move = g.walkToNextPoint(firstR, firstC, move)
	}
	if move == headNowhere {
		// our path all the way around was incomplete
		return true, nil
	}

	for _, cMap := range g.nodes {
		for _, n := range cMap {
			if n.straightLineLens != n.val || !n.seen {
				// somehow, we made a loop, but we didn't see all of the nodes
				// in the grid. therefore, this is incomplete.
				return true, nil
			}
		}
	}

	// at this point, we have a valid board, our path is a loop,
	// and we've seen all of the nodes appropriately. Therefore,
	// our board is not incomplete, and it's a solution.
	return false, nil
}

func (g *grid) walkToNextPoint(
	fromR, fromC int,
	avoid cardinal,
) (nextR, nextC int, _ cardinal) {
	efn := g.getEdgesFromNode(fromR, fromC)
	if efn == nil {
		return 0, 0, headNowhere
	}

	if efn.isabove && avoid != headUp {
		return fromR, fromC - efn.above, headUp
	}

	if efn.isleft && avoid != headLeft {
		return fromR - efn.left, fromC, headLeft
	}

	if efn.isbelow && avoid != headDown {
		return fromR, fromC + efn.below, headDown
	}

	if efn.isright && avoid != headRight {
		return fromR + efn.right, fromC, headRight
	}

	return 0, 0, headNowhere
}

func (g *grid) String() string {
	if g == nil {
		return `(*grid)<nil>`
	}
	var sb strings.Builder
	for r := range g.rows {
		var below strings.Builder
		for c := range g.cols {
			// write a node
			sb.WriteString(`(`)
			if n, ok := g.nodes[row(r)][col(c)]; ok {
				if n.nType == whiteNode {
					sb.WriteString(`w`)
				} else {
					sb.WriteString(`b`)
				}
				sb.WriteString(fmt.Sprintf("%2d", n.val))
			} else {
				sb.WriteString(`   `)
			}
			sb.WriteString(`)`)
			// now draw an edge?
			if g.IsEdge(rowDir, r, c) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			below.WriteString(`  `)
			if g.IsEdge(colDir, c, r) {
				below.WriteString(`|`)
			} else {
				below.WriteString(` `)
			}
			below.WriteString(`     `)
		}
		sb.WriteString("\n")
		sb.WriteString(below.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
