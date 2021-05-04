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

	// The following fields are all index-matched (nodes, nearby, twoArmOptions)
	// TODO store it as a slice of a separate struct
	nodes         []model.Node
	nearby        []model.NearbyNodes
	twoArmOptions [][]model.TwoArms

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
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		nodes = append(nodes, model.NewNode(nc, nl.IsWhite, nl.Value))
	}

	nearby := make([]model.NearbyNodes, len(nodes))
	gn := newNodeGrid(nodes)
	edges := state.New(numEdges)
	rules := logic.New(&edges, numEdges, nodes)

	puzz := Puzzle{
		nodes:  nodes,
		gn:     gn,
		nearby: nearby,
		edges:  edges,
		rules:  rules,
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
		nearby:           p.nearby,
		twoArmOptions:    p.twoArmOptions,
		edges:            edges,
		rules:            p.rules,
		areNodesComplete: p.areNodesComplete,
		loop:             newLoop,
	}
}

func updateCache(p *Puzzle) {
	res := make([][]model.TwoArms, 0, len(p.nodes))
	numEdges := p.numEdges()

	for i, n := range p.nodes {
		allTAs := model.BuildTwoArmOptions(n, numEdges)

		maxLensByDir := model.GetMaxArmsByDir(allTAs)
		nearbyNodes := model.BuildNearbyNodes(n, p.gn, maxLensByDir)
		p.nearby[i] = nearbyNodes

		res = append(res,
			n.GetFilteredOptions(allTAs, &p.edges, nearbyNodes),
		)
	}

	p.twoArmOptions = res
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
	for i, n := range p.nodes {
		if cs.IsCoordSeen(n.Coord()) {
			continue
		}
		nodesCopy = append(nodesCopy, n)
		tas = append(tas, n.GetFilteredOptions(
			p.twoArmOptions[i],
			&p.edges,
			p.nearby[i],
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
