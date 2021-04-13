package puzzle

type rulesQueue struct {
	toCheck []edgePair

	updated map[edgePair]struct{}
}

func newRulesQueue() *rulesQueue {
	return &rulesQueue{
		toCheck: make([]edgePair, 0, 16),
		updated: make(map[edgePair]struct{}, 16),
	}
}

func (rq *rulesQueue) push(
	others ...edgePair,
) {
	for _, other := range others {
		if !rq.containsEdgeInToCheck(other) {
			rq.toCheck = append(rq.toCheck, other)
		}
	}
}

func (rq *rulesQueue) containsEdgeInToCheck(
	ep edgePair,
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

func (rq *rulesQueue) pop() (edgePair, bool) {
	if len(rq.toCheck) == 0 {
		return edgePair{}, false
	}

	ep := rq.toCheck[0]
	rq.toCheck = rq.toCheck[1:]
	return ep, true
}

func (rq *rulesQueue) noticeUpdated(
	ep edgePair,
) {
	rq.updated[ep] = struct{}{}
}
