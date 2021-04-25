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
}

func NewQueue(
	numEdges int,
) *Queue {
	return &Queue{
		toCheck: state.NewEdgeQueue(numEdges),
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
