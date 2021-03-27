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
	headRight   cardinal = 1
	headUp      cardinal = 2
	headLeft    cardinal = 3
	headDown    cardinal = 4
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

func (efn edgesFromNode) getNumCardinals() int {
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

	return numBranches
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

	seen bool
}

func (n *node) copy() *node {
	return &node{
		nType: n.nType,
		val:   n.val,
	}
}

type grid struct {
	numEdges int

	rows  []edges
	cols  []edges
	nodes map[row]map[col]*node

	cachedEFNs [][]*edgesFromNode
}

func newGrid(
	numEdges int,
	nodeLocations []NodeLocation,
) *grid {
	if numEdges > MAX_EDGES {
		return nil
	}

	nodes := map[row]map[col]*node{}
	for _, nl := range nodeLocations {
		if _, ok := nodes[row(nl.Row)]; !ok {
			nodes[row(nl.Row)] = make(map[col]*node)
		}
		nt := blackNode
		if nl.IsWhite {
			nt = whiteNode
		}
		nodes[row(nl.Row)][col(nl.Col)] = &node{
			nType: nt,
			val:   nl.Value,
		}
	}
	rows := make([]edges, numEdges+1)
	for i := range rows {
		rows[i] = newEdges()
	}
	cols := make([]edges, numEdges+1)
	for i := range cols {
		cols[i] = newEdges()
	}
	efnCache := make([][]*edgesFromNode, numEdges+1)
	for i := range efnCache {
		efnCache[i] = make([]*edgesFromNode, numEdges+1)
	}
	return &grid{
		numEdges:   numEdges,
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

	nodes := map[row]map[col]*node{}
	for r, cMap := range g.nodes {
		nodes[r] = make(map[col]*node, len(cMap))
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

	efnCache := make([][]*edgesFromNode, len(g.cachedEFNs))
	for i := range efnCache {
		efnCache[i] = make([]*edgesFromNode, len(g.cachedEFNs[i]))
	}

	return &grid{
		numEdges:   g.numEdges,
		rows:       rows,
		cols:       cols,
		nodes:      nodes,
		cachedEFNs: efnCache,
	}
}

func (g *grid) IsEdge(
	move cardinal,
	r, c int,
) bool {
	if r < 0 || c < 0 {
		return false
	}

	switch move {
	case headLeft:
		return r < len(g.rows) && c <= len(g.cols) && g.rows[r].isEdge(c-1)
	case headRight:
		return r < len(g.rows) && c < len(g.cols) && g.rows[r].isEdge(c)
	case headUp:
		return c < len(g.cols) && r <= len(g.rows) && g.cols[c].isEdge(r-1)
	case headDown:
		return c < len(g.cols) && r < len(g.rows) && g.cols[c].isEdge(r)
	default:
		return false
	}

}

func (g *grid) AddEdge(
	move cardinal,
	r, c int,
) (*grid, error) {
	if r < 0 || c < 0 {
		return nil, fmt.Errorf("invalid g.AddEdge(%d, %d, %d)", move, r, c)
	}

	var newEdges edges
	switch move {
	case headLeft, headRight:
		if r >= len(g.rows) {
			return nil, fmt.Errorf("r >= len(g.rows) g.AddEdge(%d, %d, %d)", move, r, c)
		}
		orig := g.rows[r]
		startI := c
		if move == headLeft {
			startI = c - 1
		}
		if startI < 0 || startI >= g.numEdges {
			return nil, fmt.Errorf("invalid col g.AddEdge(%d, %d, %d)", move, r, c)
		}
		newEdges = orig.addEdge(startI)
		if newEdges == orig {
			return nil, errors.New(`did not add edge`)
		}
	case headUp, headDown:
		if c >= len(g.cols) {
			return nil, fmt.Errorf("c >= len(g.cols) g.AddEdge(%d, %d, %d)", move, r, c)
		}
		orig := g.cols[c]
		startI := r
		if move == headUp {
			startI = r - 1
		}
		if startI < 0 || startI >= g.numEdges {
			return nil, fmt.Errorf("invalid row g.AddEdge(%d, %d, %d)", move, r, c)
		}
		newEdges = orig.addEdge(startI)
		if newEdges == orig {
			return nil, errors.New(`did not add edge`)
		}
	default:
		return nil, errors.New(`invalid cardinal`)
	}

	gCpy := g.Copy()
	switch move {
	case headLeft, headRight:
		gCpy.rows[r] = newEdges
	case headUp, headDown:
		gCpy.cols[c] = newEdges
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
			if efn.getNumCardinals() > 2 {
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
	if numNodes := len(g.cachedEFNs); r >= numNodes || c >= numNodes {
		return nil
	}

	if g.cachedEFNs[r][c] != nil {
		return g.cachedEFNs[r][c]
	}

	above := g.getLenStraightPath(headUp, r, c)
	below := g.getLenStraightPath(headDown, r, c)
	left := g.getLenStraightPath(headLeft, r, c)
	right := g.getLenStraightPath(headRight, r, c)
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
	move cardinal,
	r, c int,
) int {
	eR, eC := r, c

	numStraight := 0
	for g.IsEdge(move, eR, eC) {
		numStraight++

		switch move {
		case headUp:
			eR--
		case headDown:
			eR++
		case headLeft:
			eC--
		case headRight:
			eC++
		default:
			return numStraight
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
	return n.val < efn.totalEdges
}

func (g *grid) IsIncomplete(
	firstR, firstC int,
) (bool, error) {
	if g.IsInvalid() {
		return true, errors.New(`invalid grid`)
	}

	efn := g.getEdgesFromNode(firstR, firstC)
	if efn.getNumCardinals() != 2 {
		// don't need to walk the whole path if we see
		// from the start that it's not going to complete.
		return true, nil
	} else {
		// we were given bad intel. let's find a node.
		firstR, firstC = -1, -1
		for r := 0; r < len(g.rows) && firstR < 0; r++ {
			for c := 0; c < len(g.cols) && firstC < 0; c++ {
				hasRow := g.IsEdge(headRight, r, c)
				hasCol := g.IsEdge(headDown, r, c)
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
	}

	// from the firstEdge, make your way around the grid until we get back
	// to the start. if we can't complete, then it's incomplete. as we see nodes,
	// mark them as seen
	curR, curC, move := g.walkToNextPoint(firstR, firstC, headDown)
	g.markNodesAsSeen(firstR, curR, firstC, curC)
	var nextR, nextC int
	for curR != firstR || curC != firstC {
		nextR, nextC, move = g.walkToNextPoint(curR, curC, move)
		g.markNodesAsSeen(curR, nextR, curC, nextC)
		if move == headNowhere || curR == nextR && curC == nextC {
			break
		}
		curR, curC = nextR, nextC
	}
	if move == headNowhere {
		// our path all the way around was incomplete
		return true, nil
	}

	for r, cMap := range g.nodes {
		for c, n := range cMap {
			efn := g.getEdgesFromNode(int(r), int(c))
			if efn == nil || efn.totalEdges != n.val || !n.seen {
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
		return fromR - efn.above, fromC, headDown
	}

	if efn.isleft && avoid != headLeft {
		return fromR, fromC - efn.left, headRight
	}

	if efn.isbelow && avoid != headDown {
		return fromR + efn.below, fromC, headUp
	}

	if efn.isright && avoid != headRight {
		return fromR, fromC + efn.right, headLeft
	}

	return 0, 0, headNowhere
}

func (g *grid) markNodesAsSeen(
	fromR, toR,
	fromC, toC int,
) {
	minR := fromR
	maxR := toR
	if maxR < minR {
		minR, maxR = maxR, minR
	}
	minC := fromC
	maxC := toC
	if maxC < minC {
		minC, maxC = maxC, minC
	}
	for r := minR; r <= maxR; r++ {
		for c := minC; c <= maxC; c++ {
			if n, ok := g.nodes[row(r)][col(c)]; ok {
				n.seen = true
			}
		}
	}
}

func (g *grid) String() string {
	if g == nil {
		return `(*grid)<nil>`
	}
	var sb strings.Builder
	for r := range g.cachedEFNs {
		var below strings.Builder
		for c := range g.cachedEFNs[r] {
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
				sb.WriteString(`XXX`)
			}
			sb.WriteString(`)`)
			// now draw an edge?
			if g.IsEdge(headRight, r, c) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			below.WriteString(`  `)
			if g.IsEdge(headDown, r, c) {
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
