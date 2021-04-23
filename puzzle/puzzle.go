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
		nodes:         p.nodes,
		twoArmOptions: make([]nodeWithOptions, 0, len(p.twoArmOptions)),
		edges:         edges,
		rules:         p.rules,
	}

	for _, tao := range p.twoArmOptions {
		newOptions := getPossibleTwoArmsFromCache(
			&out.edges,
			out,
			tao,
		)
		if len(newOptions) == 0 {
			continue
		}

		out.twoArmOptions = append(out.twoArmOptions, nodeWithOptions{
			Node:    tao.Node,
			Options: newOptions,
		})
	}

	return out
}

func (p Puzzle) getIncompleteNodes(
	ge model.GetEdger,
	minSize int8,
) map[model.NodeCoord]model.Node {
	filtered := make(map[model.NodeCoord]model.Node, len(p.nodes))
	if p.twoArmOptions == nil {
		for _, n := range p.nodes {
			if n.Value() < minSize {
				continue
			}
			switch getNodeState(n, ge) {
			case model.Incomplete:
				filtered[n.Coord()] = n
			}
		}
	}
	for _, tao := range p.twoArmOptions {
		if tao.Value() < minSize {
			continue
		}
		if len(tao.Options) > 1 {
			filtered[tao.Coord()] = tao.Node
		}
	}
	return filtered
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
	if p.GetState(cur.Node.Coord()) == model.Complete {
		return model.Target{}, model.Complete
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
	if p.GetState(model.InvalidNodeCoord) == model.Complete {
		return model.Target{}, model.Complete
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
