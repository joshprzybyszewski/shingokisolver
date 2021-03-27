package puzzle

import (
	"errors"
	"fmt"
	"strings"
)

type rowIndex int8
type colIndex int8

type nodeCoord struct {
	row rowIndex
	col colIndex
}

func (nc nodeCoord) translate(
	move cardinal,
) nodeCoord {
	switch move {
	case headUp:
		nc.row--
	case headDown:
		nc.row++
	case headLeft:
		nc.col--
	case headRight:
		nc.col++
	}
	return nc
}

type puzzle struct {
	numEdges uint8
	nodes    map[nodeCoord]node

	outgoingEdges gridNoder
}

func newGrid(
	numEdges int,
	nodeLocations []NodeLocation,
) *puzzle {
	if numEdges > MAX_EDGES {
		return nil
	}

	nodes := map[nodeCoord]node{}
	for _, nl := range nodeLocations {
		nt := blackNode
		if nl.IsWhite {
			nt = whiteNode
		}
		nodes[nodeCoord{
			row: rowIndex(nl.Row),
			col: colIndex(nl.Col),
		}] = node{
			nType: nt,
			val:   nl.Value,
		}
	}

	var outgoingEdges gridNoder
	switch numEdges {
	case 2:
		outgoingEdges = &grid3x3{}
	default:
		outgoingEdges = &maxGrid{}
	}

	return &puzzle{
		numEdges:      uint8(numEdges),
		nodes:         nodes,
		outgoingEdges: outgoingEdges,
	}
}

func (g *puzzle) deepCopy() *puzzle {
	if g == nil {
		return nil
	}

	nodes := map[nodeCoord]node{}
	for nc, n := range g.nodes {
		nodes[nc] = n.copy()
	}

	return &puzzle{
		numEdges:      g.numEdges,
		nodes:         nodes,
		outgoingEdges: g.outgoingEdges.copy(),
	}
}

func (g *puzzle) numNodes() int {
	return int(g.numEdges) + 1
}

func (g *puzzle) IsEdge(
	move cardinal,
	nc nodeCoord,
) bool {
	isEdge, _ := g.isEdge(move, nc)
	return isEdge
}

func (g *puzzle) isEdge(
	move cardinal,
	nc nodeCoord,
) (bool, bool) {
	if !g.isInBounds(nc) {
		return false, false
	}
	maxIndex := g.numEdges

	switch move {
	case headLeft:
		return nc.col != 0 && g.outgoingEdges.get(nc).isleft(), true
	case headRight:
		return uint8(nc.col) != maxIndex && g.outgoingEdges.get(nc).isright(), true
	case headUp:
		return nc.row != 0 && g.outgoingEdges.get(nc).isabove(), true
	case headDown:
		return uint8(nc.row) != maxIndex && g.outgoingEdges.get(nc).isbelow(), true
	default:
		return false, false
	}
}

func (g *puzzle) isInBounds(
	nc nodeCoord,
) bool {
	if nc.row < 0 || nc.col < 0 {
		return false
	}
	maxIndex := g.numEdges
	return uint8(nc.row) <= maxIndex && uint8(nc.col) <= maxIndex
}

func (g *puzzle) AddEdge(
	move cardinal,
	nc nodeCoord,
) (nodeCoord, *puzzle, error) {
	isEdge, isValid := g.isEdge(move, nc)
	if !isValid {
		return nodeCoord{}, nil, errors.New(`invalid input`)
	}
	if isEdge {
		return nodeCoord{}, nil, errors.New(`already had edge`)
	}

	newCoord := nc.translate(move)
	if !g.isInBounds(newCoord) {
		return nodeCoord{}, nil, errors.New(`next point is out of bounds`)
	}

	gCpy := g.deepCopy()
	gCpy.updateConnections(nc, move)

	return newCoord, gCpy, nil
}

func (g *puzzle) updateConnections(
	startCoord nodeCoord,
	motion cardinal,
) {
	start := g.outgoingEdges.get(startCoord)
	endCoord := startCoord.translate(motion)

	end := g.outgoingEdges.get(endCoord)

	switch motion {
	case headLeft:
		start.left = end.left + 1
		end.right = start.right + 1
	case headRight:
		start.right = end.right + 1
		end.left = start.left + 1
	case headUp:
		start.above = end.above + 1
		end.below = start.below + 1
	case headDown:
		start.below = end.below + 1
		end.above = start.above + 1
	}

	g.outgoingEdges.set(startCoord, start)
	g.outgoingEdges.set(endCoord, end)

	switch motion {
	case headLeft:
		g.updateRowConnections(endCoord, startCoord)
	case headRight:
		g.updateRowConnections(startCoord, endCoord)
	case headUp:
		g.updateColConnections(endCoord, startCoord)
	case headDown:
		g.updateColConnections(startCoord, endCoord)
	}
}

func (g *puzzle) updateRowConnections(
	leftNode, rightNode nodeCoord,
) {
	curCoord := rightNode
	cur := g.outgoingEdges.get(curCoord)
	for cur.isright() {
		nextCoord := curCoord.translate(headRight)
		if !g.isInBounds(nextCoord) {
			break
		}
		next := g.outgoingEdges.get(nextCoord)
		next.left = cur.left + 1
		g.outgoingEdges.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = leftNode
	cur = g.outgoingEdges.get(curCoord)
	for cur.isleft() {
		nextCoord := curCoord.translate(headLeft)
		if !g.isInBounds(nextCoord) {
			break
		}
		next := g.outgoingEdges.get(nextCoord)
		next.right = cur.right + 1
		g.outgoingEdges.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}

func (g *puzzle) updateColConnections(
	topNode, bottomNode nodeCoord,
) {
	curCoord := topNode
	cur := g.outgoingEdges.get(curCoord)
	for cur.isabove() {
		nextCoord := curCoord.translate(headUp)
		if !g.isInBounds(nextCoord) {
			break
		}
		next := g.outgoingEdges.get(nextCoord)
		next.below = cur.below + 1
		g.outgoingEdges.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = bottomNode
	cur = g.outgoingEdges.get(curCoord)
	for cur.isbelow() {
		nextCoord := curCoord.translate(headDown)
		if !g.isInBounds(nextCoord) {
			break
		}
		next := g.outgoingEdges.get(nextCoord)
		next.above = cur.above + 1
		g.outgoingEdges.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}

func (g *puzzle) isRangeInvalidWithBoundsCheck(
	startR, stopR rowIndex,
	startC, stopC colIndex,
) bool {
	if startR < 0 {
		startR = 0
	}
	if maxR := rowIndex(g.numNodes()); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := colIndex(g.numNodes()); stopC > maxC {
		stopC = maxC
	}

	return g.isRangeInvalid(startR, stopR, startC, stopC)
}

func (g *puzzle) isRangeInvalid(
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
			efn, ok := g.getEdgesFromNode(nc)
			if !ok {
				// the coordinate must be out of bounds
				return true
			}
			if efn.getNumOutgoingDirections() > 2 {
				// either we can't get the node, or this is a branch.
				// therefore, this puzzle is invalid
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

func (g *puzzle) getEdgesFromNode(
	coord nodeCoord,
) (edgesFromNode, bool) {
	if coord.row < 0 || coord.col < 0 {
		return edgesFromNode{}, false
	}

	if numNodes := g.numNodes(); int(coord.row) >= numNodes || int(coord.col) >= numNodes {
		return edgesFromNode{}, false
	}

	return g.outgoingEdges.get(coord), true
}

func (g *puzzle) isInvalidNode(
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
	return n.val < efn.totalEdges()
}

func (g *puzzle) IsIncomplete(
	coord nodeCoord,
) (bool, error) {
	if g.isRangeInvalid(0, rowIndex(g.numNodes()), 0, colIndex(g.numNodes())) {
		return true, errors.New(`invalid puzzle`)
	}

	efn, ok := g.getEdgesFromNode(coord)
	if !ok {
		return true, errors.New(`bad input`)
	}

	if nOutgoingEdges := efn.getNumOutgoingDirections(); nOutgoingEdges == 0 {
		// we were given bad intel. let's find a node with any edges.
		found := false
		for r := rowIndex(0); int(r) < g.numNodes() && !found; r++ {
			for c := colIndex(0); int(c) < g.numNodes() && !found; c++ {
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
			// this puzzle had no edges!
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
		efn, ok := g.getEdgesFromNode(nc)
		if !ok {
			return true, errors.New(`bad input`)
		}
		if sumStraightLines := efn.totalEdges(); sumStraightLines > n.val {
			// this node has too many edges coming out of it
			return true, errors.New(`this graph has too many connections for a node`)
		} else if sumStraightLines != n.val {
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

func (g *puzzle) String() string {
	if g == nil {
		return `(*puzzle)<nil>`
	}
	var sb strings.Builder
	for r := 0; r < g.numNodes(); r++ {
		var below strings.Builder
		for c := 0; c < g.numNodes(); c++ {
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
