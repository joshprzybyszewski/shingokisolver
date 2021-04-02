package solvers

import (
	"sort"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type target struct {
	coord model.NodeCoord
	node  model.Node

	next *target
}

func buildTargets(p *puzzle.Puzzle) []*target {
	targets := make([]*target, 0, len(p.NodeTargets()))

	for nc, n := range p.NodeTargets() {
		targets = append(targets, &target{
			node:  n,
			coord: nc,
		})
	}

	maxRowColVal := p.NumEdges()
	isOnTheSide := func(coord model.NodeCoord) bool {
		return coord.Row == 0 ||
			coord.Row == model.RowIndex(maxRowColVal) ||
			coord.Col == 0 ||
			coord.Col == model.ColIndex(maxRowColVal)
	}
	isACorner := func(coord model.NodeCoord) bool {
		return (coord.Row == 0 ||
			coord.Row == model.RowIndex(maxRowColVal)) &&
			(coord.Col == 0 ||
				coord.Col == model.ColIndex(maxRowColVal))
	}
	sort.Slice(targets, func(i, j int) bool {
		// rank _lower_ valued nodes at the start of the target list
		iPossibilities := possibleConfigurationsForNode(targets[i].node)
		jPossibilities := possibleConfigurationsForNode(targets[j].node)
		if iPossibilities != jPossibilities {
			// this is counter-intuitive to me. I would think that I should
			// solve for "big rocks" first. But it makes sense that a computer
			// can process all of the "size 2" nodes first, because they have
			// 2 solutions (for white) or 4 solutions (for black) instead of
			// the many possible solutions a larger node has. Coupling this with
			// the DFS search on targeted nodes provides marked improvements.
			return iPossibilities < jPossibilities
		}

		// put nodes with more limitations (i.e. on the sides or
		// the corners of the graph) higher up on the list
		iIsEdge := isOnTheSide(targets[i].coord)
		jIsEdge := isOnTheSide(targets[j].coord)
		if iIsEdge && jIsEdge {
			iIsACorner := isACorner(targets[i].coord)
			jIsACorner := isACorner(targets[j].coord)
			if iIsACorner && !jIsACorner {
				return true
			} else if !iIsACorner && jIsACorner {
				return false
			}
		} else if iIsEdge {
			return true
		} else if jIsEdge {
			return false
		}

		// put white nodes in front of black nodes
		// because white nodes are more restrictive than black nodes
		iIsWhite := targets[i].node.Type() == model.WhiteNode
		jIsWhite := targets[j].node.Type() == model.WhiteNode
		if iIsWhite != jIsWhite {
			return iIsWhite
		}

		// at this point, we just want a consistent ordering.
		// let's put nodes closer to (0,0) higher up in the list
		if targets[i].coord.Row != targets[j].coord.Row {
			return targets[i].coord.Row < targets[j].coord.Row
		}
		return targets[i].coord.Col < targets[j].coord.Col
	})

	for i := 1; i < len(targets); i++ {
		targets[i-1].next = targets[i]
	}

	return targets
}

func possibleConfigurationsForNode(
	n model.Node,
) int {
	switch n.Type() {
	case model.WhiteNode:
		return int(n.Value()-1) * 2
	case model.BlackNode:
		return int(n.Value()-1) * 4
	}

	return 0
}
