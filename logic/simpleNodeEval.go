package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

func newSimpleNodeEvaluator(
	node model.Node,
	myDir model.Cardinal,
) evaluator {
	otherSideOfNode := model.NewEdgePair(
		node.Coord(),
		myDir.Opposite(),
	)

	switch node.Type() {
	case model.WhiteNode:
		perps := make([]model.EdgePair, 0, 2)
		for _, perpDir := range myDir.Perpendiculars() {
			perps = append(perps,
				model.NewEdgePair(node.Coord(), perpDir),
			)
		}

		return simpleWhiteNode{
			oppEdge: otherSideOfNode,
			perps:   perps,
		}
	case model.BlackNode:
		return simpleBlackNode{
			oppEdge: otherSideOfNode,
		}
	default:
		return errEval(`newSimpleNodeEvaluator received unexpected node type`)
	}
}

var _ evaluator = simpleWhiteNode{}

type simpleWhiteNode struct {
	perps   []model.EdgePair
	oppEdge model.EdgePair
}

func (sn simpleWhiteNode) evaluate(ge model.GetEdger) model.EdgeState {
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
}

var _ evaluator = simpleBlackNode{}

type simpleBlackNode struct {
	oppEdge model.EdgePair
}

func (sn simpleBlackNode) evaluate(ge model.GetEdger) model.EdgeState {
	// if the opposite edge is defined, then we can know what we are
	switch ge.GetEdge(sn.oppEdge) {
	case model.EdgeExists:
		return model.EdgeAvoided
	case model.EdgeAvoided, model.EdgeOutOfBounds:
		return model.EdgeExists
	}

	// not enough info
	return model.EdgeUnknown
}
