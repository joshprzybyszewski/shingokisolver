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
	edges *state.TriEdges
	rules *logic.RuleSet

	twoArmOptions map[model.NodeCoord]nodeWithOptions
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

	return Puzzle{
		nodes: nodes,
		edges: state.New(numEdges),
		rules: logic.New(numEdges, nodes),
	}
}

func (p Puzzle) DeepCopy() Puzzle {
	// I don't think we need to copy nodes because it should only
	// ever be read from, never written to :#

	return Puzzle{
		nodes:         p.nodes,
		twoArmOptions: p.twoArmOptions,
		edges:         p.edges.Copy(),
		rules:         p.rules,
	}
}

func (p Puzzle) withNewState(
	edges *state.TriEdges,
) Puzzle {
	return Puzzle{
		nodes:         p.nodes,
		twoArmOptions: p.twoArmOptions,
		edges:         edges,
		rules:         p.rules,
	}
}

func (p Puzzle) NumEdges() int {
	return p.edges.NumEdges()
}

func (p Puzzle) numNodes() int {
	return p.NumEdges() + 1
}

func (p Puzzle) getNode(
	nc model.NodeCoord,
) (model.Node, bool) {
	if p.twoArmOptions == nil {
		for _, n := range p.nodes {
			if n.Coord() == nc {
				return n, true
			}
		}
		return model.Node{}, false
	}

	tao, ok := p.twoArmOptions[nc]
	if !ok {
		return model.Node{}, false
	}

	return tao.Node, true
}

func (p Puzzle) getPossibleTwoArms(
	node model.Node,
) []model.TwoArms {

	if p.twoArmOptions == nil {
		return getTwoArmsForNode(node, p.NumEdges(), p.edges, p).Options
	}

	tao := p.twoArmOptions[node.Coord()]
	filteredOptions := make([]model.TwoArms, 0, len(tao.Options))

	for _, o := range tao.Options {
		if isTwoArmsPossible(node, o, p.edges) {
			filteredOptions = append(filteredOptions, o)
		}
	}

	return filteredOptions
}

func getTwoArmsCache(
	nodes []model.Node,
	numEdges int,
	ge model.GetEdger,
	gn getNoder,
) map[model.NodeCoord]nodeWithOptions {
	nwoByNC := make(map[model.NodeCoord]nodeWithOptions, len(nodes))

	for _, node := range nodes {
		nwoByNC[node.Coord()] = getTwoArmsForNode(node, numEdges, ge, gn)
	}

	return nwoByNC
}

func getTwoArmsForNode(
	node model.Node,
	numEdges int,
	ge model.GetEdger,
	gn getNoder,
) nodeWithOptions {
	options := model.BuildTwoArmOptions(node, numEdges)
	filteredOptions := make([]model.TwoArms, 0, len(options))

	for _, o := range options {
		if !isTwoArmsPossible(node, o, ge) {
			continue
		}
		if isInTheWayOfOtherNodes(node, o, gn) {
			continue
		}
		filteredOptions = append(filteredOptions, o)
	}

	return nodeWithOptions{
		Node:    node,
		Options: filteredOptions,
	}
}

func isTwoArmsPossible(
	node model.Node,
	ta model.TwoArms,
	ge model.GetEdger,
) bool {

	nc := node.Coord()
	return !ge.AnyAvoided(nc, ta.One) &&
		!ge.AnyAvoided(nc, ta.Two) &&
		!ge.IsEdge(
			model.NewEdgePair(
				nc.TranslateAlongArm(ta.One),
				ta.One.Heading,
			),
		) &&
		!ge.IsEdge(
			model.NewEdgePair(
				nc.TranslateAlongArm(ta.Two),
				ta.Two.Heading,
			),
		)
}

type getNoder interface {
	getNode(model.NodeCoord) (model.Node, bool)
}

func isInTheWayOfOtherNodes(
	node model.Node,
	ta model.TwoArms,
	gn getNoder,
) bool {

	nc := node.Coord()

	a1StraightLineVal := ta.One.Len
	a2StraightLineVal := ta.Two.Len
	if node.Type() == model.WhiteNode {
		a1StraightLineVal = ta.One.Len + ta.Two.Len
		a2StraightLineVal = ta.One.Len + ta.Two.Len
	}

	for i, a1 := 1, nc; i < int(ta.One.Len); i++ {
		a1 = a1.Translate(ta.One.Heading)
		otherNode, ok := gn.getNode(a1)
		if !ok {
			continue
		}
		if otherNode.Type() == model.BlackNode {
			// this arm would pass through this node in a straight line
			// that makes this arm impossible.
			return true
		}
		if otherNode.Value() != a1StraightLineVal {
			// this arm would pass through the other node
			// in a straight line, and the value would not be tenable
			return true
		}
	}
	if otherNode, ok := gn.getNode(nc.TranslateAlongArm(ta.One)); ok {
		if otherNode.Type() == model.WhiteNode {
			// this arm would end in a white node. That's not ok because
			// we would need to continue through it
			return true
		}
		if otherNode.Value()-a1StraightLineVal < 1 {
			// this arm meets the other node, and would require going
			// next in a perpendicular path. Since this arm would
			// contribute too much to its value, we can filter it ou.
			return true
		}
	}

	for i, a2 := 1, nc; i < int(ta.Two.Len); i++ {
		a2 = a2.Translate(ta.Two.Heading)
		otherNode, ok := gn.getNode(a2)
		if !ok {
			continue
		}
		if otherNode.Type() == model.BlackNode {
			// this arm would pass through this node in a straight line
			// that makes this arm impossible.
			return true
		}
		if otherNode.Value() != a2StraightLineVal {
			// this arm would pass through the other node
			// in a straight line, and the value would not be tenable
			return true
		}
	}

	if otherNode, ok := gn.getNode(nc.TranslateAlongArm(ta.Two)); ok {
		if otherNode.Type() == model.WhiteNode {
			// this arm would end in a white node. That's not ok because
			// we would need to continue through it
			return true
		}
		if otherNode.Value()-a2StraightLineVal < 1 {
			// this arm meets the other node, and would require going
			// next in a perpendicular path. Since this arm would
			// contribute too much to its value, we can filter it ou.
			return true
		}
	}

	return false
}

func (p Puzzle) GetNextTarget(
	cur model.Target,
) (model.Target, model.State) {
	if p.GetState(cur.Node.Coord()) == model.Complete {
		return model.Target{}, model.Complete
	}

	t, ok, err := model.GetNextTarget(
		cur,
		p.nodes,
		p.getPossibleTwoArms,
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

	t, ok, err := model.GetNextTarget(
		model.InvalidTarget,
		p.nodes,
		p.getPossibleTwoArms,
	)

	if err != nil {
		return model.Target{}, model.Violation
	}
	if !ok {
		return model.Target{}, model.NodesComplete
	}
	return t, model.Incomplete
}
