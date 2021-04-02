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

	sort.Slice(targets, func(i, j int) bool {
		// rank _lower_ valued nodes at the start of the Target list
		iPossibilities := numPossibleConfigs(targets[i], numEdges)
		jPossibilities := numPossibleConfigs(targets[j], numEdges)
		if iPossibilities != jPossibilities {
			// this is counter-intuitive to me. I would think that I should
			// solve for "big rocks" first. But it makes sense that a computer
			// can process all of the "size 2" nodes first, because they have
			// 2 solutions (for white) or 4 solutions (for black) instead of
			// the many possible solutions a larger node has. Coupling this with
			// the DFS search on targeted nodes provides marked improvements.
			return iPossibilities < jPossibilities
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

func numPossibleConfigs(
	target Target,
	numEdges int,
) int {
	twoArmsOptions := len(BuildTwoArmOptions(target.Node))

	if isOnTheSide(target.Coord, numEdges) {
		if isACorner(target.Coord, numEdges) {
			twoArmsOptions /= 4
		} else {
			twoArmsOptions /= 2
		}
	}

	return twoArmsOptions
}

func isOnTheSide(coord NodeCoord, numEdges int) bool {
	return coord.Row == 0 || coord.Row == RowIndex(numEdges) ||
		coord.Col == 0 || coord.Col == ColIndex(numEdges)
}

func isACorner(coord NodeCoord, numEdges int) bool {
	return (coord.Row == 0 || coord.Row == RowIndex(numEdges)) &&
		(coord.Col == 0 || coord.Col == ColIndex(numEdges))
}
