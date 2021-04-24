package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type Puzzle struct {
	rules *logic.RuleSet

	nodes         []model.Node
	nearby        []map[model.Cardinal][]*model.Node
	twoArmOptions [][]model.TwoArms

	edges state.TriEdges
	gn    model.GetNoder
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

	nearby := make([]map[model.Cardinal][]*model.Node, len(nodes))
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
	return Puzzle{
		nodes:         p.nodes,
		gn:            p.gn,
		nearby:        p.nearby,
		twoArmOptions: p.twoArmOptions,
		edges:         edges,
		rules:         p.rules,
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

	tas := make([][]model.TwoArms, len(p.nodes))
	for i, n := range p.nodes {
		tas[i] = n.GetFilteredOptions(
			p.twoArmOptions[i],
			&p.edges,
			p.nearby[i],
		)
	}

	t, ok, err := model.GetNextTarget(
		curTarget,
		p.nodes,
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
