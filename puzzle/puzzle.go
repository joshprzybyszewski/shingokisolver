package puzzle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type Puzzle struct {
	numEdges uint8
	nodes    map[model.NodeCoord]model.Node

	nodeGrid model.Grid
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) *Puzzle {
	if numEdges > model.MAX_EDGES {
		return nil
	}

	nodes := map[model.NodeCoord]model.Node{}
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		nodes[nc] = model.NewNode(nl.IsWhite, nl.Value)
	}

	return &Puzzle{
		numEdges: uint8(numEdges),
		nodes:    nodes,
		nodeGrid: model.NewGrid(numEdges),
	}
}

func (p *Puzzle) DeepCopy() *Puzzle {
	if p == nil {
		return nil
	}

	return &Puzzle{
		numEdges: p.numEdges,
		nodes:    p.nodes,
		nodeGrid: p.nodeGrid.Copy(),
	}
}

func (p *Puzzle) GetCoordForHighestValueNode() model.NodeCoord {
	var bestCoord model.NodeCoord
	bestVal := int8(-1)
	for nc, n := range p.nodes {
		if n.Value() > bestVal {
			bestCoord = nc
			bestVal = n.Value()
		}
	}
	return bestCoord
}

func (p *Puzzle) NodeTargets() map[model.NodeCoord]model.Node {
	return p.nodes
}

func (p *Puzzle) NumEdges() int {
	return int(p.numEdges)
}

func (p *Puzzle) numNodes() int {
	return int(p.numEdges) + 1
}

func (p *Puzzle) IsEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	isEdge, _ := p.isEdge(move, nc)
	return isEdge
}

func (p *Puzzle) isEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) (bool, bool) {
	if !p.nodeGrid.IsInBounds(nc) {
		return false, false
	}
	maxIndex := p.numEdges

	switch move {
	case model.HeadLeft:
		return nc.Col != 0 && p.nodeGrid.Get(nc).IsLeft(), true
	case model.HeadRight:
		return uint8(nc.Col) != maxIndex && p.nodeGrid.Get(nc).IsRight(), true
	case model.HeadUp:
		return nc.Row != 0 && p.nodeGrid.Get(nc).IsAbove(), true
	case model.HeadDown:
		return uint8(nc.Row) != maxIndex && p.nodeGrid.Get(nc).IsBelow(), true
	default:
		return false, false
	}
}

func (p *Puzzle) AddEdge(
	move model.Cardinal,
	startNode model.NodeCoord,
) (model.NodeCoord, *Puzzle, error) {
	endNode := startNode.Translate(move)
	if !p.nodeGrid.IsInBounds(endNode) {
		return model.NodeCoord{}, nil, errors.New(`next point is out of bounds`)
	}

	isEdge, isValid := p.isEdge(move, startNode)
	if !isValid {
		return model.NodeCoord{}, nil, errors.New(`invalid input`)
	}
	if isEdge {
		return model.NodeCoord{}, nil, errors.New(`already had edge`)
	}

	puzzCopy := p.DeepCopy()
	model.UpdateGridConnections(puzzCopy.nodeGrid, startNode, endNode, move)
	return endNode, puzzCopy, nil
}

func (p *Puzzle) IsRangeInvalid(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) bool {
	if startR < 0 {
		startR = 0
	}
	if maxR := model.RowIndex(p.numNodes()); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := model.ColIndex(p.numNodes()); stopC > maxC {
		stopC = maxC
	}

	return p.isRangeInvalid(startR, stopR, startC, stopC)
}

func (p *Puzzle) isRangeInvalid(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) bool {
	for r := startR; r < stopR; r++ {
		for c := startC; c < stopC; c++ {
			// check that this point doesn't branch
			nc := model.NewCoord(r, c)
			oe, ok := p.GetOutgoingEdgesFrom(nc)
			if !ok {
				// the coordinate must be out of bounds
				return true
			}
			if oe.GetNumOutgoingDirections() > 2 {
				// either we can't get the node, or this is a branch.
				// therefore, this Puzzle is invalid
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

func (p *Puzzle) GetOutgoingEdgesFrom(
	coord model.NodeCoord,
) (model.OutgoingEdges, bool) {
	if !p.nodeGrid.IsInBounds(coord) {
		return model.OutgoingEdges{}, false
	}

	return p.nodeGrid.Get(coord), true
}

func (p *Puzzle) IsIncomplete(
	coord model.NodeCoord,
) (bool, error) {
	if p.isRangeInvalid(0, model.RowIndex(p.numNodes()), 0, model.ColIndex(p.numNodes())) {
		return true, errors.New(`invalid Puzzle`)
	}

	oe, ok := p.GetOutgoingEdgesFrom(coord)
	if !ok {
		return true, errors.New(`bad input`)
	}

	if nOutgoingEdges := oe.GetNumOutgoingDirections(); nOutgoingEdges == 0 {
		// we were given bad intel. let's find a node with any edges.
		found := false
		for r := 0; r < p.numNodes() && !found; r++ {
			for c := 0; c < p.numNodes() && !found; c++ {
				coord = model.NewCoordFromInts(r, c)
				hasRow := p.IsEdge(model.HeadRight, coord)
				hasCol := p.IsEdge(model.HeadDown, coord)
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
			// this Puzzle had no edges!
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
			return true, errors.New(`this path made a loop, but didn't see every node`)
		}

		oe, ok := p.GetOutgoingEdgesFrom(nc)
		if !ok {
			return true, errors.New(`bad input`)
		}
		if oe.TotalEdges() != n.Value() {
			// previously (in isRangeInvalid) we checked if oe.TotalEdges() > n.val
			// This check exists to verify we have exactly how many we need.
			return true, nil
		}

	}

	// at this point, we have a valid board, our path is a loop,
	// and we've seen all of the nodes appropriately. Therefore,
	// our board is not incomplete, and it's a solution.
	return false, nil
}

func (p *Puzzle) String() string {
	if p == nil {
		return `(*Puzzle)<nil>`
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for r := 0; r < p.numNodes(); r++ {
		var below strings.Builder
		for c := 0; c < p.numNodes(); c++ {
			nc := model.NewCoordFromInts(r, c)
			// write a node
			sb.WriteString(`(`)
			if n, ok := p.nodes[nc]; ok {
				if n.Type() == model.WhiteNode {
					sb.WriteString(`w`)
				} else {
					sb.WriteString(`b`)
				}
				sb.WriteString(fmt.Sprintf("%2d", n.Value()))
			} else {
				sb.WriteString(`XXX`)
			}
			sb.WriteString(`)`)

			// now draw an edge
			if p.IsEdge(model.HeadRight, nc) {
				sb.WriteString(`---`)
			} else {
				sb.WriteString(`   `)
			}

			// now draw any edges that are below
			below.WriteString(`  `)
			if p.IsEdge(model.HeadDown, nc) {
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
