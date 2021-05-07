package state

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

type bitData uint32

const (
	MaxEdges = 32 // needs to have the space for it in bitData
)

var _ model.GetEdger = (*TriEdges)(nil)

type TriEdges struct {
	rows []bitData
	cols []bitData

	avoidRows []bitData
	avoidCols []bitData

	numEdges uint16
}

func New(
	numEdges int,
) TriEdges {
	return TriEdges{
		numEdges:  uint16(numEdges),
		rows:      make([]bitData, numEdges+1),
		cols:      make([]bitData, numEdges+1),
		avoidRows: make([]bitData, numEdges+1),
		avoidCols: make([]bitData, numEdges+1),
	}
}

func (ets *TriEdges) NumEdges() int {
	return int(ets.numEdges)
}

func (ets *TriEdges) isArmInBounds(
	start model.NodeCoord,
	arm model.Arm,
) bool {
	if start.Row < 0 ||
		start.Col < 0 ||
		uint16(start.Row) > ets.numEdges ||
		uint16(start.Col) > ets.numEdges {
		return false
	}
	switch arm.Heading {
	case model.HeadUp:
		return start.Row-model.RowIndex(arm.Len) >= 0
	case model.HeadDown:
		return uint16(start.Row+model.RowIndex(arm.Len)) <= ets.numEdges
	case model.HeadLeft:
		return start.Col-model.ColIndex(arm.Len) >= 0
	case model.HeadRight:
		return uint16(start.Col+model.ColIndex(arm.Len)) <= ets.numEdges
	}

	return false
}

func (ets *TriEdges) IsInBounds(
	ep model.EdgePair,
) bool {
	if ep.Row < 0 || ep.Col < 0 {
		// negative coords are bad
		return false
	}
	switch ep.Cardinal {
	case model.HeadRight:
		return uint16(ep.Row) <= ets.numEdges && uint16(ep.Col) < ets.numEdges
	case model.HeadDown:
		return uint16(ep.Row) < ets.numEdges && uint16(ep.Col) <= ets.numEdges
	default:
		// unexpected input
		return false
	}
}

func (ets *TriEdges) AllExist(
	start model.NodeCoord,
	arm model.Arm,
) bool {
	if !ets.isArmInBounds(start, arm) {
		// the arm goes out of bounds, so we know it can't be all in
		return false
	}
	mask := getMask(start, arm)

	switch arm.Heading {
	case model.HeadRight:
		return ets.rows[start.Row]&mask == mask

	case model.HeadDown:
		return ets.cols[start.Col]&mask == mask

	case model.HeadLeft:
		return ets.rows[start.Row]&mask == mask

	case model.HeadUp:
		return ets.cols[start.Col]&mask == mask

	default:
		return false
	}
}

func (ets *TriEdges) AnyAvoided(
	start model.NodeCoord,
	arm model.Arm,
) bool {
	if !ets.isArmInBounds(start, arm) {
		return true
	}
	mask := getMask(start, arm)

	switch arm.Heading {
	case model.HeadRight:
		return ets.avoidRows[start.Row]&mask != 0

	case model.HeadDown:
		return ets.avoidCols[start.Col]&mask != 0

	case model.HeadLeft:
		return ets.avoidRows[start.Row]&mask != 0

	case model.HeadUp:
		return ets.avoidCols[start.Col]&mask != 0

	default:
		return false
	}
}

func (ets *TriEdges) Any(
	start model.NodeCoord,
	arm model.Arm,
) (bool, bool) {
	goesOutOfBounds := !ets.isArmInBounds(start, arm)
	mask := getMask(start, arm)

	switch arm.Heading {
	case model.HeadRight:
		return ets.rows[start.Row]&mask != 0, goesOutOfBounds || ets.avoidRows[start.Row]&mask != 0

	case model.HeadDown:
		return ets.cols[start.Col]&mask != 0, goesOutOfBounds || ets.avoidCols[start.Col]&mask != 0

	case model.HeadLeft:
		return ets.rows[start.Row]&mask != 0, goesOutOfBounds || ets.avoidRows[start.Row]&mask != 0

	case model.HeadUp:
		return ets.cols[start.Col]&mask != 0, goesOutOfBounds || ets.avoidCols[start.Col]&mask != 0

	default:
		return false, false
	}
}

func (ets *TriEdges) IsEdge(ep model.EdgePair) bool {
	if !ets.IsInBounds(ep) {
		return false
	}

	switch ep.Cardinal {
	case model.HeadRight:
		return (ets.rows[ep.Row] & masks[ep.Col]) != 0
	case model.HeadDown:
		return (ets.cols[ep.Col] & masks[ep.Row]) != 0
	}

	// unexpected!
	return false
}

func (ets *TriEdges) IsAvoided(ep model.EdgePair) bool {
	if !ets.IsInBounds(ep) {
		return true
	}

	switch ep.Cardinal {
	case model.HeadRight:
		return (ets.avoidRows[ep.Row] & masks[ep.Col]) != 0
	case model.HeadDown:
		return (ets.avoidCols[ep.Col] & masks[ep.Row]) != 0
	}

	return false
}

func (ets *TriEdges) IsDefined(ep model.EdgePair) bool {
	if !ets.IsInBounds(ep) {
		return true
	}

	switch ep.Cardinal {
	case model.HeadRight:
		return ((ets.avoidRows[ep.Row] | ets.rows[ep.Row]) & masks[ep.Col]) != 0
	case model.HeadDown:
		return ((ets.avoidCols[ep.Col] | ets.cols[ep.Col]) & masks[ep.Row]) != 0
	}

	// unexpected!
	return false
}

func (ets *TriEdges) GetEdge(
	ep model.EdgePair,
) model.EdgeState {

	if !ets.IsInBounds(ep) {
		return model.EdgeOutOfBounds
	}

	switch ep.Cardinal {
	case model.HeadRight:
		if (ets.rows[ep.Row] & masks[ep.Col]) != 0 {
			return model.EdgeExists
		}
		if (ets.avoidRows[ep.Row] & masks[ep.Col]) != 0 {
			return model.EdgeAvoided
		}
	case model.HeadDown:
		if (ets.cols[ep.Col] & masks[ep.Row]) != 0 {
			return model.EdgeExists
		}
		if (ets.avoidCols[ep.Col] & masks[ep.Row]) != 0 {
			return model.EdgeAvoided
		}
	default:
		return model.EdgeErrored
	}

	return model.EdgeUnknown
}

func (ets *TriEdges) UpdateEdge(
	ep model.EdgePair,
	es model.EdgeState,
) model.State {
	switch es {
	case model.EdgeAvoided:
		return ets.avoidEdge(ep)
	case model.EdgeExists:
		return ets.setEdge(ep)
	}
	return model.Violation
}

func (ets *TriEdges) setEdge(
	ep model.EdgePair,
) model.State {

	switch ets.GetEdge(ep) {
	case model.EdgeExists:
		return model.Duplicate
	case model.EdgeAvoided, model.EdgeOutOfBounds:
		return model.Violation
	case model.EdgeErrored:
		return model.Unexpected
	}

	switch ep.Cardinal {
	case model.HeadRight:
		ets.rows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		ets.cols[ep.Col] |= masks[ep.Row]
	}

	return model.Incomplete
}

func (ets *TriEdges) avoidEdge(
	ep model.EdgePair,
) model.State {

	switch ets.GetEdge(ep) {
	case model.EdgeAvoided, model.EdgeOutOfBounds:
		// we can avoid edges that are out of bounds
		// we just know that they're out of bounds!
		return model.Duplicate
	case model.EdgeExists:
		return model.Violation
	case model.EdgeErrored:
		return model.Unexpected
	}

	switch ep.Cardinal {
	case model.HeadRight:
		ets.avoidRows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		ets.avoidCols[ep.Col] |= masks[ep.Row]
	}

	return model.Incomplete
}

func (ets TriEdges) Copy() TriEdges {
	cpy := TriEdges{
		numEdges:  ets.numEdges,
		rows:      make([]bitData, len(ets.rows)),
		cols:      make([]bitData, len(ets.cols)),
		avoidRows: make([]bitData, len(ets.avoidRows)),
		avoidCols: make([]bitData, len(ets.avoidCols)),
	}

	copy(cpy.rows, ets.rows)
	copy(cpy.cols, ets.cols)
	copy(cpy.avoidRows, ets.avoidRows)
	copy(cpy.avoidCols, ets.avoidCols)

	return cpy
}
