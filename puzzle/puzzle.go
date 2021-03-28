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

	return &puzzle{
		numEdges:      uint8(numEdges),
		nodes:         nodes,
		outgoingEdges: newGridNoder(numEdges),
	}
}

func (p *puzzle) deepCopy() *puzzle {
	if p == nil {
		return nil
	}

	nodes := map[nodeCoord]node{}
	for nc, n := range p.nodes {
		nodes[nc] = n.copy()
	}

	return &puzzle{
		numEdges:      p.numEdges,
		nodes:         nodes,
		outgoingEdges: p.outgoingEdges.copy(),
	}
}

func (p *puzzle) numNodes() int {
	return int(p.numEdges) + 1
}

func (p *puzzle) IsEdge(
	move cardinal,
	nc nodeCoord,
) bool {
	isEdge, _ := p.isEdge(move, nc)
	return isEdge
}

func (p *puzzle) isEdge(
	move cardinal,
	nc nodeCoord,
) (bool, bool) {
	if !p.outgoingEdges.isInBounds(nc) {
		return false, false
	}
	maxIndex := p.numEdges

	switch move {
	case headLeft:
		return nc.col != 0 && p.outgoingEdges.get(nc).isleft(), true
	case headRight:
		return uint8(nc.col) != maxIndex && p.outgoingEdges.get(nc).isright(), true
	case headUp:
		return nc.row != 0 && p.outgoingEdges.get(nc).isabove(), true
	case headDown:
		return uint8(nc.row) != maxIndex && p.outgoingEdges.get(nc).isbelow(), true
	default:
		return false, false
	}
}

func (p *puzzle) AddEdge(
	move cardinal,
	startNode nodeCoord,
) (nodeCoord, *puzzle, error) {
	endNode := startNode.translate(move)
	if !p.outgoingEdges.isInBounds(endNode) {
		return nodeCoord{}, nil, errors.New(`next point is out of bounds`)
	}

	isEdge, isValid := p.isEdge(move, startNode)
	if !isValid {
		return nodeCoord{}, nil, errors.New(`invalid input`)
	}
	if isEdge {
		return nodeCoord{}, nil, errors.New(`already had edge`)
	}

	gCpy := p.deepCopy()
	updateConnections(gCpy.outgoingEdges, startNode, endNode, move)

	return endNode, gCpy, nil
}

func (p *puzzle) isRangeInvalidWithBoundsCheck(
	startR, stopR rowIndex,
	startC, stopC colIndex,
) bool {
	if startR < 0 {
		startR = 0
	}
	if maxR := rowIndex(p.numNodes()); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := colIndex(p.numNodes()); stopC > maxC {
		stopC = maxC
	}

	return p.isRangeInvalid(startR, stopR, startC, stopC)
}

func (p *puzzle) isRangeInvalid(
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
			oe, ok := p.getOutgoingEdgesFrom(nc)
			if !ok {
				// the coordinate must be out of bounds
				return true
			}
			if oe.getNumOutgoingDirections() > 2 {
				// either we can't get the node, or this is a branch.
				// therefore, this puzzle is invalid
				return true
			}

			// if this point is a node, check for if it's invalid
			if isInvalidNode(p, nc, oe) {
				return true
			}
		}
	}

	return false
}

func (p *puzzle) getOutgoingEdgesFrom(
	coord nodeCoord,
) (edgesFromNode, bool) {
	if !p.outgoingEdges.isInBounds(coord) {
		return edgesFromNode{}, false
	}

	return p.outgoingEdges.get(coord), true
}

func (p *puzzle) IsIncomplete(
	coord nodeCoord,
) (bool, error) {
	if p.isRangeInvalid(0, rowIndex(p.numNodes()), 0, colIndex(p.numNodes())) {
		return true, errors.New(`invalid puzzle`)
	}

	oe, ok := p.getOutgoingEdgesFrom(coord)
	if !ok {
		return true, errors.New(`bad input`)
	}

	if nOutgoingEdges := oe.getNumOutgoingDirections(); nOutgoingEdges == 0 {
		// we were given bad intel. let's find a node with any edges.
		found := false
		for r := rowIndex(0); int(r) < p.numNodes() && !found; r++ {
			for c := colIndex(0); int(c) < p.numNodes() && !found; c++ {
				coord = nodeCoord{
					row: r,
					col: c,
				}
				hasRow := p.IsEdge(headRight, coord)
				hasCol := p.IsEdge(headDown, coord)
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

	w := newWalker(p, coord)
	seenNodes, ok := w.walk()
	if !ok {
		// our path did not make it all the way around
		return true, nil
	}

	for nc, n := range p.nodes {
		if _, ok := seenNodes[nc]; !ok {
			// node was not seen
			return true, nil
		}

		oe, ok := p.getOutgoingEdgesFrom(nc)
		if !ok {
			return true, errors.New(`bad input`)
		}
		if oe.totalEdges() != n.val {
			// previously (in isRangeInvalid) we checked if oe.totalEdges() > n.val
			// This check exists to verify we have exactly how many we need.
			return true, nil
		}

	}

	// at this point, we have a valid board, our path is a loop,
	// and we've seen all of the nodes appropriately. Therefore,
	// our board is not incomplete, and it's a solution.
	return false, nil
}

func (p *puzzle) String() string {
	if p == nil {
		return `(*puzzle)<nil>`
	}
	var sb strings.Builder
	for r := 0; r < p.numNodes(); r++ {
		var below strings.Builder
		for c := 0; c < p.numNodes(); c++ {
			nc := nodeCoord{
				row: rowIndex(r),
				col: colIndex(c),
			}
			// write a node
			sb.WriteString(`(`)
			if n, ok := p.nodes[nc]; ok {
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
			if p.IsEdge(headRight, nc) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			// now draw any edges that are below
			below.WriteString(`  `)
			if p.IsEdge(headDown, nc) {
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
