package state

import "github.com/joshprzybyszewski/shingokisolver/model"

type CoordSeener interface {
	Mark(model.NodeCoord)
	UnmarkAll()
	IsCoordSeen(model.NodeCoord) bool
}

type coordSeen struct {
	rows []bitData
}

func NewCoordSeen(numEdges int) CoordSeener {
	return coordSeen{
		rows: make([]bitData, numEdges+1),
	}
}

func (s coordSeen) Mark(nc model.NodeCoord) {
	s.rows[nc.Row] = s.rows[nc.Row] | masks[nc.Col]
}

func (s coordSeen) UnmarkAll() {
	for i := range s.rows {
		s.rows[i] = 0
	}
}

func (s coordSeen) IsCoordSeen(nc model.NodeCoord) bool {
	return (s.rows[nc.Row] & masks[nc.Col]) != 0
}
