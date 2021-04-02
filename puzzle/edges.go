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
	isEdge, _ := p.isEdge(move, nc)
	return isEdge
}

func (p *Puzzle) isEdge(
	move model.Cardinal,
	nc model.NodeCoord,
) (bool, bool) {
	if !p.nodeGrid.IsInBounds(nc) {
		return false, false
	}
	maxIndex := p.numEdges

	switch move {
	case model.HeadLeft:
		return nc.Col != 0 && p.nodeGrid.Get(nc).IsLeft(), true
	case model.HeadRight:
		return uint8(nc.Col) != maxIndex && p.nodeGrid.Get(nc).IsRight(), true
	case model.HeadUp:
		return nc.Row != 0 && p.nodeGrid.Get(nc).IsAbove(), true
	case model.HeadDown:
		return uint8(nc.Row) != maxIndex && p.nodeGrid.Get(nc).IsBelow(), true
	default:
		return false, false
	}
}

func (p *Puzzle) AddEdge(
	move model.Cardinal,
	startNode model.NodeCoord,
) (model.NodeCoord, *Puzzle, error) {
	endNode := startNode.Translate(move)
	if !p.nodeGrid.IsInBounds(endNode) {
		return model.NodeCoord{}, nil, errors.New(`next point is out of bounds`)
	}

	isEdge, isValid := p.isEdge(move, startNode)
	if !isValid {
		return model.NodeCoord{}, nil, errors.New(`invalid input`)
	}
	if isEdge {
		return model.NodeCoord{}, nil, ErrEdgeAlreadyExists
	}

	model.ApplyGridConnections(
		p.nodeGrid,
		move,
		startNode, endNode,
	)

	return endNode, p, nil
}
