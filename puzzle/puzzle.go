package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type nodeWithOptions struct {
	model.Node
	Options []model.TwoArms
}

type Puzzle struct {
	numEdges uint8

	nodes         []model.Node
	twoArmOptions map[model.NodeCoord]nodeWithOptions

	edges *edgesTriState
	rules *ruleSet
	rq    *rulesQueue
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) *Puzzle {
	if numEdges > MaxEdges {
		return nil
	}

	nodes := make([]model.Node, 0, len(nodeLocations))
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		nodes = append(nodes, model.NewNode(nc, nl.IsWhite, nl.Value))
	}

	puzz := &Puzzle{
		numEdges: uint8(numEdges),
		nodes:    nodes,
		edges:    newEdgesStates(numEdges),
		rules:    newRuleSet(numEdges, nodes),
	}

	puzz.rq = newRulesQueue(puzz.edges, puzz, puzz.NumEdges())

	return puzz
}

func (p *Puzzle) DeepCopy() *Puzzle {
	if p == nil {
		return nil
	}

	// I don't think we need to copy nodes because it should only
	// ever be read from, never written to :#

	dc := &Puzzle{
		numEdges:      p.numEdges,
		nodes:         p.nodes,
		twoArmOptions: p.twoArmOptions,
		edges:         p.edges.Copy(),
		rules:         p.rules,
	}

	dc.rq = newRulesQueue(dc.edges, dc, dc.NumEdges())

	return dc
}

func (p *Puzzle) NumEdges() int {
	return int(p.numEdges)
}

func (p *Puzzle) numNodes() int {
	return int(p.numEdges) + 1
}

func (p *Puzzle) getNode(
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

func (p *Puzzle) getPossibleTwoArms(
	node model.Node,
) []model.TwoArms {

	if p.twoArmOptions == nil {
		return p.getTwoArmsForNode(node).Options
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

func (p *Puzzle) populateTwoArmsCache() {
	p.twoArmOptions = make(map[model.NodeCoord]nodeWithOptions, len(p.nodes))
	for _, node := range p.nodes {
		p.twoArmOptions[node.Coord()] = p.getTwoArmsForNode(node)
	}
}

func (p *Puzzle) getTwoArmsForNode(node model.Node) nodeWithOptions {
	options := model.BuildTwoArmOptions(node, p.NumEdges())
	filteredOptions := make([]model.TwoArms, 0, len(options))

	for _, o := range options {
		if !isTwoArmsPossible(node, o, p.edges) {
			continue
		}
		if p.isInTheWayOfOtherNodes(node, o) {
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

func (p *Puzzle) isInTheWayOfOtherNodes(
	node model.Node,
	ta model.TwoArms,
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
		otherNode, ok := p.getNode(a1)
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

	for i, a2 := 1, nc; i < int(ta.Two.Len); i++ {
		a2 = a2.Translate(ta.Two.Heading)
		otherNode, ok := p.getNode(a2)
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

	return false
}

func (p *Puzzle) GetNextTarget(
	cur *model.Target,
) (*model.Target, model.State) {
	t, ok, err := model.GetNextTarget(
		cur,
		p.nodes,
		p.getPossibleTwoArms,
	)

	if err != nil {
		return nil, model.Violation
	}
	if !ok {
		return nil, model.NodesComplete
	}
	return &t, model.Incomplete
}
