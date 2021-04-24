package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type nodeWithOptions struct {
	Options []model.TwoArms
	model.Node
}

type Puzzle struct {
	rules *logic.RuleSet

	twoArmOptions []nodeWithOptions
	nodes         []model.Node
	edges         state.TriEdges
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

	edges := state.New(numEdges)
	return Puzzle{
		nodes: nodes,
		edges: edges,
		rules: logic.New(&edges, numEdges, nodes),
	}
}

func (p Puzzle) withNewState(
	edges state.TriEdges,
) Puzzle {
	return Puzzle{
		nodes:         p.nodes,
		twoArmOptions: p.twoArmOptions,
		edges:         edges,
		rules:         p.rules,
	}
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

	// TODO determine if there's a better data structure/it's worth it
	var gn model.GetNoder = (nodeList)(p.nodes)

	tas := make([][]model.TwoArms, len(p.nodes))
	for i, n := range p.nodes {
		tas[i] = getPossibleTwoArmsWithNewEdges(n, &p.edges, gn, p.twoArmOptions)
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

func getPossibleTwoArmsWithNewEdges(
	node model.Node,
	ge model.GetEdger,
	gn model.GetNoder,
	allTAOs []nodeWithOptions,
) []model.TwoArms {

	var tao nodeWithOptions
	found := false
	for _, o := range allTAOs {
		if o.Node.Coord() == node.Coord() {
			tao = o
			found = true
			break
		}
	}
	if !found {
		panic(`not found!`)
	}

	maxLensByDir := model.GetMaxArmsByDir(tao.Options)
	nearbyNodes := model.BuildNearbyNodes(tao.Node, gn, maxLensByDir)
	return tao.Node.GetFilteredOptions(tao.Options, ge, nearbyNodes)
}
