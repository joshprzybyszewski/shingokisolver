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
	edges state.TriEdges
	rules *logic.RuleSet

	twoArmOptions []nodeWithOptions
	nodes         []model.Node
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
	out := Puzzle{
		nodes: p.nodes,
		// TODO don't copy this entirely anymore
		// twoArmOptions: make([]nodeWithOptions, 0, len(p.twoArmOptions)),
		twoArmOptions: p.twoArmOptions,
		edges:         edges,
		rules:         p.rules,
	}

	// for _, tao := range p.twoArmOptions {
	// 	newOptions := getPossibleTwoArmsFromCache(
	// 		&out.edges,
	// 		out,
	// 		tao,
	// 	)
	// 	if len(newOptions) == 0 {
	// 		continue
	// 	}

	// 	out.twoArmOptions = append(out.twoArmOptions, nodeWithOptions{
	// 		Node:    tao.Node,
	// 		Options: newOptions,
	// 	})
	// }

	return out
}

func (p Puzzle) NumEdges() int {
	return p.edges.NumEdges()
}

func (p Puzzle) GetNode(
	nc model.NodeCoord,
) (model.Node, bool) {
	for _, n := range p.nodes {
		if n.Coord() == nc {
			return n, true
		}
	}
	return model.Node{}, false
}

func getPossibleTwoArmsWithNewEdges(
	node model.Node,
	ge model.GetEdger,
	gn model.GetNoder,
	allTAOs []nodeWithOptions,
) []model.TwoArms {

	var tao nodeWithOptions
	for _, o := range allTAOs {
		if o.Node.Coord() == node.Coord() {
			tao = o
			break
		}
	}

	return getPossibleTwoArmsFromCache(ge, gn, tao)
}

func getPossibleTwoArmsFromCache(
	ge model.GetEdger,
	gn model.GetNoder,
	tao nodeWithOptions,
) []model.TwoArms {
	nearbyNodes := model.BuildNearbyNodes(tao.Node, tao.Options, gn)
	return tao.Node.GetFilteredOptions(tao.Options, ge, nearbyNodes)
}

func (p Puzzle) GetNextTarget(
	cur model.Target,
) (model.Target, model.State) {
	switch s := p.GetState(); s {
	case model.Complete:
		return model.Target{}, model.Complete
	case model.Incomplete, model.NodesComplete:
		// continue. we'll let out caller handle the 'nodes complete' state
		// TODO handle the nodes complete state now!
	default:
		return model.Target{}, s
	}

	tas := make([][]model.TwoArms, len(p.nodes))
	for i, n := range p.nodes {
		tas[i] = getPossibleTwoArmsWithNewEdges(n, &p.edges, p, p.twoArmOptions)
	}

	t, ok, err := model.GetNextTarget(
		cur,
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

func (p Puzzle) GetFirstTarget() (model.Target, model.State) {
	switch s := p.GetState(); s {
	case model.Complete:
		return model.Target{}, model.Complete
	case model.Incomplete:
		// continue. we'll let out caller handle the 'nodes complete' state
	default:
		return model.Target{}, s
	}

	tas := make([][]model.TwoArms, len(p.nodes))
	for i, n := range p.nodes {
		tas[i] = getPossibleTwoArmsWithNewEdges(n, &p.edges, p, p.twoArmOptions)
	}

	t, ok, err := model.GetNextTarget(
		model.InvalidTarget,
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
