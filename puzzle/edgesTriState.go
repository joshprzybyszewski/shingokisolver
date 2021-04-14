package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	MaxEdges = 25 // currently, the len(masks)
)

var _ getEdger = (*edgesTriState)(nil)

var (
	masks = []uint64{
		1 << 0,
		1 << 1,
		1 << 2,
		1 << 3,
		1 << 4,
		1 << 5,
		1 << 6,
		1 << 7,
		1 << 8,
		1 << 9,
		1 << 10,
		1 << 11,
		1 << 12,
		1 << 13,
		1 << 14,
		1 << 15,
		1 << 16,
		1 << 17,
		1 << 18,
		1 << 19,
		1 << 20,
		1 << 21,
		1 << 22,
		1 << 23,
		1 << 24,
		1 << 25,
	}
)

type edgesTriState struct {
	numEdges uint8

	rows []uint64
	cols []uint64

	avoidRows []uint64
	avoidCols []uint64
}

func newEdgesBits(
	numEdges uint8,
) *edgesTriState {
	return &edgesTriState{
		numEdges:  numEdges,
		rows:      make([]uint64, numEdges+1),
		cols:      make([]uint64, numEdges+1),
		avoidRows: make([]uint64, numEdges+1),
		avoidCols: make([]uint64, numEdges+1),
	}
}

func (eb *edgesTriState) isInBounds(
	ep EdgePair,
) bool {
	if ep.Row < 0 || ep.Col < 0 {
		// negative coords are bad
		return false
	}
	switch ep.Cardinal {
	case model.HeadRight:
		return uint8(ep.Row) <= eb.numEdges && uint8(ep.Col) < eb.numEdges
	case model.HeadDown:
		return uint8(ep.Row) < eb.numEdges && uint8(ep.Col) <= eb.numEdges
	default:
		// unexpected input
		return false
	}
}

func (eb *edgesTriState) GetEdge(
	ep EdgePair,
) model.EdgeState {

	if !eb.isInBounds(ep) {
		return model.EdgeOutOfBounds
	}

	switch ep.Cardinal {
	case model.HeadRight:
		if (eb.rows[ep.Row] & masks[ep.Col]) != 0 {
			return model.EdgeExists
		}
		if (eb.avoidRows[ep.Row] & masks[ep.Col]) != 0 {
			return model.EdgeAvoided
		}
	case model.HeadDown:
		if (eb.cols[ep.Col] & masks[ep.Row]) != 0 {
			return model.EdgeExists
		}
		if (eb.avoidCols[ep.Col] & masks[ep.Row]) != 0 {
			return model.EdgeAvoided
		}
	default:
		return model.EdgeErrored
	}

	return model.EdgeUnknown
}

func (eb *edgesTriState) SetEdge(
	ep EdgePair,
) model.State {

	switch eb.GetEdge(ep) {
	case model.EdgeExists:
		return model.Duplicate
	case model.EdgeAvoided, model.EdgeOutOfBounds:
		return model.Violation
	case model.EdgeErrored:
		return model.Unexpected
	}

	switch ep.Cardinal {
	case model.HeadRight:
		eb.rows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		eb.cols[ep.Col] |= masks[ep.Row]
	}

	return model.Incomplete
}

func (eb *edgesTriState) AvoidEdge(
	ep EdgePair,
) model.State {

	switch eb.GetEdge(ep) {
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
		eb.avoidRows[ep.Row] |= masks[ep.Col]
	case model.HeadDown:
		eb.avoidCols[ep.Col] |= masks[ep.Row]
	}

	return model.Incomplete
}

func (eb *edgesTriState) Copy() *edgesTriState {
	cpy := &edgesTriState{
		numEdges:  eb.numEdges,
		rows:      make([]uint64, len(eb.rows)),
		cols:      make([]uint64, len(eb.cols)),
		avoidRows: make([]uint64, len(eb.avoidRows)),
		avoidCols: make([]uint64, len(eb.avoidCols)),
	}

	copy(cpy.rows, eb.rows)
	copy(cpy.cols, eb.cols)
	copy(cpy.avoidRows, eb.avoidRows)
	copy(cpy.avoidCols, eb.avoidCols)

	return cpy
}
