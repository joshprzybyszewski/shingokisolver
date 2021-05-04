package model

import "fmt"

type nodeOption struct {
	Options []TwoArms
	Node
}

// GetNextTarget will target the next node that should be the best for solving.
// Inputs:
//   curTarget is the target that was just solved for. InvalidTarget if this is the first one.
//   nodes and tas are index-matched slices, where tas is the options of
//     TwoArms for each node in nodes.
// Output:
//   The next target to solve for, with a backwards reference to curTarget via Parent.
//   A bool that is true when we've satisfied all of the nodes in the input.
//   An error if the input is invalid (this is unexpected).
func GetNextTarget(
	curTarget Target,
	nodes []Node,
	tas [][]TwoArms,
) (Target, bool, error) {

	var best nodeOption

	for i, n := range nodes {
		options := tas[i]
		if len(options) == 0 {
			// this means that there's a node with literally zero options
			return Target{}, false, fmt.Errorf(`no options for node %s`, n)
		}

		no := nodeOption{
			Node:    n,
			Options: options,
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
		Parent:  &curTarget,
	}, true, nil
}

func isBetterThanCurrentBest(
	prev, new nodeOption,
) bool {
	if prev.Options == nil {
		return true
	}

	if len(new.Options) != len(prev.Options) {
		return len(new.Options) < len(prev.Options)
	}

	if new.Node.Type() != prev.Node.Type() {
		return new.Node.Type() == WhiteNode
	}

	return new.Value() > prev.Value()
}
