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

	toCheck []model.EdgePair

	updated map[model.EdgePair]struct{}
}

func newRulesQueue(
	bc isInBoundser,
	ed isEdgeDefineder,
	numEdges int,
) *rulesQueue {
	return &rulesQueue{
		bc:      bc,
		ed:      ed,
		toCheck: make([]model.EdgePair, 0, 2*numEdges*(numEdges-1)),
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
		if other.IsIn(rq.toCheck...) {
			continue
		}
		rq.toCheck = append(rq.toCheck, other)

	}
}

func (rq *rulesQueue) pop() (model.EdgePair, bool) {
	if len(rq.toCheck) == 0 {
		return model.EdgePair{}, false
	}

	ep := rq.toCheck[0]

	// TODO to save on re-allocs, we can move the last item to the front.
	rq.toCheck[0] = rq.toCheck[len(rq.toCheck)-1]
	rq.toCheck = rq.toCheck[:len(rq.toCheck)-1]

	// rq.toCheck = rq.toCheck[1:]

	return ep, true
}

func (rq *rulesQueue) noticeUpdated(
	ep model.EdgePair,
) {
	rq.updated[ep] = struct{}{}
}

func (rq *rulesQueue) clearUpdated() {
	// TODO there may be a faster way to do this
	rq.updated = make(map[model.EdgePair]struct{}, len(rq.updated))
}
