package model

import (
	"errors"
)

type Target struct {
	Node    Node
	Options []TwoArms

	Parent *Target
}

func GetNextTarget(
	curTarget *Target,
	nodes map[NodeCoord]Node,
	getOptions func(Node) []TwoArms,
) (Target, bool, error) {
	seenNodes := make(map[NodeCoord]struct{}, len(nodes))
	for t := curTarget; t != nil; t = t.Parent {
		seenNodes[t.Node.Coord()] = struct{}{}
	}

	bestCoord := InvalidNodeCoord
	var bestOptions []TwoArms

	for nc, n := range nodes {
		if _, ok := seenNodes[nc]; ok {
			continue
		}
		options := getOptions(n)
		if len(options) == 0 {
			// this means that there's a node with literally zero options
			return Target{}, false, errors.New(`invalid node!`)
		}

		if len(options) < len(bestOptions) || bestOptions == nil {
			bestCoord = nc
			bestOptions = options
		}
	}

	if bestCoord == InvalidNodeCoord {
		// all of the nodes have been satisfied.
		return Target{}, false, nil
	}

	return Target{
		Node:    nodes[bestCoord],
		Options: bestOptions,
		Parent:  curTarget,
	}, true, nil
}
