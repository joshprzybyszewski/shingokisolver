package puzzle

import (
	"errors"
	"fmt"
	"strings"
)

type row uint8
type col uint8

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

func (g *grid) deepCopy() *grid {
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
		return nil, errors.New(`AddEdge had negative input`)
	}

	var newEdges edges
	switch move {
	case headLeft, headRight:
		if r >= len(g.rows) {
			return nil, errors.New(`AddEdge had too many rows`)
		}
		orig := g.rows[r]
		startI := c
		if move == headLeft {
			startI = c - 1
		}
		if startI < 0 || startI >= g.numEdges {
			return nil, errors.New(`AddEdge had bad cols`)
		}
		newEdges = orig.addEdge(startI)
		if newEdges == orig {
			return nil, errors.New(`did not add edge`)
		}
	case headUp, headDown:
		if c >= len(g.cols) {
			return nil, errors.New(`AddEdge had too many cols`)
		}
		orig := g.cols[c]
		startI := r
		if move == headUp {
			startI = r - 1
		}
		if startI < 0 || startI >= g.numEdges {
			return nil, errors.New(`AddEdge had bad rows`)
		}
		newEdges = orig.addEdge(startI)
		if newEdges == orig {
			return nil, errors.New(`did not add edge`)
		}
	default:
		return nil, errors.New(`invalid cardinal`)
	}

	gCpy := g.deepCopy()
	switch move {
	case headLeft, headRight:
		gCpy.rows[r] = newEdges
	case headUp, headDown:
		gCpy.cols[c] = newEdges
	}

	return gCpy, nil
}

func (g *grid) isInvalid() bool {
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
		isabove:    above != 0,
		isbelow:    below != 0,
		isleft:     left != 0,
		isright:    right != 0,
	}

	g.cachedEFNs[r][c] = &efn

	return &efn
}

func (g *grid) getLenStraightPath(
	move cardinal,
	r, c int,
) int8 {
	eR, eC := r, c

	var numStraight int8
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
	if g.isInvalid() {
		return true, errors.New(`invalid grid`)
	}

	efn := g.getEdgesFromNode(firstR, firstC)
	if nCards := efn.getNumCardinals(); nCards == 0 {
		// we were given bad intel. let's find a node with any edges.
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
	} else if nCards != 2 {
		// don't need to walk the whole path if we see
		// from the start that it's not going to complete.
		return true, nil
	}

	// from the firstEdge, make your way around the grid until we get back
	// to the start. if we can't complete, then it's incomplete. as we see nodes,
	// mark them as seen
	curR, curC, move := g.walkToNextPoint(firstR, firstC, headDown)
	err := g.markNodesAsSeen(firstR, curR, firstC, curC)
	if err != nil {
		return true, err
	}
	var nextR, nextC int
	for curR != firstR || curC != firstC {
		nextR, nextC, move = g.walkToNextPoint(curR, curC, move)
		err = g.markNodesAsSeen(curR, nextR, curC, nextC)
		if err != nil {
			return true, err
		}
		if move == headNowhere || curR == nextR && curC == nextC {
			// check that we've completed the loop before checking the error
			// if we've complete the loop, then we don't care about the
			// error "nodes already seen"
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
			if efn == nil || !n.seen {
				// somehow, we made a loop, but we didn't see all of the nodes
				// in the grid. therefore, this is incomplete.
				return true, nil
			} else if efn.totalEdges > n.val {
				return true, errors.New(`this graph has too many connections for a node`)
			} else if efn.totalEdges != n.val {
				// TODO we could figure out if it's possible to add an edge to the end
				// of the current straight lines. If not, then we can return an error
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
		return fromR - int(efn.above), fromC, headDown
	}

	if efn.isleft && avoid != headLeft {
		return fromR, fromC - int(efn.left), headRight
	}

	if efn.isbelow && avoid != headDown {
		return fromR + int(efn.below), fromC, headUp
	}

	if efn.isright && avoid != headRight {
		return fromR, fromC + int(efn.right), headLeft
	}

	return 0, 0, headNowhere
}

func (g *grid) markNodesAsSeen(
	fromR, toR,
	fromC, toC int,
) error {
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
	for r, cMap := range g.nodes {
		for c, n := range cMap {
			if int(r) < minR || int(r) >= maxR ||
				int(c) < minC || int(c) >= maxC {
				continue
			}
			if n.seen {
				return errors.New(`already seen node`)
			}
			n.seen = true
		}
	}
	return nil
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

			// now draw an edge
			if g.IsEdge(headRight, r, c) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			// now draw any edges that are below
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
