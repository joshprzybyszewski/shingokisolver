package puzzle

import (
	"errors"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

var (
	ErrEdgeAlreadyExists = errors.New(`already had edge`)
)

func (p *Puzzle) IsEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	ep := newEdgePair(nc, move)
	return p.edges.GetEdge(ep) == model.EdgeExists
}

func (p *Puzzle) AddEdge(
	startNode model.NodeCoord,
	move model.Cardinal,
) model.State {

	p.printMsg("AddEdge(%s, %s)",
		startNode,
		move,
	)

	ep := newEdgePair(startNode, move)

	if !p.edges.isInBounds(ep) {
		return model.Violation
	}

	rq := newRulesQueue()

	switch s := p.addEdge(rq, ep); s {
	case model.Ok, model.Incomplete, model.Duplicate:
	default:
		return s
	}

	switch s := p.runQueue(rq); s {
	case model.Ok, model.Incomplete, model.Duplicate:
		return model.Incomplete
	default:
		return s
	}
}

func (p *Puzzle) addEdge(
	rq *rulesQueue,
	ep edgePair,
) model.State {

	p.printMsg("addEdge(%s)",
		ep,
	)

	switch state := p.edges.SetEdge(ep); state {
	case model.Ok, model.Incomplete, model.Complete:
		rq.noticeUpdated(ep)

		// see if I'm breaking any rules or I can make any more moves
		switch state := p.checkRuleset(rq, ep, model.EdgeExists); state {
		case model.Ok, model.Incomplete:
			return model.Incomplete
		default:
			p.printMsg("addEdge(%s) checkRuleset returned %s",
				ep,
				state,
			)
			return state
		}

	case model.Duplicate:
		// not technically updated, but we tried to.
		rq.noticeUpdated(ep)
		return model.Incomplete
	default:
		p.printMsg("addEdge(%s) edges.SetEdge returned %s",
			ep,
			state,
		)
		return state
	}

}

func (p *Puzzle) AvoidEdge(
	startNode model.NodeCoord,
	move model.Cardinal,
) model.State {

	p.printMsg("AvoidEdge(%s, %s)",
		startNode, move,
	)

	rq := newRulesQueue()
	ep := newEdgePair(startNode, move)

	switch s := p.avoidEdge(rq, ep); s {
	case model.Ok, model.Incomplete, model.Complete:
		return p.runQueue(rq)
	default:
		return s
	}

}

func (p *Puzzle) avoidEdge(
	rq *rulesQueue,
	ep edgePair,
) model.State {

	p.printMsg("avoidEdge(%s)",
		ep,
	)

	switch state := p.edges.AvoidEdge(ep); state {
	case model.Ok, model.Incomplete, model.Complete:
		rq.noticeUpdated(ep)

		// see if I'm breaking any rules or I can make any more moves
		return p.checkRuleset(rq, ep, model.EdgeAvoided)
	default:
		p.printMsg("avoidEdge(%s) edges returned %s",
			ep,
			state,
		)
		return state
	}
}

func (p *Puzzle) runQueue(
	rq *rulesQueue,
) model.State {
	for ep, ok := rq.pop(); ok; ep, ok = rq.pop() {
		switch s := p.updateEdgeFromRules(rq, ep); s {
		case model.Violation,
			model.Unexpected:
			return s
		}
	}

	for ep := range rq.updated {
		eval := p.rules.getRules(ep).getEdgeState(p.edges)
		if eval == model.EdgeUnknown || eval == model.EdgeOutOfBounds {
			// this is ok. It means that our algorithm is trying out
			// edges, and we cannot determine what they are
			continue
		}
		exp := p.edges.GetEdge(ep)
		if eval != exp {
			p.printMsg("runQueue(%s) evaled %s but expected %s",
				ep,
				eval,
				exp,
			)
			return model.Violation
		}
	}

	return model.Incomplete
}
