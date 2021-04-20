package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) GetUnknownEdge() (model.EdgePair, bool) {

	// start from a point that is guaranteed to have an edge.
	// so we choose a node!
	var nc model.NodeCoord
	for _, n := range p.nodes {
		nc = n.Coord()
		break
	}

	// now let's walk to the end of the line
	w := newWalker(p.edges, nc)
	nc, isLoop := w.walkToTheEndOfThePath()
	if isLoop {
		// This is an error case. We made a loop, but we weren't expecting to.
		return model.EdgePair{}, false
	}

	// now from this end of the path, choose a random edge
	// off of it that is unknown.
	for dir := range model.AllCardinalsMap {
		ep := model.NewEdgePair(nc, dir)
		if !p.edges.IsDefined(ep) {
			return ep, true
		}
	}

	// we walked to the end of the path and did not find an unknown edge.
	// this is an error case.
	return model.EdgePair{}, false
}

func (p Puzzle) IsEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	ep := model.NewEdgePair(nc, move)
	return p.GetEdgeState(ep) == model.EdgeExists
}

func (p Puzzle) GetEdgeState(
	ep model.EdgePair,
) model.EdgeState {
	return p.edges.GetEdge(ep)
}

func (p Puzzle) AddEdge(
	ep model.EdgePair,
) (Puzzle, model.State) {
	var state model.State

	if !p.edges.IsInBounds(ep) {
		return Puzzle{}, model.Violation
	}

	state = p.addEdge(ep)
	switch state {
	case model.Incomplete, model.Complete, model.Duplicate:
	default:
		return Puzzle{}, state
	}

	state = p.runQueue()
	switch state {
	case model.Incomplete, model.Complete, model.Duplicate:
		// TODO return a copy of p?
		return p, state
	default:
		return Puzzle{}, state
	}
}

func (p Puzzle) AvoidEdge(
	ep model.EdgePair,
) (Puzzle, model.State) {
	if !p.edges.IsInBounds(ep) {
		return Puzzle{}, model.Violation
	}

	state := p.avoidEdge(ep)
	switch state {
	case model.Incomplete, model.Duplicate:
	default:
		return Puzzle{}, state
	}

	state = p.runQueue()
	switch state {
	case model.Incomplete, model.Complete, model.Duplicate:
		// TODO return a copy of p?
		return p, state
	default:
		return Puzzle{}, state
	}
}

func (p Puzzle) AddTwoArms(
	start model.NodeCoord,
	ta model.TwoArms,
) (Puzzle, model.State) {
	var state model.State

	for _, ep := range ta.GetAllEdges(start) {
		if !p.edges.IsInBounds(ep) {
			return Puzzle{}, model.Violation
		}

		state = p.addEdge(ep)
		switch state {
		case model.Incomplete, model.Complete, model.Duplicate:
		default:
			return Puzzle{}, state
		}
	}

	state = p.runQueue()
	switch state {
	case model.Incomplete, model.Complete, model.Duplicate:
		// TODO return a copy of p?
		return p, state
	default:
		return Puzzle{}, state
	}
}

func (p Puzzle) addEdge(
	ep model.EdgePair,
) model.State {
	switch state := p.edges.SetEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.NoticeUpdated(ep)

		// TODO return a copy of p?
		return p.checkRuleset(ep, model.EdgeExists)

	default:
		return state
	}
}

func (p Puzzle) avoidEdge(
	ep model.EdgePair,
) model.State {

	switch state := p.edges.AvoidEdge(ep); state {
	case model.Incomplete, model.Complete:
		p.rq.NoticeUpdated(ep)

		// TODO return a copy of p?
		// see if I'm breaking any rules or I can make any more moves
		return p.checkRuleset(ep, model.EdgeAvoided)
	default:
		return state
	}
}

func (p Puzzle) runQueue() model.State {
	defer p.rq.ClearUpdated()

	for ep, ok := p.rq.Pop(); ok; ep, ok = p.rq.Pop() {
		switch s := p.updateEdgeFromRules(ep); s {
		case model.Violation,
			model.Unexpected:
			return s
		}
	}

	for _, ep := range p.rq.Updated() {
		eval := p.rules.Get(ep).GetEvaluatedState(p.edges)
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
