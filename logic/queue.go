package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type EdgeChecker interface {
	IsInBounds(model.EdgePair) bool
	IsDefined(model.EdgePair) bool
}

type Queue struct {
	toCheck state.EdgeQueue

	nodeSeener state.CoordSeener
	nodes      []model.Node
}

func NewQueue(
	numEdges int,
) *Queue {
	return &Queue{
		toCheck:    state.NewEdgeQueue(numEdges),
		nodeSeener: state.NewCoordSeen(numEdges),
		nodes:      make([]model.Node, 0, 16),
	}
}

func (rq *Queue) Push(
	ec EdgeChecker,
	others []model.EdgePair,
) {
	for _, other := range others {
		if !ec.IsDefined(other) {
			rq.toCheck.Push(other)
		}
	}
}

func (rq *Queue) Pop() (model.EdgePair, bool) {
	return rq.toCheck.Pop()
}

func (rq *Queue) PushNodes(
	others []model.Node,
) {
	for _, other := range others {
		if rq.nodeSeener.IsCoordSeen(other.Coord()) {
			continue
		}
		rq.nodeSeener.Mark(other.Coord())
		rq.nodes = append(rq.nodes, other)
	}
}

func (rq *Queue) PopAllNodes() []model.Node {
	all := rq.nodes

	rq.nodeSeener.UnmarkAll()
	rq.nodes = rq.nodes[:0]

	return all

}
