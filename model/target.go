package model

import (
	"errors"
)

type Target struct {
	Node Node

	Parent *Target
}

func GetNextTarget(
	curTarget *Target,
	nodes map[NodeCoord]Node,
	getNumOptions func(Node) int,
) (Target, bool, error) {
	seenNodes := make(map[NodeCoord]struct{}, len(nodes))
	for t := curTarget; t != nil; t = t.Parent {
		seenNodes[t.Node.Coord()] = struct{}{}
	}

	bestCoord := InvalidNodeCoord
	bestVal := -1

	for nc, n := range nodes {
		if _, ok := seenNodes[nc]; ok {
			continue
		}
		nOptions := getNumOptions(n)
		if nOptions == 0 {
			// this means that there's a node with literally zero options
			return Target{}, false, errors.New(`invalid node!`)
		}

		if nOptions < bestVal || bestVal == -1 {
			bestCoord = nc
		}
	}

	if bestCoord == InvalidNodeCoord {
		// all of the nodes have been satisfied.
		return Target{}, false, nil
	}

	return Target{
		Node:   nodes[bestCoord],
		Parent: curTarget,
	}, true, nil
}
