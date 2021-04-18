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

	var best nodeOption

	for _, n := range nodes {
		if _, ok := seenNodes[n.Coord()]; ok {
			continue
		}
		options := getOptions(n)
		if len(options) == 0 {
			// this means that there's a node with literally zero options
			return Target{}, false, errors.New(`invalid node!`)
		}

		no := nodeOption{
			Node:    n,
			Options: options,
			MinDist: getMinDist(n, seenNodes),
		}

		if isBetterThanCurrentBest(best, no) {
			best = no
		}
	}

	if best.Options == nil {
		// all of the nodes have been satisfied.
		return Target{}, false, nil
	}

	return Target{
		Node:    best.Node,
		Options: best.Options,
		Parent:  curTarget,
	}, true, nil
}

func getMinDist(
	n Node,
	seen map[NodeCoord]struct{},
) int {
	lowest := -1

	for nc := range seen {
		dNew := nc.DistanceTo(n.Coord())
		if lowest == -1 || dNew < lowest {
			lowest = dNew
		}
	}

	return lowest
}

type nodeOption struct {
	Node
	Options []TwoArms
	MinDist int
}

func isBetterThanCurrentBest(
	prev, new nodeOption,
) bool {
	if prev.Options == nil {
		return true
	}

	if len(new.Options) == 1 {
		return true
	}

	if new.MinDist < prev.MinDist {
		return true
	}

	if len(new.Options) < len(prev.Options) {
		return true
	}

	if new.Node.Type() == WhiteNode {
		return true
	}

	return false
}
