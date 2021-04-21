package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

type EdgeChecker interface {
	IsInBounds(model.EdgePair) bool
	IsDefined(model.EdgePair) bool
}

type Queue struct {
	ec EdgeChecker

	toCheck map[model.EdgePair]struct{}

	updated map[model.EdgePair]struct{}
}

func NewQueue(
	ec EdgeChecker,
	numEdges int,
) *Queue {
	return &Queue{
		ec:      ec,
		toCheck: make(map[model.EdgePair]struct{}, 2*numEdges*(numEdges-1)),
		updated: make(map[model.EdgePair]struct{}, 2*numEdges*(numEdges-1)),
	}
}

func (rq *Queue) Push(
	others []model.EdgePair,
) {
	for _, other := range others {
		if !rq.ec.IsInBounds(other) {
			continue
		}
		if rq.ec.IsDefined(other) {
			continue
		}
		rq.toCheck[other] = struct{}{}
	}
}

func (rq *Queue) Pop() (model.EdgePair, bool) {
	for ep := range rq.toCheck {
		// Delete this edge from the map, and return it
		delete(rq.toCheck, ep)
		return ep, true
	}

	return model.EdgePair{}, false
}

func (rq *Queue) NoticeUpdated(
	ep model.EdgePair,
) {
	rq.updated[ep] = struct{}{}
}

func (rq *Queue) Updated() []model.EdgePair {
	u := make([]model.EdgePair, 0, len(rq.updated))
	for uep := range rq.updated {
		u = append(u, uep)
	}

	return u
}

func (rq *Queue) ClearUpdated() {
	for k := range rq.updated {
		delete(rq.updated, k)
	}
}
