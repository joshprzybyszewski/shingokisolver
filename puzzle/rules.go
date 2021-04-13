package puzzle

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type rules struct {
	couldAffect []edgePair

	evals map[string]func(getEdger) model.EdgeState
}

func newRules(
	ep edgePair,
	numEdges int,
) *rules {

	r := rules{
		couldAffect: make([]edgePair, 0, 6),
		evals:       make(map[string]func(getEdger) model.EdgeState, 2),
	}

	otherStartEdges := getOtherEdgeInputs(ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ep.Translate(ep.Cardinal), ep.Opposite())

	r.couldAffect = append(r.couldAffect, otherStartEdges...)
	r.couldAffect = append(r.couldAffect, otherEndEdges...)

	r.evals[`otherNodeInputs`] = getStandardNodeRules(ep, otherStartEdges, otherEndEdges)

	return &r
}

func (r *rules) affects() []edgePair {
	return r.couldAffect
}

func (r *rules) getEdgeState(ge getEdger) model.EdgeState {
	if r == nil {
		return model.EdgeOutOfBounds
	}

	es := model.EdgeUnknown

	for _, eval := range r.evals {
		switch newES := eval(ge); newES {
		case model.EdgeErrored:
			return newES
		case model.EdgeAvoided, model.EdgeExists:
			if es != model.EdgeUnknown && es != newES {
				return model.EdgeErrored
			}
			es = newES
		case model.EdgeOutOfBounds, model.EdgeUnknown:
			// ok
		default:
			// unsupported response
			return model.EdgeErrored
		}
	}

	return es
}

func (r *rules) addRulesForNode(
	node model.Node,
	dir model.Cardinal,
) {
	if r == nil {
		// this node is actually out of bounds...
		return
	}

	otherSideOfNode := newEdgePair(
		node.Coord(),
		dir.Opposite(),
	)

	perps := make([]edgePair, 0, 2)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			newEdgePair(node.Coord(), perpDir),
		)
	}

	r.evals[fmt.Sprintf(`node(%s) edge(%s)`, node, otherSideOfNode)] = getNodeRules(node, otherSideOfNode, perps)
}

func getOtherEdgeInputs(
	coord model.NodeCoord,
	dir model.Cardinal,
) []edgePair {

	perps := make([]edgePair, 0, 3)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			newEdgePair(coord, perpDir),
		)
	}

	return append(perps, newEdgePair(coord, dir.Opposite()))
}
