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
	ep, err := standardizeInput(nc, move)
	return err == nil && p.edges.GetEdge(ep) == model.EdgeExists
}

func (p *Puzzle) AddEdge(
	startNode model.NodeCoord,
	move model.Cardinal,
) model.State {

	p.printMsg("AddEdge(%s, %s)",
		startNode,
		move,
	)

	ep, err := standardizeInput(startNode, move)

	if err != nil {
		return model.Unexpected
	} else if !p.edges.isInBounds(ep) {
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
			// here
			return state
		}

	case model.Duplicate:
		// not technically updated, but we tried to.
		rq.noticeUpdated(ep)
		return model.Incomplete
	default:
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

	ep, err := standardizeInput(startNode, move)

	if err != nil {
		return model.Unexpected
	}

	rq := newRulesQueue()

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
		if p.rules.getRules(ep).getEdgeState(p.edges) != p.edges.GetEdge(ep) {
			return model.Violation
		}
	}

	return model.Incomplete
}
