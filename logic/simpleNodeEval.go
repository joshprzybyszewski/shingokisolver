package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ evaluator = simpleNode{}

type simpleNode struct {
	nodeType model.NodeType
	oppEdge  model.EdgePair
	perps    []model.EdgePair
}

func newSimpleNodeEvaluator(
	node model.Node,
	otherSideOfNode model.EdgePair,
	perps []model.EdgePair,
) evaluator {
	if node.Type() == model.BlackNode {
		// black nodes don't care about perps.
		// don't retain a reference to them.
		perps = nil
	}

	return simpleNode{
		nodeType: node.Type(),
		oppEdge:  otherSideOfNode,
		perps:    perps,
	}
}

func (sn simpleNode) evaluate(ge model.GetEdger) model.EdgeState {
	switch sn.nodeType {
	case model.WhiteNode:
		// if the opposite edge is defined, then we are the same
		switch ge.GetEdge(sn.oppEdge) {
		case model.EdgeExists:
			return model.EdgeExists
		case model.EdgeAvoided, model.EdgeOutOfBounds:
			return model.EdgeAvoided
		}

		// if a perpendicular edge is defined, then we are the opposite
		for _, perpEP := range sn.perps {
			switch ge.GetEdge(perpEP) {
			case model.EdgeExists:
				return model.EdgeAvoided
			case model.EdgeAvoided, model.EdgeOutOfBounds:
				return model.EdgeExists
			}
		}

		// not enough info
		return model.EdgeUnknown

	case model.BlackNode:
		// if the opposite edge is defined, then we can know what we are
		switch ge.GetEdge(sn.oppEdge) {
		case model.EdgeExists:
			return model.EdgeAvoided
		case model.EdgeAvoided, model.EdgeOutOfBounds:
			return model.EdgeExists
		}

		// not enough info
		return model.EdgeUnknown

	default:
		// unsupported node type!
		return model.EdgeErrored
	}
}
