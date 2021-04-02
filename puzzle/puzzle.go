package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type Puzzle struct {
	numEdges uint8
	nodes    map[model.NodeCoord]model.Node

	paths *paths

	nodeGrid model.Grid
}

func NewPuzzle(
	numEdges int,
	nodeLocations []model.NodeLocation,
) *Puzzle {
	if numEdges > model.MAX_EDGES {
		return nil
	}

	nodes := map[model.NodeCoord]model.Node{}
	for _, nl := range nodeLocations {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		nodes[nc] = model.NewNode(nl.IsWhite, nl.Value)
	}

	return &Puzzle{
		numEdges: uint8(numEdges),
		nodes:    nodes,
		nodeGrid: model.NewGrid(numEdges),
		paths:    &paths{},
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
		nodeGrid: p.nodeGrid.Copy(),
		paths:    p.paths.copy(),
	}
}

func (p *Puzzle) GetNode(coord model.NodeCoord) (model.Node, bool) {
	n, ok := p.nodes[coord]
	return n, ok
}

func (p *Puzzle) GetLooseEnd() (model.NodeCoord, bool) {
	return p.paths.getLooseEnd()
}

func (p *Puzzle) NumLooseEnds() int {
	return p.paths.numLooseEnds()
}

func (p *Puzzle) NumEdges() int {
	return int(p.numEdges)
}

func (p *Puzzle) numNodes() int {
	return int(p.numEdges) + 1
}
