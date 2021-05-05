package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type Puzzle struct {
	edges state.TriEdges
	gn    model.GetNoder
	rules *logic.RuleSet
	loop  looper

	metas []*model.NodeMeta
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) Puzzle {
	if numEdges > state.MaxEdges {
		return Puzzle{}
	}

	nodes := make([]model.Node, 0, len(nodeLocations))
	for _, nl := range nodeLocations {
		nodes = append(nodes,
			model.NewNode(
				model.NewCoordFromInts(nl.Row, nl.Col),
				nl.IsWhite,
				nl.Value,
			),
		)
	}
	gn := newNodeGrid(nodes)

	metas := make([]*model.NodeMeta, 0, len(nodeLocations))
	for _, n := range nodes {
		tao := model.BuildTwoArmOptions(n, numEdges)
		maxLensByDir := model.GetMaxArmsByDir(tao)
		nearby := model.BuildNearbyNodes(n, gn, maxLensByDir)

		metas = append(metas, &model.NodeMeta{
			Node:          n,
			TwoArmOptions: tao,
			Nearby:        nearby,
		})
	}

	edges := state.New(numEdges)
	rules := logic.New(&edges, numEdges, metas)

	return Puzzle{
		metas: metas,
		gn:    gn,
		edges: edges,
		rules: rules,
	}
}

func (p Puzzle) withNewState(
	edges state.TriEdges,
	newNMs []*model.NodeMeta,
) Puzzle {
	if len(newNMs) != len(p.metas) {
		// TODO remove
		panic(`dev error`)
	}

	var newLoop looper
	if p.loop != nil {
		newLoop = p.loop.withUpdatedEdges(&edges)
	}
	return Puzzle{
		metas: newNMs,
		gn:    p.gn,
		edges: edges,
		rules: p.rules,
		loop:  newLoop,
	}
}

func (p Puzzle) numEdges() int {
	return p.edges.NumEdges()
}

func (p Puzzle) areNodesComplete() bool {
	return p.loop != nil
}

func (p Puzzle) getMetasCopy() []*model.NodeMeta {
	metas := make([]*model.NodeMeta, 0, len(p.metas))
	for _, n := range p.metas {
		metas = append(metas, n.Copy())
	}
	return metas
}

func (p Puzzle) GetNextTarget(
	cur model.Target,
) (model.Target, model.State) {
	return p.getNextTarget(cur)
}

func (p Puzzle) GetFirstTarget() (model.Target, model.State) {
	return p.getNextTarget(model.InvalidTarget)
}

func (p Puzzle) getNextTarget(
	curTarget model.Target,
) (model.Target, model.State) {
	switch s := p.GetState(); s {
	case model.Complete:
		return model.Target{}, model.Complete
	case model.Incomplete:
		// continue on...
	default:
		// Note: If we're NodesComplete, then we'll let our caller handle it.
		return model.Target{}, s
	}

	cs := buildSeenState(p.numEdges(), curTarget)

	metasCpy := make([]*model.NodeMeta, 0, len(p.metas))
	for _, nm := range p.metas {
		if nm.IsComplete || cs.IsCoordSeen(nm.Coord()) {
			continue
		}
		metasCpy = append(metasCpy, nm.Copy())
	}

	t, ok, err := model.GetNextTarget(
		curTarget,
		metasCpy,
	)

	if err != nil {
		return model.Target{}, model.Violation
	}
	if !ok {
		return model.Target{}, model.NodesComplete
	}
	return t, model.Incomplete
}

func buildSeenState(
	numEdges int,
	curTarget model.Target,
) state.CoordSeener {
	cs := state.NewCoordSeen(numEdges)

	for t := curTarget; t.Node != model.InvalidNode; t = *t.Parent {
		cs.Mark(t.Node.Coord())
	}

	return cs
}
