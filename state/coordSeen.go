package state

import "github.com/joshprzybyszewski/shingokisolver/model"

type CoordSeener interface {
	Mark(model.NodeCoord)
	IsCoordSeen(model.NodeCoord) bool
}

// TODO this would be a lot less space using bits
type slices struct {
	seenCoords [][]bool
}

func NewCoordSeen(numEdges int) CoordSeener {
	// TODO provide an alternative that requires less space (bit masking?)
	s := make([][]bool, numEdges+1)
	for i := range s {
		s[i] = make([]bool, numEdges+1)
	}
	return slices{
		seenCoords: s,
	}
}

func (s slices) Mark(nc model.NodeCoord) {
	s.seenCoords[nc.Row][nc.Col] = true
}

func (s slices) IsCoordSeen(nc model.NodeCoord) bool {
	return s.seenCoords[nc.Row][nc.Col]
}
