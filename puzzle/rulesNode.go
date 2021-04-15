package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func getNodeRules(
	node model.Node,
	otherSideOfNode EdgePair,
	perps []EdgePair,
) []func(ge getEdger) model.EdgeState {
	return []func(ge getEdger) model.EdgeState{
		func(ge getEdger) model.EdgeState {
			oppState := ge.GetEdge(otherSideOfNode)
			switch node.Type() {
			case model.WhiteNode:
				// if the opposite edge is defined, then we are the same
				switch oppState {
				case model.EdgeExists:
					return model.EdgeExists
				case model.EdgeAvoided, model.EdgeOutOfBounds:
					return model.EdgeAvoided
				}

				// if a perpendicular edge is defined, then we are the opposite
				for _, perpEP := range perps {
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
				switch oppState {
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
		},
	}
}
