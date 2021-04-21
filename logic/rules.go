package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type evaluator interface {
	evaluate(model.GetEdger) model.EdgeState
}

type Rules struct {
	couldAffect []model.EdgePair

	evals []evaluator
	me    model.EdgePair
}

func newRules(
	ep model.EdgePair,
) *Rules {

	r := Rules{
		me:          ep,
		couldAffect: make([]model.EdgePair, 0, 8),
		evals:       make([]evaluator, 0, 4),
	}

	otherStartEdges := getOtherEdgeInputs(ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ep.Translate(ep.Cardinal), ep.Opposite())

	r.addAffected(otherStartEdges...)
	r.addAffected(otherEndEdges...)

	r.addEvaluation(
		standardInput{
			otherInputs: otherStartEdges,
		},
	)
	r.addEvaluation(
		standardInput{
			otherInputs: otherEndEdges,
		},
	)

	return &r
}

func (r *Rules) Affects() []model.EdgePair {
	return r.couldAffect
}

func (r *Rules) addAffected(couldAffect ...model.EdgePair) {
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

// TODO instead of relying on the execution of evals of other
// nodes _after_ these rules have been checked, we should detect
// what other nodes change when "I" go into the Exists/Avoided state.
func (r *Rules) addEvaluation(eval evaluator) {
	if r == nil {
		return
	}

	if eval == nil {
		panic(`dev error`)
	}
	r.evals = append(r.evals, eval)
}

func (r *Rules) GetEvaluatedState(ge model.GetEdger) model.EdgeState {
	if r == nil {
		return model.EdgeOutOfBounds
	}

	es := model.EdgeUnknown

	for _, e := range r.evals {
		newES := e.evaluate(ge)
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

func (r *Rules) addRulesForNode(
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

	r.addEvaluation(
		newSimpleNodeEvaluator(node, otherSideOfNode, perps),
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
