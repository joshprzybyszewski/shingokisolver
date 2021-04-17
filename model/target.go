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
	nodes []Node,
	getOptions func(Node) []TwoArms,
) (Target, bool, error) {
	seenNodes := make(map[NodeCoord]struct{}, len(nodes))
	for t := curTarget; t != nil; t = t.Parent {
		seenNodes[t.Node.Coord()] = struct{}{}
	}

	var bestNode Node
	var bestOptions []TwoArms

	for _, n := range nodes {
		if _, ok := seenNodes[n.Coord()]; ok {
			continue
		}
		options := getOptions(n)
		if len(options) == 0 {
			// this means that there's a node with literally zero options
			return Target{}, false, errors.New(`invalid node!`)
		}

		if len(options) < len(bestOptions) || bestOptions == nil {
			bestNode = n
			bestOptions = options
		}
	}

	if bestOptions == nil {
		// all of the nodes have been satisfied.
		return Target{}, false, nil
	}

	return Target{
		Node:    bestNode,
		Options: bestOptions,
		Parent:  curTarget,
	}, true, nil
}
