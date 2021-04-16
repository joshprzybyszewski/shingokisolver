package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) GetUnknownEdge() (model.EdgePair, bool) {
	// TODO find a loose end instead
	// we can do that by choosing a random node, then walking to the end
	// of it's path.
	for r := 0; r <= p.NumEdges(); r++ {
		for c := 0; c <= p.NumEdges(); c++ {
			nc := model.NewCoordFromInts(r, c)

			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}

			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}
		}
	}
	return model.EdgePair{}, false
}

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
	return p.AddEdges(model.NewEdgePair(startNode, move))
}

func (p *Puzzle) AddEdges(
	pairs ...model.EdgePair,
) model.State {
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
	switch state := p.edges.SetEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.noticeUpdated(ep)

		return p.checkRuleset(ep, model.EdgeExists)

	default:
		return state
	}
}

func (p *Puzzle) AvoidEdge(
	ep model.EdgePair,
) model.State {
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

	switch state := p.edges.AvoidEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.noticeUpdated(ep)

		// see if I'm breaking any rules or I can make any more moves
		return p.checkRuleset(ep, model.EdgeAvoided)
	default:
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
			return model.Violation
		}
	}

	return model.Incomplete
}
