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
	return AddEdges(p, []model.EdgePair{ep})
}

func AddTwoArms(
	p Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) (Puzzle, model.State) {

	existing, avoided := ta.GetAllEdges(start)

	return performUpdates(p, updates{
		edgesToAdd:             existing,
		edgesToAvoid:           avoided,
		allowOutOfBoundsAvoids: true,
	})
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

type updates struct {
	nodes      map[model.Node]map[model.Cardinal]int8
	edgesToAdd []model.EdgePair

	edgesToAvoid           []model.EdgePair
	allowOutOfBoundsAvoids bool
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

		ms = addEdge(&newState, ep, rq, rules, p.areNodesComplete)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	for _, ep := range u.edgesToAvoid {
		if !newState.IsInBounds(ep) {
			if u.allowOutOfBoundsAvoids {
				continue
			}
			return Puzzle{}, model.Violation
		}

		ms = avoidEdge(&newState, ep, rq, rules, p.areNodesComplete)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	for n, maxArmsByDir := range u.nodes {
		for _, dir := range model.AllCardinals {
			dir := dir
			ep := model.NewEdgePair(n.Coord(), dir)

			maxLen := maxArmsByDir[dir]
			for i := int8(0); i <= maxLen; i++ {
				ms = updateEdgeFromRules(
					&newState,
					ep,
					rq,
					rules,
					p.areNodesComplete,
				)
				if ms != model.Incomplete && ms != model.Duplicate {
					return Puzzle{}, ms
				}
				ep = ep.Next(dir)
			}
		}
	}

	ms = runQueue(&newState, rq, rules, p.areNodesComplete)
	if ms != model.Incomplete {
		return Puzzle{}, ms
	}

	return p.withNewState(newState), ms
}

func runQueue(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
	areNodesComplete bool,
) model.State {

	var ms model.State
	for ep, ok := rq.Pop(); ok; ep, ok = rq.Pop() {
		ms = updateEdgeFromRules(edges, ep, rq, rules, areNodesComplete)
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
	areNodesComplete bool,
) model.State {
	return updateEdge(model.EdgeExists, ep, edges, rq, rules.Get(ep), areNodesComplete)
}

func avoidEdge(
	edges *state.TriEdges,
	ep model.EdgePair,
	rq *logic.Queue,
	rules *logic.RuleSet,
	areNodesComplete bool,
) model.State {
	return updateEdge(model.EdgeAvoided, ep, edges, rq, rules.Get(ep), areNodesComplete)
}

func updateEdge(
	es model.EdgeState,
	ep model.EdgePair,
	edges *state.TriEdges,
	rq *logic.Queue,
	r *logic.Rules,
	areNodesComplete bool,
) model.State {

	if ms := edges.UpdateEdge(ep, es); ms != model.Incomplete {
		return ms
	}

	// if NodesComplete, then I don't need to EvaluateFullState.
	// This will by-pass the extended "advanced node" and "simple node"
	// evaluators everywhere.

	var evalState model.EdgeState
	if areNodesComplete {
		evalState = r.EvaluateQuickState(edges)
	} else {
		evalState = r.EvaluateFullState(edges)
	}
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
	areNodesComplete bool,
) model.State {
	switch es := rules.Get(ep).EvaluateQuickState(edges); es {
	case model.EdgeAvoided:
		return avoidEdge(edges, ep, rq, rules, areNodesComplete)
	case model.EdgeExists:
		return addEdge(edges, ep, rq, rules, areNodesComplete)
	case model.EdgeUnknown, model.EdgeOutOfBounds:
		return model.Incomplete
	default:
		return model.Unexpected
	}
}
