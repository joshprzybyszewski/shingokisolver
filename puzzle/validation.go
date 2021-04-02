package puzzle

import (
	"errors"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) IsRangeInvalid(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) bool {
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

	return p.isRangeInvalid(startR, stopR, startC, stopC)
}

func (p *Puzzle) isRangeInvalid(
	startR, stopR model.RowIndex,
	startC, stopC model.ColIndex,
) bool {
	for r := startR; r < stopR; r++ {
		for c := startC; c < stopC; c++ {
			if p.isCoordInvalid(model.NewCoord(r, c)) {
				return true
			}
		}
	}

	// it's cheaper for us to just iterate all of the nodes
	// and check for their validity than it is to check every
	// (r, c) or filtering out to only be in the range
	for nc, n := range p.nodes {
		// if this point is a node, check for if it's invalid
		oe, ok := p.GetOutgoingEdgesFrom(nc)
		if !ok || n.IsInvalid(oe) {
			return true
		}
	}

	return false
}

func (p *Puzzle) isCoordInvalid(
	coord model.NodeCoord,
) bool {
	// check that this point doesn't branch
	oe, ok := p.GetOutgoingEdgesFrom(coord)
	// either we can't get the node (the coordinate
	// must be out of bounds or this node branches.
	// therefore, this Puzzle is invalid
	return !ok || oe.GetNumOutgoingDirections() > 2
}

func (p *Puzzle) checkIsInvalidFrom(
	coord model.NodeCoord,
) bool {
	if p.isCoordInvalid(coord) {
		return true
	}

	// TODO iterate out from coord to find out if it's invalid
	return p.isRangeInvalid(
		0,
		model.RowIndex(p.numNodes()),
		0,
		model.ColIndex(p.numNodes()),
	)
}

func (p *Puzzle) IsInViolation(
	coord model.NodeCoord,
) (bool, error) {
	if p.checkIsInvalidFrom(coord) {
		return true, errors.New(`invalid Puzzle`)
	}

	oe, ok := p.GetOutgoingEdgesFrom(coord)
	if !ok {
		return true, errors.New(`bad input`)
	}

	return oe.GetNumOutgoingDirections() != 2, nil
}

func (p *Puzzle) IsIncomplete(
	coord model.NodeCoord,
) (bool, error) {

	if violates, err := p.IsInViolation(coord); err != nil || violates {
		return violates, err
	}

	w := newWalker(p, coord)
	seenNodes, ok := w.walk()
	if !ok {
		// our path did not make it all the way around
		return true, nil
	}

	for nc, n := range p.nodes {
		if _, ok := seenNodes[nc]; !ok {
			// node was not seen
			return true, errors.New(`this path made a loop, but didn't see every node`)
		}

		oe, ok := p.GetOutgoingEdgesFrom(nc)
		if !ok {
			return true, errors.New(`bad input`)
		}
		if oe.TotalEdges() != n.Value() {
			// previously (in isRangeInvalid) we checked if oe.TotalEdges() > n.val
			// This check exists to verify we have exactly how many we need.
			return true, nil
		}

	}

	// at this point, we have a valid board, our path is a loop,
	// and we've seen all of the nodes appropriately. Therefore,
	// our board is not incomplete, and it's a solution.
	return false, nil
}
