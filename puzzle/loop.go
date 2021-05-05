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
	coord    model.NodeCoord
	outbound model.Cardinal
}

type pathSegment struct {
	start segmentCap
	end   segmentCap

	numNodes int
}

type allSegments struct {
	all []pathSegment

	numNodesSeen int
}

func newAllSegmentsFromNodesComplete(
	metas []*model.NodeMeta,
	ge model.GetEdger,
) *allSegments {
	all := make([]pathSegment, 0, len(metas)/2)

	isSeen := make([]bool, len(metas))

	for i, nm := range metas {
		if isSeen[i] {
			continue
		}
		ps, cs, isLoop := getPathSegment(ge, nm.Coord())

		for j := i; j < len(metas); j++ {
			if cs.IsCoordSeen(metas[j].Coord()) {
				ps.numNodes++
				isSeen[j] = true
			}
		}

		if isLoop {
			return &allSegments{
				all:          nil,
				numNodesSeen: ps.numNodes,
			}
		}
		all = append(all, ps)
	}

	return &allSegments{
		all: all,
	}
}

func getPathSegment(
	ge model.GetEdger,
	nc model.NodeCoord,
) (pathSegment, state.CoordSeener, bool) {
	w := newWalker(ge, nc)
	ps, isLoop := w.walkSegment()
	return ps, w.seen, isLoop
}

func (as *allSegments) NumNodesInLoop() int {
	if !as.IsLoop() {
		return 0
	}
	return as.numNodesSeen
}

func (as *allSegments) IsLoop() bool {
	return as != nil && len(as.all) == 0
}

func (as *allSegments) GetUnknownEdge(
	d isDefineder,
) (model.EdgePair, model.State) {
	if as == nil || as.IsLoop() {
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

	newSegments := make([]pathSegment, 0, len(as.all))

	targets := state.NewCoordSeen(ge.NumEdges())
	for _, pc := range as.all {
		targets.Mark(pc.start.coord)
		targets.Mark(pc.end.coord)
	}

	isConnected := make([]bool, len(as.all))
	for i := 0; i < len(as.all); i++ {
		if isConnected[i] {
			continue
		}
		newPS, othersSeen, isLoop := extendOut(ge, as.all[i], targets, as.all[i+1:])

		nSeen := as.all[i].numNodes
		numSeen := 0
		for j, isSeen := range othersSeen {
			if isSeen {
				isConnected[i+1+j] = true
				nSeen += as.all[i+1+j].numNodes
				numSeen++
			}
		}

		if isLoop {
			return &allSegments{
				all:          nil,
				numNodesSeen: nSeen,
			}
		}

		newPS.numNodes = nSeen

		newSegments = append(newSegments, newPS)
	}

	return &allSegments{
		all: newSegments,
	}
}

func extendOut(
	ge model.GetEdger,
	ps pathSegment,
	targets state.CoordSeener,
	possibleConnections []pathSegment,
) (pathSegment, []bool, bool) {

	seenOthers := make([]bool, len(possibleConnections))
	seenPtrs := make([]*bool, len(possibleConnections))
	for i := range seenPtrs {
		seenPtrs[i] = &seenOthers[i]
	}

	sw := newWalker(ge, ps.end.coord)

	// find a defined outgoing edge from the start cap, and walk that path, until we find another segment
	newSC, isLoop := extendFrom(ps, ps.start, sw, ge, targets, possibleConnections, seenPtrs)
	if isLoop {
		// it's a loop!
		return pathSegment{}, seenOthers, true
	}

	// find a defined outgoing edge from the end cap, and walk that path, until we find another segment
	newEC, _ := extendFrom(ps, ps.end, sw, ge, targets, possibleConnections, seenPtrs)
	// if this were a loop, then we would have found it when extending the start
	// cap. As it is, this end cap _cannot_ be a loop.

	return pathSegment{
		start: newSC,
		end:   newEC,
	}, seenOthers, false
}

func extendFrom(
	ps pathSegment,
	sc segmentCap,
	w *simpleWalker,
	ge model.GetEdger,
	targets state.CoordSeener,
	possibleConnections []pathSegment,
	seenOthers []*bool,
) (segmentCap, bool) {

	startGoing := getNextEdge(ge, sc.coord, sc.outbound)
	if startGoing == model.HeadNowhere {
		// this end cap did not extend out at all. Just return it as it is.
		return sc, false
	}

	nc, lastMove, foundTarget := w.walkToTargets(sc.coord, startGoing, targets)
	if !foundTarget {
		return segmentCap{
			coord:    nc,
			outbound: lastMove.Opposite(),
		}, false
	}

	if nc == ps.start.coord || nc == ps.end.coord {
		// this is a loop
		return segmentCap{
			coord:    nc,
			outbound: lastMove.Opposite(),
		}, true
	}

	// Ok. so we foundTarget. That means, we need to find the possibleConnection
	// that has a start/end at nc. Then, we need to take that possibility, and extend
	// the opposite end. And we need to recurse that down until we've exhausted this
	// newfound segment
	for i, pc := range possibleConnections {
		if *seenOthers[i] {
			continue
		}
		if pc.start.coord == nc {
			*seenOthers[i] = true
			// extend out from the end.
			return extendFrom(ps, pc.end, w, ge, targets, possibleConnections, seenOthers)
		} else if pc.end.coord == nc {
			*seenOthers[i] = true
			// extend out from start.
			return extendFrom(ps, pc.start, w, ge, targets, possibleConnections, seenOthers)
		}
	}

	panic(`unexpected. dev should have been able to find connection`)
}
