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

	nodes []nodeMeta

	// TODO replace this as a method that checks is loop != nil
	areNodesComplete bool
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) Puzzle {
	if numEdges > state.MaxEdges {
		return Puzzle{}
	}

	nodes := make([]model.Node, 0, len(nodeLocations))
	nodeMetas := make([]nodeMeta, 0, len(nodeLocations))
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		n := model.NewNode(nc, nl.IsWhite, nl.Value)
		nodes = append(nodes, n)
		nodeMetas = append(nodeMetas, nodeMeta{
			node: n,
		})
	}

	gn := newNodeGrid(nodes)
	edges := state.New(numEdges)
	rules := logic.New(&edges, numEdges, nodes)

	puzz := Puzzle{
		nodes: nodeMetas,
		gn:    gn,
		edges: edges,
		rules: rules,
	}

	updateCache(&puzz)

	return puzz
}

func (p Puzzle) withNewState(
	edges state.TriEdges,
) Puzzle {
	var newLoop looper
	if p.loop != nil {
		newLoop = p.loop.withUpdatedEdges(&edges)
	}
	return Puzzle{
		nodes:            p.nodes,
		gn:               p.gn,
		edges:            edges,
		rules:            p.rules,
		areNodesComplete: p.areNodesComplete,
		loop:             newLoop,
	}
}

func updateCache(p *Puzzle) {
	numEdges := p.numEdges()
	oldMetas := p.nodes
	newMetas := make([]nodeMeta, len(oldMetas))

	for i, nm := range oldMetas {
		allTAs := model.BuildTwoArmOptions(nm.node, numEdges)

		maxLensByDir := model.GetMaxArmsByDir(allTAs)
		nearbyNodes := model.BuildNearbyNodes(nm.node, p.gn, maxLensByDir)

		newMetas[i].node = nm.node
		newMetas[i].nearby = nearbyNodes
		newMetas[i].twoArmOptions = nm.node.GetFilteredOptions(
			allTAs,
			&p.edges,
			nearbyNodes,
		)
	}
	p.nodes = newMetas
}

func (p Puzzle) numEdges() int {
	return p.edges.NumEdges()
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

	nodesCopy := make([]model.Node, 0, len(p.nodes))
	tas := make([][]model.TwoArms, 0, len(p.nodes))
	for _, nm := range p.nodes {
		if cs.IsCoordSeen(nm.node.Coord()) {
			continue
		}

		nodesCopy = append(nodesCopy, nm.node)
		tas = append(tas, nm.node.GetFilteredOptions(
			nm.twoArmOptions,
			&p.edges,
			nm.nearby,
		))
	}

	t, ok, err := model.GetNextTarget(
		curTarget,
		nodesCopy,
		tas,
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
