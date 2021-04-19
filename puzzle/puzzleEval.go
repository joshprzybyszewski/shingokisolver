package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) checkRuleset(
	ep model.EdgePair,
	expState model.EdgeState,
) model.State {
	r := p.rules.getRules(ep)

	// check if the rules for this edge are broken
	newEdge := r.getEdgeState(p.edges)

	switch newEdge {
	case model.EdgeAvoided, model.EdgeExists:
		if expState != newEdge {
			return model.Violation
		}
	}

	// Now let's look at all of the other affected rules
	p.rq.push(r.affects()...)

	return model.Incomplete
}

func (p Puzzle) updateEdgeFromRules(
	ep model.EdgePair,
) model.State {
	switch es := p.rules.getRules(ep).getEdgeState(p.edges); es {
	case model.EdgeAvoided:
		return p.avoidEdge(ep)
	case model.EdgeExists:
		return p.addEdge(ep)
	case model.EdgeUnknown, model.EdgeOutOfBounds:
		return model.Incomplete
	default:
		return model.Unexpected
	}
}
