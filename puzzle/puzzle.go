package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type Puzzle struct {
	numEdges uint8
	nodes    map[model.NodeCoord]model.Node

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
	}
}

func (p *Puzzle) GetNode(coord model.NodeCoord) (model.Node, bool) {
	n, ok := p.nodes[coord]
	return n, ok
}

func (p *Puzzle) GetLooseEnd() (model.NodeCoord, model.State) {
	for r := model.RowIndex(0); r <= model.RowIndex(p.numEdges); r++ {
		for c := model.ColIndex(0); c <= model.ColIndex(p.numEdges); c++ {
			nc := model.NewCoord(r, c)
			oe, ok := p.GetOutgoingEdgesFrom(nc)
			if !ok {
				return model.NodeCoord{}, model.Violation
			}

			switch numEdges := oe.GetNumOutgoingDirections(); {
			case numEdges > 2:
				return model.NodeCoord{}, model.Violation
			case numEdges == 1:
				return nc, model.Incomplete
			}
		}
	}
	return model.NodeCoord{}, model.NodesComplete
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
