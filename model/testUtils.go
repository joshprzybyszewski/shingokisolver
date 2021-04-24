package model

var _ GetEdger = testGetEdger{}

type testGetEdger struct {
	existing []EdgePair
	avoided  []EdgePair

	numEdges int
}

func (ge testGetEdger) GetEdge(ep EdgePair) EdgeState {
	for _, ex := range ge.existing {
		if ex == ep {
			return EdgeExists
		}
	}
	for _, av := range ge.avoided {
		if av == ep {
			return EdgeAvoided
		}
	}
	return EdgeUnknown
}

func (ge testGetEdger) IsEdge(ep EdgePair) bool {
	return ge.GetEdge(ep) == EdgeExists
}

func (ge testGetEdger) IsAvoided(ep EdgePair) bool {
	return ge.GetEdge(ep) == EdgeAvoided
}

func (ge testGetEdger) IsInBounds(ep EdgePair) bool {
	// This is copied from triEdges
	if ep.Row < 0 || ep.Col < 0 {
		// negative coords are bad
		return false
	}
	switch ep.Cardinal {
	case HeadRight:
		return int(ep.Row) <= ge.numEdges && int(ep.Col) < ge.numEdges
	case HeadDown:
		return int(ep.Row) < ge.numEdges && int(ep.Col) <= ge.numEdges
	default:
		// unexpected input
		return false
	}
}

func (ge testGetEdger) AllExist(nc NodeCoord, a Arm) bool {
	allExisting := true
	ep := NewEdgePair(nc, a.Heading)
	for i := 1; i <= int(a.Len); i++ {
		if !ge.IsEdge(ep) {
			allExisting = false
			break
		}
		ep = ep.Next(a.Heading)
	}
	return allExisting
}

func (ge testGetEdger) Any(nc NodeCoord, a Arm) (bool, bool) {
	anyExisting, anyAvoided := false, false
	ep := NewEdgePair(nc, a.Heading)
	for i := 1; i <= int(a.Len); i++ {
		if ge.IsEdge(ep) {
			anyExisting = true
		}
		if ge.IsAvoided(ep) {
			anyAvoided = true
		}
		ep = ep.Next(a.Heading)
	}
	return anyExisting, anyAvoided
}

func (ge testGetEdger) AnyAvoided(nc NodeCoord, a Arm) bool {
	_, any := ge.Any(nc, a)
	return any
}

func (ge testGetEdger) NumEdges() int {
	return ge.numEdges
}
