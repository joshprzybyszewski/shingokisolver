package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

type isDefineder interface {
	IsDefined(model.EdgePair) bool
}

type looper interface {
	NumNodesInLoop() int
	IsLoop() bool
	GetUnknownEdge(isDefineder) (model.EdgePair, model.State)

	withUpdatedEdges(model.GetEdger) looper
}

var _ looper = (*allSegments)(nil)

type segmentCap struct {
	coord model.NodeCoord
	edge  model.EdgePair
}

type pathSegment struct {
	start segmentCap
	end   segmentCap

	seenNodes []model.Node
}

type allSegments struct {
	all []pathSegment

	looseEnds state.CoordSeener

	isComplete   bool
	numNodesSeen int
}

func newAllSegmentsFromNodesComplete(
	nodes []model.Node,
	ge model.GetEdger,
) *allSegments {
	// TODO
	return nil
}

func (as *allSegments) NumNodesInLoop() int {
	if as == nil {
		return 0
	}
	return as.numNodesSeen
}

func (as *allSegments) IsLoop() bool {
	if as == nil {
		return false
	}
	return as.isComplete
}

func (as *allSegments) GetUnknownEdge(
	d isDefineder,
) (model.EdgePair, model.State) {
	if as == nil || as.IsLoop() || len(as.all) == 0 {
		// cannot find an unknown edge on a loop'ed graph!
		return model.InvalidEdgePair, model.Violation
	}

	looseEnd := as.all[0].start.coord

	for _, dir := range model.AllCardinals {
		ep := model.NewEdgePair(looseEnd, dir)
		if !d.IsDefined(ep) {
			return ep, model.Incomplete
		}
	}

	return model.InvalidEdgePair, model.Violation
}

func (as *allSegments) withUpdatedEdges(
	ge model.GetEdger,
) looper {
	if as == nil {
		// return a nil *allSegments
		return as
	}
	// TODO
	return nil
}
