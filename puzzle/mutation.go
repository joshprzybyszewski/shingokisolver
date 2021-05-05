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
	metas []*model.NodeMeta

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

		ms = addEdge(&newState, ep, rq, rules)
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

		ms = avoidEdge(&newState, ep, rq, rules)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	var nms []*model.NodeMeta
	if u.metas != nil {
		nms = u.metas
		for _, nm := range u.metas {
			rq.PushNodes([]model.Node{nm.Node})
		}
	} else {
		nms = p.getMetasCopy()
	}

	ms = runQueue(&newState, rq, rules, nms, p.areNodesComplete())
	if ms != model.Incomplete {
		return Puzzle{}, ms
	}

	return p.withNewState(newState, nms), ms
}

func runQueue(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
	nms []*model.NodeMeta,
	areNodesComplete bool,
) model.State {

	var ms model.State

	for {

		// empty out the edge queue entirely
		for ep, ok := rq.Pop(); ok; ep, ok = rq.Pop() {
			ms = updateEdgeFromRules(edges, ep, rq, rules)
			if ms != model.Incomplete && ms != model.Duplicate {
				return ms
			}
		}
		if areNodesComplete {
			// no use checking nodes for when we know they're all complete.
			break
		}

		nodesToCheck := rq.PopAllNodes()
		if len(nodesToCheck) == 0 {
			break
		}

		// TODO this is n squared
		for _, n := range nodesToCheck {
			for _, nm := range nms {
				if nm.IsComplete || nm.Coord() != n.Coord() {
					continue
				}

				ms = checkAdvancedRules(edges, rq, rules, nm)
				if ms != model.Incomplete && ms != model.Duplicate {
					return ms
				}
				break // out of nms
			}
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

	evalState := r.EvaluateQuickState(edges)
	if evalState != model.EdgeUnknown && evalState != es {
		return model.Violation
	}

	rq.Push(edges, r.Affects())

	rq.PushNodes(r.InterestedNodes())

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
