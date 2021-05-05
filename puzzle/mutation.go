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

	var nms []*model.NodeMeta
	if u.metas != nil {
		nms = u.metas
		for _, nm := range u.metas {
			rq.PushNodes([]model.Node{nm.Node})
		}
	} else {
		nms = p.getMetasCopy()
	}

	su := newStateUpdater(
		&newState,
		rq,
		p.rules,
		nms,
		p.areNodesComplete(),
	)

	for _, ep := range u.edgesToAdd {
		if !newState.IsInBounds(ep) {
			return Puzzle{}, model.Violation
		}

		ms = su.addEdge(ep)
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

		ms = su.avoidEdge(ep)
		if ms != model.Incomplete && ms != model.Duplicate {
			return Puzzle{}, ms
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered: %+v", r)
			log.Printf("old puzzle: %s", p.String())
			log.Printf("new puzzle: %s", p.withNewState(newState, nms).String())
			panic(r)
		}
	}()

	ms = su.runQueue()
	if ms != model.Incomplete {
		return Puzzle{}, ms
	}

	return p.withNewState(newState, nms), ms
}

type stateUpdater struct {
	edges *state.TriEdges
	rq    *logic.Queue
	rules *logic.RuleSet
	nms   []*model.NodeMeta

	areNodesComplete bool
}

func newStateUpdater(
	edges *state.TriEdges,
	rq *logic.Queue,
	rules *logic.RuleSet,
	nms []*model.NodeMeta,
	areNodesComplete bool,
) stateUpdater {
	return stateUpdater{
		edges:            edges,
		rq:               rq,
		rules:            rules,
		nms:              nms,
		areNodesComplete: areNodesComplete,
	}
}

func (su stateUpdater) runQueue() model.State {

	var ms model.State

	for {

		// empty out the edge queue entirely
		for ep, ok := su.rq.Pop(); ok; ep, ok = su.rq.Pop() {
			ms = su.updateEdgeFromRules(ep)
			if ms != model.Incomplete && ms != model.Duplicate {
				return ms
			}
		}
		if su.areNodesComplete {
			// no use checking nodes for when we know they're all complete.
			break
		}

		nodesToCheck := su.rq.PopAllNodes()
		if len(nodesToCheck) == 0 {
			break
		}

		// TODO this is n squared
		for _, n := range nodesToCheck {
			ms = su.checkAdvancedRules(n)
			switch ms {
			case model.Complete:
				// TODO check for a loop?
			case model.Incomplete, model.Duplicate:
				// just keep going
			default:
				return ms
			}
		}
	}

	return model.Incomplete
}

func (su stateUpdater) addEdge(
	ep model.EdgePair,
) model.State {
	return su.updateEdge(ep, model.EdgeExists)
}

func (su stateUpdater) avoidEdge(
	ep model.EdgePair,
) model.State {
	return su.updateEdge(ep, model.EdgeAvoided)
}

func (su stateUpdater) updateEdge(
	ep model.EdgePair,
	es model.EdgeState,
) model.State {

	if ms := su.edges.UpdateEdge(ep, es); ms != model.Incomplete {
		return ms
	}

	r := su.rules.Get(ep)

	var evalState model.EdgeState
	if su.areNodesComplete {
		evalState = r.EvaluateQuickState(su.edges)
	} else {
		evalState = r.EvaluateFullState(su.edges)
	}
	if evalState != model.EdgeUnknown && evalState != es {
		return model.Violation
	}

	su.rq.Push(su.edges, r.Affects())

	su.rq.PushNodes(r.InterestedNodes())

	return model.Incomplete
}

func (su stateUpdater) updateEdgeFromRules(
	ep model.EdgePair,
) model.State {
	switch es := su.rules.Get(ep).EvaluateQuickState(su.edges); es {
	case model.EdgeAvoided:
		return su.avoidEdge(ep)
	case model.EdgeExists:
		return su.addEdge(ep)
	case model.EdgeUnknown, model.EdgeOutOfBounds:
		return model.Incomplete
	default:
		return model.Unexpected
	}
}
