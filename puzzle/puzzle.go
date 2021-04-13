package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type Puzzle struct {
	numEdges uint8
	nodes    map[model.NodeCoord]model.Node

	edges *edgesTriState
	rules *ruleSet
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

	return &Puzzle{
		numEdges: uint8(numEdges),
		nodes:    nodes,
		edges:    newEdgesBits(uint8(numEdges)),
		rules:    newRuleSet(numEdges, nodes),
	}
}

func (p *Puzzle) DeepCopy() *Puzzle {
	if p == nil {
		return nil
	}

	// I don't think we need to copy nodes because it should only
	// ever be read from, never written to :#
	return &Puzzle{
		numEdges: p.numEdges,
		nodes:    p.nodes,
		edges:    p.edges.Copy(),
		rules:    p.rules,
	}
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

func (p *Puzzle) Targets() []model.Target {
	return model.BuildTargets(p.nodes, p.NumEdges())
}
