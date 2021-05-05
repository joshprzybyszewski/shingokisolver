package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type evaluator interface {
	evaluate(model.GetEdger) model.EdgeState
}

type Rules struct {
	couldAffectMap map[model.EdgePair]struct{}
	couldAffect    []model.EdgePair

	interestedNodesSeener state.CoordSeener
	interestedNodes       []model.Node

	mustRunEvals []standardInput
	otherEvals   []evaluator

	me model.EdgePair
}

func newRules(
	ge model.GetEdger,
	ep model.EdgePair,
) *Rules {

	r := Rules{
		me:                    ep,
		couldAffectMap:        make(map[model.EdgePair]struct{}, 8),
		couldAffect:           make([]model.EdgePair, 8),
		interestedNodesSeener: state.NewCoordSeen(ge.NumEdges()),
		interestedNodes:       make([]model.Node, 8),
		mustRunEvals:          make([]standardInput, 0, 2),
		otherEvals:            make([]evaluator, 0, 4),
	}

	otherStartEdges := getOtherEdgeInputs(ge, ep.NodeCoord, ep.Cardinal)
	otherEndEdges := getOtherEdgeInputs(ge, ep.NodeCoord.Translate(ep.Cardinal), ep.Cardinal.Opposite())

	r.addAffected(otherStartEdges...)
	r.addEvaluation(newStandardInputEvaluator(otherStartEdges))

	r.addAffected(otherEndEdges...)
	r.addEvaluation(newStandardInputEvaluator(otherEndEdges))

	return &r
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

func (r *Rules) Affects() []model.EdgePair {
	return r.couldAffect
}

func (r *Rules) InterestedNodes() []model.Node {
	return r.interestedNodes
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
		if _, ok := r.couldAffectMap[other]; ok {
			continue
		}
		r.couldAffectMap[other] = struct{}{}
		r.couldAffect = append(r.couldAffect, other)
	}
}

func (r *Rules) addInterestedNode(node model.Node) {
	if r == nil {
		return
	}
	if r.interestedNodesSeener.IsCoordSeen(node.Coord()) {
		return
	}

	r.interestedNodesSeener.Mark(node.Coord())

	r.interestedNodes = append(r.interestedNodes, node)
}

func (r *Rules) addEvaluation(eval evaluator) {
	if r == nil {
		return
	}

	if eval == nil {
		panic(`dev error: addEvaluation should not have been nil`)
	}

	// type checking isn't necessarily the best decision, but
	// I chose to do it here because there's currently only
	// one type of evaluator that is a "must run". If there were
	// more kinds, I would add another method to the interface
	// called "must() bool" and then each evaluator could define
	// it if they think they're a must-run.
	if si, ok := eval.(standardInput); ok {
		r.mustRunEvals = append(r.mustRunEvals, si)
	} else {
		r.otherEvals = append(r.otherEvals, eval)
	}
}

func (r *Rules) EvaluateQuickState(ge model.GetEdger) model.EdgeState {
	return r.evaluateState(ge, false)
}

func (r *Rules) EvaluateFullState(ge model.GetEdger) model.EdgeState {
	return r.evaluateState(ge, true)
}

func (r *Rules) evaluateState(
	ge model.GetEdger,
	fullCheck bool,
) model.EdgeState {

	if r == nil {
		return model.EdgeOutOfBounds
	}

	// These are the rules that _must_ be evaluated every time.
	// They are the "standard input" rules that protect against branching.
	es := r.evaluateMusts(ge)

	if fullCheck || es == model.EdgeUnknown {
		// These rules are additional node-related rules on this edge.
		// They need not be evaluated all of the time, but they can
		// help detect when we know what an edge is.
		return r.evaluateAdditionals(ge, es, fullCheck)
	}

	return es
}

func (r *Rules) evaluateMusts(
	ge model.GetEdger,
) model.EdgeState {
	es := model.EdgeUnknown
	for _, e := range r.mustRunEvals {
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

func (r *Rules) evaluateAdditionals(
	ge model.GetEdger,
	found model.EdgeState,
	fullCheck bool,
) model.EdgeState {
	for _, e := range r.otherEvals {
		newES := e.evaluate(ge)
		switch newES {
		case model.EdgeErrored:
			return newES
		case model.EdgeAvoided, model.EdgeExists:
			if found != model.EdgeUnknown && found != newES {
				return model.EdgeErrored
			}
			if !fullCheck {
				return newES
			}
			found = newES
		case model.EdgeUnknown:
			// ok
		default:
			// unsupported response
			return model.EdgeErrored
		}
	}
	return found
}

func (r *Rules) addSimpleNodeRules(
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
