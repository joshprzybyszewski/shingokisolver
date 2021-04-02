package model

import (
	"sort"
)

type Target struct {
	Coord NodeCoord
	Node  Node

	Next *Target
}

func BuildTargets(
	nodes map[NodeCoord]Node,
	numEdges int,
) []Target {
	targets := make([]Target, 0, len(nodes))

	for nc, n := range nodes {
		targets = append(targets, Target{
			Node:  n,
			Coord: nc,
		})
	}

	isOnTheSide := func(coord NodeCoord) bool {
		return coord.Row == 0 ||
			coord.Row == RowIndex(numEdges) ||
			coord.Col == 0 ||
			coord.Col == ColIndex(numEdges)
	}
	isACorner := func(coord NodeCoord) bool {
		return (coord.Row == 0 ||
			coord.Row == RowIndex(numEdges)) &&
			(coord.Col == 0 ||
				coord.Col == ColIndex(numEdges))
	}
	sort.Slice(targets, func(i, j int) bool {
		// rank _lower_ valued nodes at the start of the Target list
		iPossibilities := possibleConfigurationsForNode(targets[i].Node)
		jPossibilities := possibleConfigurationsForNode(targets[j].Node)
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
		iIsEdge := isOnTheSide(targets[i].Coord)
		jIsEdge := isOnTheSide(targets[j].Coord)
		if iIsEdge && jIsEdge {
			iIsACorner := isACorner(targets[i].Coord)
			jIsACorner := isACorner(targets[j].Coord)
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
		iIsWhite := targets[i].Node.Type() == WhiteNode
		jIsWhite := targets[j].Node.Type() == WhiteNode
		if iIsWhite != jIsWhite {
			return iIsWhite
		}

		// at this point, we just want a consistent ordering.
		// let's put nodes closer to (0,0) higher up in the list
		if targets[i].Coord.Row != targets[j].Coord.Row {
			return targets[i].Coord.Row < targets[j].Coord.Row
		}
		return targets[i].Coord.Col < targets[j].Coord.Col
	})

	for i := 1; i < len(targets); i++ {
		targets[i-1].Next = &targets[i]
	}

	return targets
}

func possibleConfigurationsForNode(
	n Node,
) int {
	// TODO once we've got a cache on buildTwoArmOptions, use:
	// return len(buildTwoArmOptions(n))
	switch n.Type() {
	case WhiteNode:
		return int(n.Value()-1) * 2
	case BlackNode:
		return int(n.Value()-1) * 4
	}

	return 0
}
