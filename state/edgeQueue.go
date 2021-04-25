package state

import "github.com/joshprzybyszewski/shingokisolver/model"

type EdgeQueue interface {
	Push(model.EdgePair)
	IsEdgeSeen(model.EdgePair) bool

	// Pop will get and remove an edge pair
	Pop() (model.EdgePair, bool)
}

type edgeState struct {
	rows []bitData
	cols []bitData
}

func NewEdgeQueue(numEdges int) EdgeQueue {
	return edgeState{
		rows: make([]bitData, numEdges+1),
		cols: make([]bitData, numEdges+1),
	}
}

func (es edgeState) IsEdgeSeen(ep model.EdgePair) bool {
	switch ep.Cardinal {
	case model.HeadRight:
		return (es.rows[ep.Row] & masks[ep.Col]) != 0
	case model.HeadDown:
		return (es.cols[ep.Col] & masks[ep.Row]) != 0
	}
	return false
}

func (es edgeState) Push(ep model.EdgePair) {
	switch ep.Cardinal {
	case model.HeadRight:
		es.rows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		es.cols[ep.Col] |= masks[ep.Row]
	}
}

func (es edgeState) Pop() (model.EdgePair, bool) {
	for r, row := range es.rows {
		if row == 0 {
			continue
		}
		for c := 0; c < len(es.cols); c++ {
			if row&masks[c] == 0 {
				continue
			}
			es.rows[r] = row ^ masks[c]

			return model.NewEdgePair(
				model.NewCoordFromInts(r, c),
				model.HeadRight,
			), true
		}
	}

	for c, col := range es.cols {
		if col == 0 {
			continue
		}
		for r := 0; r < len(es.cols); r++ {
			if col&masks[r] == 0 {
				continue
			}
			es.cols[c] = col ^ masks[r]

			return model.NewEdgePair(
				model.NewCoordFromInts(r, c),
				model.HeadDown,
			), true
		}
	}

	return model.InvalidEdgePair, false
}
