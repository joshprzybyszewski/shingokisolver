package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type rules struct {
	me model.EdgePair

	couldAffect []model.EdgePair

	evals []func(ge model.GetEdger) model.EdgeState
}

func newRules(
	ep model.EdgePair,
	numEdges int,
) *rules {

	r := rules{
		me:          ep,
		couldAffect: make([]model.EdgePair, 0, 6),
		evals:       make([]func(ge model.GetEdger) model.EdgeState, 0, 2),
	}

	otherStartEdges := getOtherEdgeInputs(ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ep.Translate(ep.Cardinal), ep.Opposite())

	r.addAffected(otherStartEdges...)
	r.addAffected(otherEndEdges...)

	r.addEvaluations(getStandardNodeRules(
		ep,
		otherStartEdges,
		otherEndEdges,
	)...)

	return &r
}

func (r *rules) affects() []model.EdgePair {
	return r.couldAffect
}

func (r *rules) addAffected(couldAffect ...model.EdgePair) {
	if r == nil {
		return
	}
	for _, other := range couldAffect {
		if other == r.me || other.IsIn(r.couldAffect...) {
			continue
		}
		r.couldAffect = append(r.couldAffect, other)
	}
}

// TODO consider having a separate set of evals for node changes:#
func (r *rules) addEvaluations(evals ...func(ge model.GetEdger) model.EdgeState) {
	if r == nil {
		return
	}

	for _, eval := range evals {
		if eval == nil {
			panic(`dev error`)
		}
		r.evals = append(r.evals, eval)
	}
}

func (r *rules) getEdgeState(ge model.GetEdger) model.EdgeState {
	if r == nil {
		return model.EdgeOutOfBounds
	}

	es := model.EdgeUnknown

	for _, eval := range r.evals {
		newES := eval(ge)
		switch newES {
		case model.EdgeErrored:
			return newES
		case model.EdgeAvoided, model.EdgeExists:
			if es != model.EdgeUnknown && es != newES {
				return model.EdgeErrored
			}
			es = newES
		case model.EdgeUnknown:
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

	otherSideOfNode := model.NewEdgePair(
		node.Coord(),
		dir.Opposite(),
	)

	perps := make([]model.EdgePair, 0, 2)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			model.NewEdgePair(node.Coord(), perpDir),
		)
	}

	r.addAffected(otherSideOfNode)
	r.addAffected(perps...)

	r.addEvaluations(
		getSimpleNodeRule(node, otherSideOfNode, perps),
	)
}

func getOtherEdgeInputs(
	coord model.NodeCoord,
	dir model.Cardinal,
) []model.EdgePair {

	perps := make([]model.EdgePair, 0, 3)
	for _, perpDir := range dir.Perpendiculars() {
		perps = append(perps,
			model.NewEdgePair(coord, perpDir),
		)
	}

	return append(perps, model.NewEdgePair(coord, dir.Opposite()))
}
