package puzzle

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type rowIndex int8
type colIndex int8

type nodeCoord struct {
	row rowIndex
	col colIndex
}

type grid struct {
	rows  []edges
	cols  []edges
	nodes map[nodeCoord]*node

	cachedEFNs [][]edgesFromNode
	efnLock    sync.RWMutex
}

func newGrid(
	numEdges int,
	nodeLocations []NodeLocation,
) *grid {
	if numEdges > MAX_EDGES {
		return nil
	}

	nodes := map[nodeCoord]*node{}
	for _, nl := range nodeLocations {
		nt := blackNode
		if nl.IsWhite {
			nt = whiteNode
		}
		nodes[nodeCoord{
			row: rowIndex(nl.Row),
			col: colIndex(nl.Col),
		}] = &node{
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
	efnCache := make([][]edgesFromNode, numEdges+1)
	for i := range efnCache {
		efnCache[i] = make([]edgesFromNode, numEdges+1)
	}
	return &grid{
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

	nodes := map[nodeCoord]*node{}
	for nc, n := range g.nodes {
		nodes[nc] = n.copy()
	}

	rows := make([]edges, len(g.rows))
	for i := range rows {
		rows[i] = g.rows[i]
	}

	cols := make([]edges, len(g.cols))
	for i := range cols {
		cols[i] = g.cols[i]
	}

	efnCache := make([][]edgesFromNode, len(g.cachedEFNs))
	for i := range efnCache {
		efnCache[i] = make([]edgesFromNode, len(g.cachedEFNs[i]))
	}

	return &grid{
		rows:       rows,
		cols:       cols,
		nodes:      nodes,
		cachedEFNs: efnCache,
	}
}

func (g *grid) IsEdge(
	move cardinal,
	nc nodeCoord,
) bool {
	if nc.row < 0 || nc.col < 0 {
		return false
	}

	switch move {
	case headLeft:
		return int(nc.row) < len(g.rows) && int(nc.col) <= len(g.cols) && g.rows[nc.row].isEdge(edgeIndex(nc.col-1))
	case headRight:
		return int(nc.row) < len(g.rows) && int(nc.col) < len(g.cols) && g.rows[nc.row].isEdge(edgeIndex(nc.col))
	case headUp:
		return int(nc.col) < len(g.cols) && int(nc.row) <= len(g.rows) && g.cols[nc.col].isEdge(edgeIndex(nc.row-1))
	case headDown:
		return int(nc.col) < len(g.cols) && int(nc.row) < len(g.rows) && g.cols[nc.col].isEdge(edgeIndex(nc.row))
	default:
		return false
	}
}

func (g *grid) AddEdge(
	move cardinal,
	nc nodeCoord,
) (nodeCoord, *grid, error) {
	if nc.row < 0 || nc.col < 0 {
		return nodeCoord{}, nil, errors.New(`AddEdge had negative input`)
	}

	var newEdges edges

	switch move {
	case headLeft, headRight:
		if int(nc.row) >= len(g.rows) {
			return nodeCoord{}, nil, errors.New(`AddEdge had too many rows`)
		}
		orig := g.rows[nc.row]
		startI := nc.col
		if move == headLeft {
			startI = nc.col - 1
		}
		if startI < 0 || int(startI) >= g.numEdges() {
			return nodeCoord{}, nil, errors.New(`AddEdge had bad cols`)
		}
		newEdges = orig.addEdge(edgeIndex(startI))
		if newEdges == orig {
			return nodeCoord{}, nil, errors.New(`did not add edge`)
		}
	case headUp, headDown:
		if int(nc.col) >= len(g.cols) {
			return nodeCoord{}, nil, errors.New(`AddEdge had too many cols`)
		}
		orig := g.cols[nc.col]
		startI := nc.row
		if move == headUp {
			startI = nc.row - 1
		}
		if startI < 0 || int(startI) >= g.numEdges() {
			return nodeCoord{}, nil, errors.New(`AddEdge had bad rows`)
		}
		newEdges = orig.addEdge(edgeIndex(startI))
		if newEdges == orig {
			return nodeCoord{}, nil, errors.New(`did not add edge`)
		}
	default:
		return nodeCoord{}, nil, errors.New(`invalid cardinal`)
	}

	// TODO on copy, we could update the efn cache too...
	gCpy := g.deepCopy()
	newCoord := nc

	switch move {
	case headUp:
		newCoord.row--
		gCpy.cols[nc.col] = newEdges
	case headDown:
		newCoord.row++
		gCpy.cols[nc.col] = newEdges
	case headLeft:
		newCoord.col--
		gCpy.rows[nc.row] = newEdges
	case headRight:
		newCoord.col++
		gCpy.rows[nc.row] = newEdges
	}

	return newCoord, gCpy, nil
}

func (g *grid) numNodes() int {
	return len(g.cachedEFNs)
}

func (g *grid) numEdges() int {
	return len(g.cachedEFNs) - 1
}

func (g *grid) isRangeInvalidWithBoundsCheck(
	startR, stopR rowIndex,
	startC, stopC colIndex,
) bool {
	if startR < 0 {
		startR = 0
	}
	if maxR := rowIndex(len(g.rows)); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := colIndex(len(g.cols)); stopC > maxC {
		stopC = maxC
	}
	return g.isRangeInvalid(startR, stopR, startC, stopC)
}

func (g *grid) isRangeInvalid(
	startR, stopR rowIndex,
	startC, stopC colIndex,
) bool {
	for r := startR; r < stopR; r++ {
		for c := startC; c < stopC; c++ {
			// check that this point doesn't branch
			nc := nodeCoord{
				row: r,
				col: c,
			}
			efn := g.getEdgesFromNode(nc)
			if !efn.isPopulated || efn.getNumOutgoingDirections() > 2 {
				// either we can't get the node, or this is a branch.
				// therefore, this grid is invalid
				return true
			}

			// if this point is a node, check for if it's invalid
			if g.isInvalidNode(nc, efn) {
				return true
			}
		}
	}

	return false
}

func (g *grid) getEdgesFromNode(
	coord nodeCoord,
) edgesFromNode {
	if coord.row < 0 || coord.col < 0 {
		return edgesFromNode{}
	}
	if numNodes := g.numNodes(); int(coord.row) >= numNodes || int(coord.col) >= numNodes {
		return edgesFromNode{}
	}

	efn := g.readCacheForEdgesFromNode(coord)
	if efn.isPopulated {
		return efn
	}

	g.efnLock.Lock()
	defer g.efnLock.Unlock()

	efn = g.cachedEFNs[coord.row][coord.col]
	if efn.isPopulated {
		return efn
	}

	efn = newEdgesFromNode(
		g.findStraightLinePathLen(headUp, coord),
		g.findStraightLinePathLen(headDown, coord),
		g.findStraightLinePathLen(headLeft, coord),
		g.findStraightLinePathLen(headRight, coord),
	)

	g.cachedEFNs[coord.row][coord.col] = efn

	return efn
}

func (g *grid) readCacheForEdgesFromNode(
	coord nodeCoord,
) edgesFromNode {
	g.efnLock.RLock()
	defer g.efnLock.RUnlock()
	return g.cachedEFNs[coord.row][coord.col]
}

func (g *grid) findStraightLinePathLen(
	move cardinal,
	start nodeCoord,
) int8 {
	cur := start

	var numStraight int8
	for g.IsEdge(move, cur) {
		numStraight++

		switch move {
		case headUp:
			cur.row--
		case headDown:
			cur.row++
		case headLeft:
			cur.col--
		case headRight:
			cur.col++
		default:
			return numStraight
		}
	}
	return numStraight
}

func (g *grid) isInvalidNode(
	nc nodeCoord,
	efn edgesFromNode,
) bool {
	n, ok := g.nodes[nc]
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
	coord nodeCoord,
) (bool, error) {
	if g.isRangeInvalid(0, rowIndex(len(g.rows)), 0, colIndex(len(g.cols))) {
		return true, errors.New(`invalid grid`)
	}

	efn := g.getEdgesFromNode(coord)

	if nOutgoingEdges := efn.getNumOutgoingDirections(); !efn.isPopulated || nOutgoingEdges == 0 {
		// we were given bad intel. let's find a node with any edges.
		found := false
		for r := rowIndex(0); int(r) < len(g.rows) && !found; r++ {
			for c := colIndex(0); int(c) < len(g.cols) && !found; c++ {
				coord = nodeCoord{
					row: r,
					col: c,
				}
				hasRow := g.IsEdge(headRight, coord)
				hasCol := g.IsEdge(headDown, coord)
				if hasRow || hasCol {
					if !hasRow || !hasCol {
						// don't need to walk the whole path if we see
						// from the start that it's not going to complete.
						return true, nil
					}
					// found firstEdge
					found = true
				}
			}
		}
		if !found {
			// this grid had no edges!
			return true, nil
		}
	} else if nOutgoingEdges != 2 {
		// don't need to walk the whole path if we see
		// from the start that it's not going to complete.
		return true, nil
	}

	w := newWalker(g, coord)
	seenNodes, ok := w.walk()
	if !ok {
		// our path did not make it all the way around
		return true, nil
	}

	for nc, n := range g.nodes {
		efn := g.getEdgesFromNode(nc)
		if !efn.isPopulated {
			// somehow, we made a loop, but we didn't see all of the nodes
			// in the grid. therefore, this is incomplete.
			return true, nil
		} else if efn.totalEdges > n.val {
			// this node has too many edges coming out of it
			return true, errors.New(`this graph has too many connections for a node`)
		} else if efn.totalEdges != n.val {
			// this node has fewer than the number of edges it needs. Therefore,
			// we consider this graph incomplete
			return true, nil
		}

		if _, ok := seenNodes[nc]; !ok {
			// node was not seen
			return true, nil
		}
	}

	// at this point, we have a valid board, our path is a loop,
	// and we've seen all of the nodes appropriately. Therefore,
	// our board is not incomplete, and it's a solution.
	return false, nil
}

func (g *grid) String() string {
	if g == nil {
		return `(*grid)<nil>`
	}
	var sb strings.Builder
	for r := range g.cachedEFNs {
		var below strings.Builder
		for c := range g.cachedEFNs[r] {
			nc := nodeCoord{
				row: rowIndex(r),
				col: colIndex(c),
			}
			// write a node
			sb.WriteString(`(`)
			if n, ok := g.nodes[nc]; ok {
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
			if g.IsEdge(headRight, nc) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			// now draw any edges that are below
			below.WriteString(`  `)
			if g.IsEdge(headDown, nc) {
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
