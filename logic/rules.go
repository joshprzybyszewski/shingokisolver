package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type evaluator interface {
	evaluate(model.GetEdger) model.EdgeState
}

type Rules struct {
	couldAffect map[model.EdgePair]struct{}

	evals []evaluator
	me    model.EdgePair
}

func newRules(
	ge model.GetEdger,
	ep model.EdgePair,
) *Rules {

	r := Rules{
		me:          ep,
		couldAffect: make(map[model.EdgePair]struct{}, 8),
		evals:       make([]evaluator, 0, 4),
	}

	otherStartEdges := getOtherEdgeInputs(ge, ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ge, ep.NodeCoord.Translate(ep.Cardinal), ep.Cardinal.Opposite())

	r.addAffected(otherStartEdges...)
	r.addEvaluation(newStandardInputEvaluator(otherStartEdges))

	r.addAffected(otherEndEdges...)
	r.addEvaluation(newStandardInputEvaluator(otherEndEdges))

	return &r
}

func (r *Rules) Affects() []model.EdgePair {
	// TODO this may be costly!
	res := make([]model.EdgePair, len(r.couldAffect))
	for ep := range r.couldAffect {
		res = append(res, ep)
	}
	return res
}

func (r *Rules) addAffected(otherEPs ...model.EdgePair) {
	if r == nil {
		return
	}
	for _, other := range otherEPs {
		if other == r.me {
			// don't allow self-reference, or multiple of
			// the same reference
			continue
		}
		r.couldAffect[other] = struct{}{}
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
		panic(`dev error: addEvaluation should not have been nil`)
	}

	r.evals = append(r.evals, eval)
}

func (r *Rules) GetEvaluatedState(ge model.GetEdger) model.EdgeState {
	if r == nil {
		return model.EdgeOutOfBounds
	}
	// log.Printf("GetEvaluatedState(%s)", r.me)
	es := model.EdgeUnknown

	for _, e := range r.evals {
		newES := e.evaluate(ge)
		// log.Printf("%s = %T %+v", newES, e, e)
		switch newES {
		case model.EdgeErrored:
			return newES
		case model.EdgeAvoided, model.EdgeExists:
			// return newES
			// TODO go back to this if I don't trust myself anymore
			if es != model.EdgeUnknown && es != newES {
				// log.Printf("(%s) previously evaluated: %s, now evaluated: %s\n\tCurrent evaluator: %T = %+v", r.me, es, newES, e, e)
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
	myDir model.Cardinal,
) {
	if r == nil {
		// this node is actually out of bounds...
		return
	}

	r.addEvaluation(
		newSimpleNodeEvaluator(node, myDir),
	)
}

func getOtherEdgeInputs(
	ge model.GetEdger,
	coord model.NodeCoord,
	myDir model.Cardinal,
) []model.EdgePair {

	otherInputs := make([]model.EdgePair, 0, 3)
	for dir := range model.AllCardinalsMap {
		if dir == myDir {
			continue
		}

		ep := model.NewEdgePair(coord, dir)
		if !ge.IsInBounds(ep) {
			continue
		}

		otherInputs = append(otherInputs, ep)
	}

	return otherInputs
}
