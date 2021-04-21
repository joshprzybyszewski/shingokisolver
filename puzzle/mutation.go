package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

func AddEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	if !p.edges.IsInBounds(ep) {
		return Puzzle{}, model.Violation
	}

	var ms model.State

	// TODO make these structs if possible
	newState := p.edges.Copy()
	rq := logic.NewQueue(newState, newState.NumEdges())
	rules := p.rules

	ms = addEdge(newState, ep, rq, rules)
	switch ms {
	case model.Incomplete, model.Complete, model.Duplicate:
	default:
		return Puzzle{}, ms
	}

	ms = runQueue(newState, rq, rules)
	switch ms {
	case model.Incomplete, model.Complete, model.Duplicate:
		return p.withNewState(newState), ms
	default:
		return Puzzle{}, ms
	}
}

func AvoidEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	if !p.edges.IsInBounds(ep) {
		return Puzzle{}, model.Violation
	}

	// TODO make these structs if possible
	newState := p.edges.Copy()
	rq := logic.NewQueue(newState, newState.NumEdges())
	rules := p.rules

	ms := avoidEdge(newState, ep, rq, rules)
	switch ms {
	case model.Incomplete, model.Duplicate:
	default:
		return Puzzle{}, ms
	}

	ms = runQueue(newState, rq, rules)
	switch ms {
	case model.Incomplete, model.Complete, model.Duplicate:
		// TODO return a copy of p?
		return p.withNewState(newState), ms
	default:
		return Puzzle{}, ms
	}
}

func AddTwoArms(
	p Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) (Puzzle, model.State) {
	var ms model.State

	// TODO make these structs if possible
	newState := p.edges.Copy()
	rq := logic.NewQueue(newState, newState.NumEdges())
	rules := p.rules

	for _, ep := range ta.GetAllEdges(start) {
		if !newState.IsInBounds(ep) {
			return Puzzle{}, model.Violation
		}

		ms = addEdge(newState, ep, rq, rules)
		switch ms {
		case model.Incomplete, model.Complete, model.Duplicate:
		default:
			return Puzzle{}, ms
		}
	}

	ms = runQueue(newState, rq, rules)
	switch ms {
	case model.Incomplete, model.Complete, model.Duplicate:
		// TODO return a copy of p?
		return p.withNewState(newState), ms
	default:
		return Puzzle{}, ms
	}
}

func addEdge(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	switch ms := edges.SetEdge(ep); ms {
	case model.Incomplete, model.Complete:
		rq.NoticeUpdated(ep)

		// TODO return a copy of p?
		return checkRuleset(edges, ep, model.EdgeExists, rq, rules)

	default:
		return ms
	}
}

func avoidEdge(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {

	switch ms := edges.AvoidEdge(ep); ms {
	case model.Incomplete, model.Complete:
		rq.NoticeUpdated(ep)

		// see if I'm breaking any rules or I can make any more moves
		return checkRuleset(edges, ep, model.EdgeAvoided, rq, rules)
	default:
		return ms
	}
}

func runQueue(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	defer rq.ClearUpdated()

	for ep, ok := rq.Pop(); ok; ep, ok = rq.Pop() {
		switch s := updateEdgeFromRules(edges, ep, rq, rules); s {
		case model.Violation,
			model.Unexpected:
			return s
		}
	}

	for _, ep := range rq.Updated() {
		eval := rules.Get(ep).GetEvaluatedState(edges)
		if eval == model.EdgeUnknown || eval == model.EdgeOutOfBounds {
			// this is ok. It means that our algorithm is trying out
			// edges, and we cannot determine what they are
			continue
		}

		exp := edges.GetEdge(ep)
		if eval != exp {
			return model.Violation
		}
	}

	return model.Incomplete
}

func checkRuleset(
	edges *state.TriEdges,
	ep model.EdgePair,
	expState model.EdgeState,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	r := rules.Get(ep)

	// check if the rules for this edge are broken
	newEdge := r.GetEvaluatedState(edges)

	switch newEdge {
	case model.EdgeAvoided, model.EdgeExists:
		if expState != newEdge {
			return model.Violation
		}
	}

	// Now let's look at all of the other affected rules
	rq.Push(r.Affects()...)

	return model.Incomplete
}

func updateEdgeFromRules(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	switch es := rules.Get(ep).GetEvaluatedState(edges); es {
	case model.EdgeAvoided:
		return avoidEdge(edges, ep, rq, rules)
	case model.EdgeExists:
		return addEdge(edges, ep, rq, rules)
	case model.EdgeUnknown, model.EdgeOutOfBounds:
		return model.Incomplete
	default:
		return model.Unexpected
	}
}
