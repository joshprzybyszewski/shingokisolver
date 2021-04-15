package puzzle

type isInBoundser interface {
	isInBounds(EdgePair) bool
}

type isEdgeDefineder interface {
	isEdgeDefined(EdgePair) bool
}

type rulesQueue struct {
	bc isInBoundser
	ed isEdgeDefineder

	toCheck []EdgePair

	updated map[EdgePair]struct{}
}

func newRulesQueue(
	bc isInBoundser,
	ed isEdgeDefineder,
	numEdges int,
) *rulesQueue {
	return &rulesQueue{
		bc:      bc,
		ed:      ed,
		toCheck: make([]EdgePair, 0, 2*numEdges*(numEdges-1)),
		updated: make(map[EdgePair]struct{}, 2*numEdges*(numEdges-1)),
	}
}

func (rq *rulesQueue) push(
	others ...EdgePair,
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

func (rq *rulesQueue) pop() (EdgePair, bool) {
	if len(rq.toCheck) == 0 {
		return EdgePair{}, false
	}

	ep := rq.toCheck[0]

	// TODO to save on re-allocs, we can move the last item to the front.
	// rq.toCheck[0] = rq.toCheck[len(rq.toCheck)-1]
	// rq.toCheck = rq.toCheck[:len(rq.toCheck)-1]

	rq.toCheck = rq.toCheck[1:]

	return ep, true
}

func (rq *rulesQueue) noticeUpdated(
	ep EdgePair,
) {
	rq.updated[ep] = struct{}{}
}

func (rq *rulesQueue) clearUpdated() {
	// TODO there may be a faster way to do this
	rq.updated = make(map[EdgePair]struct{}, len(rq.updated))
}
