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
) (Puzzle, model.State) {
	allNodeEdgesToCheck := make(map[model.Node]map[model.Cardinal]int8, len(p.nodes))
	for _, n := range p.nodes {
		allNodeEdgesToCheck[n] = model.GetMaxArmsByDir(
			model.BuildTwoArmOptions(n, p.NumEdges()),
		)
	}
	obviousFilled, ms := performUpdates(p, updates{
		nodes: allNodeEdgesToCheck,
	})

	if ms != model.Incomplete {
		log.Printf("ClaimGimmes() first performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	// now we're going to add all of the extended rules
	allNodeEdgesToCheck = make(map[model.Node]map[model.Cardinal]int8, len(obviousFilled.nodes))
	for _, n := range obviousFilled.nodes {
		allTAs := model.BuildTwoArmOptions(n, obviousFilled.NumEdges())
		nearbyNodes := model.BuildNearbyNodes(n, allTAs, obviousFilled)
		possibleTAs := n.GetFilteredOptions(allTAs, &obviousFilled.edges, nearbyNodes)
		obviousFilled.rules.AddAllTwoArmRules(
			n,
			obviousFilled,
			possibleTAs,
		)
		allNodeEdgesToCheck[n] = model.GetMaxArmsByDir(possibleTAs)
	}

	return performUpdates(obviousFilled, updates{
		nodes: allNodeEdgesToCheck,
	})

}

type updates struct {
	edgesToAdd   []model.EdgePair
	edgesToAvoid []model.EdgePair

	nodes map[model.Node]map[model.Cardinal]int8
}

func performUpdates(
	p Puzzle,
	u updates,
) (Puzzle, model.State) {
	var ms model.State

	newState := p.edges.Copy()
	rq := logic.NewQueue(newState.NumEdges())
	rules := p.rules

	for _, ep := range u.edgesToAdd {
		if !newState.IsInBounds(ep) {
			return Puzzle{}, model.Violation
		}

		ms = addEdge(&newState, ep, rq, rules)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	for _, ep := range u.edgesToAvoid {
		if !newState.IsInBounds(ep) {
			return Puzzle{}, model.Violation
		}

		ms = avoidEdge(&newState, ep, rq, rules)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	for n, maxArmsByDir := range u.nodes {
		for dir, maxLen := range maxArmsByDir {
			ep := model.NewEdgePair(n.Coord(), dir)
			for i := int8(0); i <= maxLen; i++ {
				ms = updateEdgeFromRules(
					&newState,
					ep,
					rq,
					rules,
				)
				if ms != model.Incomplete && ms != model.Duplicate {
					return Puzzle{}, ms
				}
				ep = ep.Next(dir)
			}
		}
	}

	ms = runQueue(&newState, rq, rules)
	if ms != model.Incomplete {
		return Puzzle{}, ms
	}

	return p.withNewState(newState), ms
}

func runQueue(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {

	var ms model.State
	for ep, ok := rq.Pop(); ok; ep, ok = rq.Pop() {
		ms = updateEdgeFromRules(edges, ep, rq, rules)
		if ms != model.Incomplete && ms != model.Duplicate {
			return ms
		}
	}

	return model.Incomplete
}

func addEdge(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	return updateEdge(model.EdgeExists, ep, edges, rq, rules.Get(ep))
}

func avoidEdge(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	return updateEdge(model.EdgeAvoided, ep, edges, rq, rules.Get(ep))
}

func updateEdge(
	es model.EdgeState,
	ep model.EdgePair,
	edges *state.TriEdges,
	rq *logic.Queue,
	r *logic.Rules,
) model.State {

	if ms := edges.UpdateEdge(ep, es); ms != model.Incomplete {
		return ms
	}

	evalState := r.EvaluateFullState(edges)
	if evalState != model.EdgeUnknown && evalState != es {
		return model.Violation
	}

	rq.Push(edges, r.Affects())
	return model.Incomplete
}

func updateEdgeFromRules(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
) model.State {
	switch es := rules.Get(ep).EvaluateQuickState(edges); es {
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
