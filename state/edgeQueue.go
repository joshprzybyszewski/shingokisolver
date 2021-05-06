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

	edges []model.EdgePair
}

func NewEdgeQueue(numEdges int) EdgeQueue {
	return &edgeState{
		rows:  make([]bitData, numEdges+1),
		cols:  make([]bitData, numEdges+1),
		edges: make([]model.EdgePair, 16),
	}
}

func (es *edgeState) IsEdgeSeen(ep model.EdgePair) bool {
	switch ep.Cardinal {
	case model.HeadRight:
		return (es.rows[ep.Row] & masks[ep.Col]) != 0
	case model.HeadDown:
		return (es.cols[ep.Col] & masks[ep.Row]) != 0
	}
	return false
}

func (es *edgeState) Push(ep model.EdgePair) {
	es.edges = append(es.edges, ep)
	switch ep.Cardinal {
	case model.HeadRight:
		es.rows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		es.cols[ep.Col] |= masks[ep.Row]
	}
}

func (es *edgeState) Pop() (model.EdgePair, bool) {
	if len(es.edges) == 0 {
		return model.InvalidEdgePair, false
	}
	ep := es.edges[0]
	es.rows[ep.Row] = es.rows[ep.Row] ^ masks[ep.Col]

	es.edges[0] = es.edges[len(es.edges)-1]
	es.edges = es.edges[:len(es.edges)-1]

	return ep, true
}
