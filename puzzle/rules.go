package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type rules struct {
	me EdgePair

	couldAffect []EdgePair

	evals []func(ge getEdger) model.EdgeState
}

func newRules(
	ep EdgePair,
	numEdges int,
) *rules {

	r := rules{
		me:          ep,
		couldAffect: make([]EdgePair, 0, 6),
		evals:       make([]func(ge getEdger) model.EdgeState, 0, 2),
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

func (r *rules) affects() []EdgePair {
	return r.couldAffect
}

func (r *rules) addAffected(couldAffect ...EdgePair) {
	if r == nil {
		return
	}
	for _, other := range couldAffect {
		if other.IsIn(r.couldAffect...) || other == r.me {
			continue
		}
		r.couldAffect = append(r.couldAffect, other)
	}
}

func (r *rules) addEvaluations(evals ...func(ge getEdger) model.EdgeState) {
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

var _ getEdger = (*edgeStateCache)(nil)

type edgeStateCache struct {
	ge    getEdger
	cache map[EdgePair]model.EdgeState
}

func newEdgeStateCache(ge getEdger) *edgeStateCache {
	return &edgeStateCache{
		ge:    ge,
		cache: make(map[EdgePair]model.EdgeState, 50),
	}
}

func (c *edgeStateCache) GetEdge(
	ep EdgePair,
) model.EdgeState {
	s, ok := c.cache[ep]
	if ok {
		return s
	}

	s = c.ge.GetEdge(ep)
	c.cache[ep] = s
	return s
}

func (r *rules) getEdgeState(ge getEdger) model.EdgeState {
	if r == nil {
		return model.EdgeOutOfBounds
	}

	geCache := newEdgeStateCache(ge)

	es := model.EdgeUnknown

	for i, eval := range r.evals {
		newES := eval(geCache)
		printDebugMsg(
			"r.evals[%d] chose %s",
			i,
			newES,
		)
		switch newES {
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
	printDebugMsg(
		"getEdgeState checked %d evals, and chose %s",
		len(r.evals),
		es,
	)

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

	r.addAffected(otherSideOfNode)
	r.addAffected(perps...)

	r.addEvaluations(
		getNodeRules(node, otherSideOfNode, perps)...,
	)
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
