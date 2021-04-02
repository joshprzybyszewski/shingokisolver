package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) GetState() model.State {
	return p.getRangeState(
		0,
		model.RowIndex(p.numNodes()),
		0,
		model.ColIndex(p.numNodes()),
	)
}

func (p *Puzzle) GetRangeState(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) model.State {
	if startR < 0 {
		startR = 0
	}
	if maxR := model.RowIndex(p.numNodes()); stopR > maxR {
		stopR = maxR
	}
	if startC < 0 {
		startC = 0
	}
	if maxC := model.ColIndex(p.numNodes()); stopC > maxC {
		stopC = maxC
	}

	return p.getRangeState(
		startR, stopR, startC, stopC,
	)
}

func (p *Puzzle) getRangeState(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) model.State {

	switch nodeState := p.getStateOfNodes(); nodeState {
	case model.Incomplete, model.Complete:
		// keep going through checks...
	default:
		return nodeState
	}

	var coord model.NodeCoord
	for nc := range p.nodes {
		// just need a random starting node for the walker
		coord = nc
		break
	}

	w := newWalker(p, coord)
	seenNodes, walkerState := w.walk()
	switch walkerState {
	case model.Complete:
		// keep going through checks...
	default:
		return walkerState
	}

	for nc := range p.nodes {
		if _, ok := seenNodes[nc]; !ok {
			// node was not seen. therefore, we completed a loop that
			// doesn't see all nodes!
			return model.Violation
		}
	}

	return model.Complete
}

func (p *Puzzle) getStateOfNodes() model.State {
	hasIncompleteNodes := false
	// it's cheaper for us to just iterate all of the nodes
	// and check for their validity than it is to check every
	// (r, c) or filtering out to only be in the range
	for nc, n := range p.nodes {
		oe, ok := p.GetOutgoingEdgesFrom(nc)
		if !ok {
			// something really weird happened
			panic(`da heck`)
			return model.Unexpected
		}
		switch s := n.GetState(oe); s {
		case model.Violation, model.Unexpected:
			return s
		case model.Incomplete:
			hasIncompleteNodes = true
		}
	}

	if hasIncompleteNodes {
		return model.Incomplete
	}

	return model.Complete
}

func (p *Puzzle) getStateForCoord(
	coord model.NodeCoord,
) model.State {
	// check that this point doesn't branch
	oe, ok := p.GetOutgoingEdgesFrom(coord)
	if !ok {
		return model.Violation
	}

	switch numEdges := oe.GetNumOutgoingDirections(); {
	case numEdges > 2:
		return model.Violation
	case numEdges == 2:
		return model.Complete
	default:
		return model.Incomplete
	}
}
