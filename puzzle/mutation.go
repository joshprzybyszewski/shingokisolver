package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func AddEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	var state model.State

	// TODO get rid of the need for this
	p = p.DeepCopy()

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

func AvoidEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	if !p.edges.IsInBounds(ep) {
		return Puzzle{}, model.Violation
	}

	// TODO get rid of the need for this
	p = p.DeepCopy()

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

func AddTwoArms(
	p Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) (Puzzle, model.State) {
	var state model.State

	// TODO get rid of the need for this
	p = p.DeepCopy()

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

func (p Puzzle) checkRuleset(
	ep model.EdgePair,
	expState model.EdgeState,
) model.State {
	r := p.rules.Get(ep)

	// check if the rules for this edge are broken
	newEdge := r.GetEvaluatedState(p.edges)

	switch newEdge {
	case model.EdgeAvoided, model.EdgeExists:
		if expState != newEdge {
			return model.Violation
		}
	}

	// Now let's look at all of the other affected rules
	p.rq.Push(r.Affects()...)

	return model.Incomplete
}

func (p Puzzle) updateEdgeFromRules(
	ep model.EdgePair,
) model.State {
	switch es := p.rules.Get(ep).GetEvaluatedState(p.edges); es {
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
