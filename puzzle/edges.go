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
	ep := model.NewEdgePair(nc, move)
	return p.GetEdgeState(ep) == model.EdgeExists
}

func (p *Puzzle) GetEdgeState(
	ep model.EdgePair,
) model.EdgeState {
	return p.edges.GetEdge(ep)
}

func (p *Puzzle) isEdgeDefined(ep model.EdgePair) bool {
	switch p.GetEdgeState(ep) {
	case model.EdgeAvoided, model.EdgeExists:
		return true
	}
	return false
}

func (p *Puzzle) AddEdge(
	startNode model.NodeCoord,
	move model.Cardinal,
) model.State {

	// p.printMsg("AddEdge(%s, %s)",
	// 	startNode,
	// 	move,
	// )

	return p.AddEdges(model.NewEdgePair(startNode, move))
}

func (p *Puzzle) AddEdges(
	pairs ...model.EdgePair,
) model.State {

	// p.printMsg("AddEdges(%+v)",
	// 	pairs,
	// )

	for _, ep := range pairs {
		if !p.edges.isInBounds(ep) {
			return model.Violation
		}

		switch s := p.addEdge(ep); s {
		case model.Incomplete, model.Complete, model.Duplicate:
		default:
			return s
		}
	}

	return p.runQueue()
}

func (p *Puzzle) addEdge(
	ep model.EdgePair,
) model.State {

	// p.printMsg("addEdge(%s)",
	// 	ep,
	// )

	switch state := p.edges.SetEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.noticeUpdated(ep)

		return p.checkRuleset(ep, model.EdgeExists)

	case model.Duplicate:
		return state
	default:
		// p.printMsg("addEdge(%s) edges.SetEdge returned %s",
		// 	ep,
		// 	state,
		// )
		return state
	}
}

func (p *Puzzle) AvoidEdge(
	ep model.EdgePair,
) model.State {
	// p.printMsg("AvoidEdge(%+v)",
	// 	ep,
	// )

	if !p.edges.isInBounds(ep) {
		return model.Violation
	}

	switch s := p.avoidEdge(ep); s {
	case model.Incomplete, model.Duplicate:
	default:
		return s
	}

	return p.runQueue()
}

func (p *Puzzle) avoidEdge(
	ep model.EdgePair,
) model.State {

	// p.printMsg("avoidEdge(%s)",
	// 	ep,
	// )

	switch state := p.edges.AvoidEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.noticeUpdated(ep)

		// see if I'm breaking any rules or I can make any more moves
		return p.checkRuleset(ep, model.EdgeAvoided)
	default:
		// p.printMsg("avoidEdge(%s) edges returned %s",
		// 	ep,
		// 	state,
		// )
		return state
	}
}

func (p *Puzzle) runQueue() model.State {
	defer p.rq.clearUpdated()

	for ep, ok := p.rq.pop(); ok; ep, ok = p.rq.pop() {
		switch s := p.updateEdgeFromRules(ep); s {
		case model.Violation,
			model.Unexpected:
			return s
		}
	}

	for ep := range p.rq.updated {
		eval := p.rules.getRules(ep).getEdgeState(p.edges)
		if eval == model.EdgeUnknown || eval == model.EdgeOutOfBounds {
			// this is ok. It means that our algorithm is trying out
			// edges, and we cannot determine what they are
			continue
		}

		exp := p.edges.GetEdge(ep)
		if eval != exp {
			// p.printMsg("runQueue(%s) evaled %s but expected %s",
			// 	ep,
			// 	eval,
			// 	exp,
			// )
			return model.Violation
		}
	}

	return model.Incomplete
}
