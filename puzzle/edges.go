package puzzle

import (
	"errors"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

var (
	ErrEdgeAlreadyExists = errors.New(`already had edge`)
)

func (p *Puzzle) GetOutgoingEdgesFrom(
	coord model.NodeCoord,
) (model.OutgoingEdges, bool) {
	if !p.nodeGrid.IsInBounds(coord) {
		return model.OutgoingEdges{}, false
	}

	return p.nodeGrid.Get(coord), true
}

func (p *Puzzle) IsEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	return p.isEdge(move, nc)
}

func (p *Puzzle) isEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) bool {
	if !p.nodeGrid.IsInBounds(nc) {
		return false
	}
	maxIndex := p.numEdges

	switch move {
	case model.HeadLeft:
		return nc.Col != 0 && p.nodeGrid.Get(nc).IsLeft()
	case model.HeadRight:
		return uint8(nc.Col) != maxIndex && p.nodeGrid.Get(nc).IsRight()
	case model.HeadUp:
		return nc.Row != 0 && p.nodeGrid.Get(nc).IsAbove()
	case model.HeadDown:
		return uint8(nc.Row) != maxIndex && p.nodeGrid.Get(nc).IsBelow()
	default:
		return false
	}
}

func (p *Puzzle) AddEdge(
	move model.Cardinal,
	startNode model.NodeCoord,
) (model.NodeCoord, model.State) {
	if p.isEdge(move, startNode) {
		return model.NodeCoord{}, model.Duplicate
	}

	endNode := startNode.Translate(move)
	if !p.nodeGrid.IsInBounds(endNode) {
		return model.NodeCoord{}, model.Violation
	}

	model.ApplyGridConnections(
		p.nodeGrid,
		move,
		startNode, endNode,
	)

	switch snState := p.getStateForCoord(startNode); snState {
	case model.Complete, model.Incomplete:
		// is fine
	default:
		return model.NodeCoord{}, snState
	}

	switch enState := p.getStateForCoord(endNode); enState {
	case model.Complete, model.Incomplete:
		// is fine
	default:
		return model.NodeCoord{}, enState
	}

	p.paths.add(startNode, endNode)

	return endNode, model.Incomplete
}
