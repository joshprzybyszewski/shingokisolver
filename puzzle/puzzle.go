package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type Puzzle struct {
	numEdges uint8
	nodes    map[model.NodeCoord]model.Node

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

	nodes := map[model.NodeCoord]model.Node{}
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		nodes[nc] = model.NewNode(nc, nl.IsWhite, nl.Value)
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
		numEdges: p.numEdges,
		nodes:    p.nodes,
		edges:    p.edges.Copy(),
		rules:    p.rules,
	}

	dc.rq = newRulesQueue(dc.edges, dc, dc.NumEdges())

	return dc
}

func (p *Puzzle) GetNode(coord model.NodeCoord) (model.Node, bool) {
	n, ok := p.nodes[coord]
	return n, ok
}

func (p *Puzzle) NumEdges() int {
	return int(p.numEdges)
}

func (p *Puzzle) numNodes() int {
	return int(p.numEdges) + 1
}

func (p *Puzzle) GetPossibleTwoArms(
	node model.Node,
) []model.TwoArms {
	options := model.BuildTwoArmOptions(node, p.NumEdges())
	filteredOptions := make([]model.TwoArms, 0, len(options))

	nc := node.Coord()
	for _, o := range options {
		if p.edges.AnyAvoided(nc, o.One) || p.edges.AnyAvoided(nc, o.Two) {
			continue
		}
		filteredOptions = append(filteredOptions, o)
	}

	return filteredOptions
}

func (p *Puzzle) GetNextTarget(
	cur *model.Target,
) (*model.Target, model.State) {
	t, ok, err := model.GetNextTarget(
		cur,
		p.nodes,
		func(n model.Node) int {
			return len(p.GetPossibleTwoArms(n))
		},
	)

	if err != nil {
		return nil, model.Violation
	}
	if !ok {
		return nil, model.NodesComplete
	}
	return &t, model.Incomplete
}
