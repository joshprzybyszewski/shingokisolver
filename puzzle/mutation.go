package puzzle

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

func AddEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	return AddEdges(p, []model.EdgePair{ep})
}

func AddTwoArms(
	p Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) (Puzzle, model.State) {

	return AddEdges(p, ta.GetAllEdges(start))
}

func AddEdges(
	p Puzzle,
	eps []model.EdgePair,
) (Puzzle, model.State) {
	return performUpdates(p, updates{
		edgesToAdd: eps,
	})
}

func AvoidEdge(
	p Puzzle,
	ep model.EdgePair,
) (Puzzle, model.State) {
	return performUpdates(p, updates{
		edgesToAvoid: []model.EdgePair{ep},
	})
}

func claimGimmes(
	p Puzzle,
	nodes []model.Node,
) (Puzzle, model.State) {
	obviousFilled, ms := performUpdates(p, updates{
		nodes: p.nodes,
	})

	switch ms {
	case model.Violation, model.Unexpected:
		log.Printf("ClaimGimmes() first performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	// now we're going to add all of the extended rules
	for _, n := range obviousFilled.nodes {
		allTAs := model.BuildTwoArmOptions(n, obviousFilled.NumEdges())
		nearbyNodes := model.BuildNearbyNodes(n, allTAs, obviousFilled)
		possibleTAs := n.GetFilteredOptions(allTAs, &obviousFilled.edges, nearbyNodes)
		obviousFilled.rules.AddAllTwoArmRules(
			n,
			obviousFilled,
			possibleTAs,
		)
	}

	return performUpdates(obviousFilled, updates{
		nodes: obviousFilled.nodes,
	})

}

type updates struct {
	edgesToAdd   []model.EdgePair
	edgesToAvoid []model.EdgePair

	nodes []model.Node
}

func performUpdates(
	p Puzzle,
	u updates,
) (Puzzle, model.State) {
	var ms model.State

	newState := p.edges.Copy()
	rq := logic.NewQueue(&newState, newState.NumEdges())
	rules := p.rules

	for _, ep := range u.edgesToAdd {
		if !newState.IsInBounds(ep) {
			// log.Printf("performUpdates(%+v) attempted to add an out-of-bounds edge: %v", u, ep)
			return Puzzle{}, model.Violation
		}

		ms = addEdge(&newState, ep, rq, rules)
		switch ms {
		case model.Violation, model.Unexpected:
			// TODO remove this
			// log.Printf("performUpdates(%+v) got bad state (%s) on addEdge: %v", u, ms, ep)
			return Puzzle{}, ms
		}
	}

	for _, ep := range u.edgesToAvoid {
		if !newState.IsInBounds(ep) {
			// log.Printf("performUpdates(%+v) attempted to avoid an out-of-bounds edge: %v", u, ep)
			return Puzzle{}, model.Violation
		}

		ms = avoidEdge(&newState, ep, rq, rules)
		switch ms {
		case model.Violation, model.Unexpected:
			// TODO remove this
			// log.Printf("performUpdates(%+v) got bad state (%s) on avoidEdge: %v", u, ms, ep)
			return Puzzle{}, ms
		}
	}

	for _, n := range u.nodes {
		for _, dir := range model.AllCardinals {
			ms = updateEdgeFromRules(
				&newState,
				model.NewEdgePair(n.Coord(), dir),
				rq,
				rules,
			)
			switch ms {
			case model.Violation, model.Unexpected:
				// log.Printf("performUpdates(%+v) got bad state (%s) on updateEdgeFromRules: %+v %s", u, ms, n, dir)
				return Puzzle{}, ms
			}
		}
	}

	ms = runQueue(&newState, rq, rules)
	switch ms {
	case model.Violation, model.Unexpected:
		// TODO remove this
		// log.Printf("performUpdates(%+v) got bad state on runQueue: %v", u, ms)
		return Puzzle{}, ms
	}

	return p.withNewState(newState), ms
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
		// log.Printf("avoidEdge(%+v) got bad state on AvoidEdge: %v", ep, ms)
		return ms
	}
}

func runQueue(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	defer rq.ClearUpdated()

	var ms model.State
	for ep, ok := rq.Pop(); ok; ep, ok = rq.Pop() {
		switch ms = updateEdgeFromRules(edges, ep, rq, rules); ms {
		case model.Violation, model.Unexpected:
			// log.Printf("runQueue(%+v) got bad state on updateEdgeFromRules: %v", ep, ms)
			return ms
		}
	}

	// TODO can I?
	return model.Incomplete
	for _, ep := range rq.Updated() {
		eval := rules.Get(ep).GetEvaluatedState(edges)
		if eval == model.EdgeUnknown || eval == model.EdgeOutOfBounds {
			// this is ok. It means that our algorithm is trying out
			// edges, and we cannot determine what they are
			continue
		}

		exp := edges.GetEdge(ep)
		if eval != exp {
			// log.Printf("runQueue(%+v, %+v) got bad state on GetEdge: %v", exp, eval, ms)
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
			// log.Printf("checkRuleset(%+v) got bad state on GetEvaluatedState: %s %s", ep, expState, newEdge)
			return model.Violation
		}
		// TODO this may have broken everything...
		// Now let's look at all of the other affected rules
		rq.Push(edges, r.Affects())
	}

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
