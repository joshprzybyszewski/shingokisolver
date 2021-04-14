package puzzle

type rulesQueue struct {
	toCheck []EdgePair

	updated map[EdgePair]struct{}
}

func newRulesQueue() *rulesQueue {
	return &rulesQueue{
		toCheck: make([]EdgePair, 0, 16),
		updated: make(map[EdgePair]struct{}, 16),
	}
}

func (rq *rulesQueue) push(
	others ...EdgePair,
) {
	for _, other := range others {
		if !rq.containsEdgeInToCheck(other) {
			rq.toCheck = append(rq.toCheck, other)
		}
	}
}

func (rq *rulesQueue) containsEdgeInToCheck(
	ep EdgePair,
) bool {
	if _, ok := rq.updated[ep]; ok {
		return true
	}

	for _, c := range rq.toCheck {
		if ep == c {
			return true
		}
	}

	return false
}

func (rq *rulesQueue) pop() (EdgePair, bool) {
	if len(rq.toCheck) == 0 {
		return EdgePair{}, false
	}

	ep := rq.toCheck[0]
	rq.toCheck = rq.toCheck[1:]
	return ep, true
}

func (rq *rulesQueue) noticeUpdated(
	ep EdgePair,
) {
	rq.updated[ep] = struct{}{}
}
