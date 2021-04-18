package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type isInBoundser interface {
	isInBounds(model.EdgePair) bool
}

type isEdgeDefineder interface {
	isEdgeDefined(model.EdgePair) bool
}

type rulesQueue struct {
	bc isInBoundser
	ed isEdgeDefineder

	// toCheck []model.EdgePair
	toCheck map[model.EdgePair]struct{}

	updated map[model.EdgePair]struct{}
}

func newRulesQueue(
	bc isInBoundser,
	ed isEdgeDefineder,
	numEdges int,
) *rulesQueue {
	return &rulesQueue{
		bc: bc,
		ed: ed,
		// toCheck: make([]model.EdgePair, 0, 2*numEdges*(numEdges-1)),
		toCheck: make(map[model.EdgePair]struct{}, 2*numEdges*(numEdges-1)),
		updated: make(map[model.EdgePair]struct{}, 2*numEdges*(numEdges-1)),
	}
}

func (rq *rulesQueue) push(
	others ...model.EdgePair,
) {
	for _, other := range others {
		if !rq.bc.isInBounds(other) {
			continue
		}
		if rq.ed.isEdgeDefined(other) {
			continue
		}
		rq.toCheck[other] = struct{}{}
	}
}

func (rq *rulesQueue) pop() (model.EdgePair, bool) {
	for ep := range rq.toCheck {
		// Delete this edge from the map, and return it
		delete(rq.toCheck, ep)
		return ep, true
	}

	return model.EdgePair{}, false
}

func (rq *rulesQueue) noticeUpdated(
	ep model.EdgePair,
) {
	rq.updated[ep] = struct{}{}
}

func (rq *rulesQueue) clearUpdated() {
	// TODO there may be a faster way to clear this map.
	rq.updated = make(map[model.EdgePair]struct{}, len(rq.updated))
}
