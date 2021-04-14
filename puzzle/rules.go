package puzzle

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type rules struct {
	couldAffect []EdgePair

	evals map[string]func(getEdger) model.EdgeState
}

func newRules(
	ep EdgePair,
	numEdges int,
) *rules {

	r := rules{
		couldAffect: make([]EdgePair, 0, 6),
		evals:       make(map[string]func(getEdger) model.EdgeState, 2),
	}

	otherStartEdges := getOtherEdgeInputs(ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ep.Translate(ep.Cardinal), ep.Opposite())

	r.couldAffect = append(r.couldAffect, otherStartEdges...)
	r.couldAffect = append(r.couldAffect, otherEndEdges...)

	r.evals[`otherNodeInputs`] = getStandardNodeRules(ep, otherStartEdges, otherEndEdges)

	return &r
}

func (r *rules) affects() []EdgePair {
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

	otherSideOfNode := NewEdgePair(
		node.Coord(),
		dir.Opposite(),
	)

	perps := make([]EdgePair, 0, 2)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			NewEdgePair(node.Coord(), perpDir),
		)
	}

	r.evals[fmt.Sprintf(`node(%s) edge(%s)`, node, otherSideOfNode)] = getNodeRules(node, otherSideOfNode, perps)
}

func getOtherEdgeInputs(
	coord model.NodeCoord,
	dir model.Cardinal,
) []EdgePair {

	perps := make([]EdgePair, 0, 3)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			NewEdgePair(coord, perpDir),
		)
	}

	return append(perps, NewEdgePair(coord, dir.Opposite()))
}

func (r *rules) addExtendedEval(
	node model.Node,
	ext extendedRules,
) {
	if r == nil {
		// this node is actually out of bounds...
		return
	}

	for _, other := range ext.couldAffect {
		if other.IsIn(r.couldAffect...) {
			continue
		}
		r.couldAffect = append(r.couldAffect, other)
	}

	for i, eval := range ext.evals {
		r.evals[fmt.Sprintf(`extEval[%d] from node(%s)`, i, node)] = eval
	}
}

type extendedRules struct {
	couldAffect []EdgePair
	evals       []func(getEdger) model.EdgeState
}
