package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) checkRuleset(
	rq *rulesQueue,
	ep edgePair,
	expState model.EdgeState,
) model.State {
	p.printMsg("checkRuleset(%s, %s)",
		ep, expState,
	)
	r := p.rules.getRules(ep)

	// check if the rules for this edge are broken
	newEdge := r.getEdgeState(p.edges)

	switch newEdge {
	case model.EdgeAvoided, model.EdgeExists:
		if expState != newEdge {
			p.printMsg("checkRuleset(%s, %s) previously was %s, now is %s",
				ep,
				expState,
				expState,
				newEdge,
			)
			return model.Violation
		}
	}

	// Now let's look at all of the other affected rules
	rq.push(r.affects()...)
	p.printMsg("checkRuleset(%s, %s) affects %+v",
		ep,
		expState,
		r.affects(),
	)

	return model.Incomplete
}

func (p *Puzzle) updateEdgeFromRules(
	rq *rulesQueue,
	ep edgePair,
) model.State {
	p.printMsg("updateEdgeFromRules(%s)", ep)

	switch es := p.rules.getRules(ep).getEdgeState(p.edges); es {
	case model.EdgeAvoided:
		return p.avoidEdge(rq, ep)
	case model.EdgeExists:
		return p.addEdge(rq, ep)
	case model.EdgeUnknown:
		return model.Incomplete
	case model.EdgeOutOfBounds:
		return model.Incomplete
	default:
		return model.Unexpected
	}
}
