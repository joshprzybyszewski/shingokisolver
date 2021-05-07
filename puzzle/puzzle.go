package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type Puzzle struct {
	edges state.TriEdges
	rules *logic.RuleSet
	loop  looper

	// TODO keep a slice of structs, not pointers
	metas []*model.NodeMeta
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) Puzzle {
	if numEdges > state.MaxEdges {
		return Puzzle{}
	}

	metas := make([]*model.NodeMeta, 0, len(nodeLocations))
	for _, nl := range nodeLocations {

		n := model.NewNode(
			model.NewCoordFromInts(nl.Row, nl.Col),
			nl.IsWhite,
			nl.Value,
		)

		metas = append(metas, &model.NodeMeta{
			Node: n,
		})
	}

	// we need this GetNoder to build nearby nodes
	gn := newNodeGrid(metas)

	for _, m := range metas {
		m.TwoArmOptions = model.BuildTwoArmOptions(m.Node, numEdges)

		maxLensByDir := model.GetMaxArmsByDir(m.TwoArmOptions)

		m.Nearby = model.BuildNearbyNodes(m.Node, gn, maxLensByDir)
	}

	edges := state.New(numEdges)
	rules := logic.New(&edges, numEdges, metas)

	return Puzzle{
		metas: metas,
		edges: edges,
		rules: rules,
	}
}

func (p Puzzle) buildGetNoder() model.GetNoder {
	if p.metas == nil {
		return nil
	}
	return newNodeGrid(p.metas)
}

func (p Puzzle) withNewState(
	edges state.TriEdges,
	newNMs []*model.NodeMeta,
) Puzzle {
	// if len(newNMs) != len(p.metas) {
	// 	panic(`dev error`)
	// }

	var newLoop looper
	if p.loop != nil {
		newLoop = p.loop.withUpdatedEdges(&edges)
	}
	return Puzzle{
		metas: newNMs,
		edges: edges,
		rules: p.rules,
		loop:  newLoop,
	}
}

func (p Puzzle) numEdges() int {
	return p.edges.NumEdges()
}

func (p Puzzle) areNodesComplete() bool {
	// TODO convert to looking at all node metas for complete.
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
	_, loopState := p.getStateOfLoop(cur.Node.Coord())
	switch loopState {
	case model.Incomplete:
		// continue down
	default:
		return model.InvalidTarget, loopState
	}
	return p.getNextTarget(cur)
}

func (p Puzzle) GetFirstTarget() (model.Target, model.State) {
	return p.getNextTarget(model.InvalidTarget)
}

func (p Puzzle) getNextTarget(
	curTarget model.Target,
) (model.Target, model.State) {

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
