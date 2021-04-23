package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

type EdgeChecker interface {
	IsInBounds(model.EdgePair) bool
	IsDefined(model.EdgePair) bool
}

type Queue struct {
	toCheck map[model.EdgePair]struct{}
}

func NewQueue(
	numEdges int,
) *Queue {
	return &Queue{
		toCheck: make(map[model.EdgePair]struct{}, 2*numEdges*(numEdges-1)),
	}
}

func (rq *Queue) Push(
	ec EdgeChecker,
	others []model.EdgePair,
) {
	for _, other := range others {
		if !ec.IsDefined(other) {
			rq.toCheck[other] = struct{}{}
		}
	}
}

func (rq *Queue) Pop() (model.EdgePair, bool) {
	for ep := range rq.toCheck {
		// Delete this edge from the map, and return it
		delete(rq.toCheck, ep)
		return ep, true
	}

	return model.InvalidEdgePair, false
}
