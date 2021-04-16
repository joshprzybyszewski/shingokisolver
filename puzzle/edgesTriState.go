package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	MaxEdges = 25 // currently, the len(masks)
)

var _ model.GetEdger = (*edgesTriState)(nil)

type bitData uint32

var (
	masks       = make([]bitData, MaxEdges)
	armLenMasks = make([]bitData, MaxEdges+1)
)

func init() {
	buildMasks()
	buildArmLenMasks()
}

func buildMasks() {
	for i := 0; i < MaxEdges; i++ {
		masks[i] = 1 << i
	}
}

func buildArmLenMasks() {
	for i := 0; i < len(armLenMasks); i++ {
		var lenMask bitData
		for mi := 0; mi < i; mi++ {
			lenMask |= masks[mi]
		}
		armLenMasks[i] = lenMask
	}
}

func getMask(
	start model.NodeCoord,
	arm model.Arm,
) bitData {
	switch arm.Heading {
	case model.HeadRight:
		return armLenMasks[arm.Len] << start.Col

	case model.HeadDown:
		return armLenMasks[arm.Len] << start.Row

	case model.HeadLeft:
		lmi := int(arm.Len)
		shift := start.Col - model.ColIndex(arm.Len)
		if shift < 0 {
			lmi += int(shift)
			if lmi < 0 {
				return 0
			}
			shift = 0
		}
		return armLenMasks[lmi] << shift

	case model.HeadUp:
		lmi := int(arm.Len)
		shift := start.Row - model.RowIndex(arm.Len)
		if shift < 0 {
			lmi += int(shift)
			if lmi < 0 {
				return 0
			}
			shift = 0
		}
		return armLenMasks[lmi] << shift

	default:
		return 0
	}
}

type edgesTriState struct {
	numEdges uint16

	rows []bitData
	cols []bitData

	avoidRows []bitData
	avoidCols []bitData
}

func newEdgesStates(
	numEdges int,
) *edgesTriState {
	return &edgesTriState{
		numEdges:  uint16(numEdges),
		rows:      make([]bitData, numEdges+1),
		cols:      make([]bitData, numEdges+1),
		avoidRows: make([]bitData, numEdges+1),
		avoidCols: make([]bitData, numEdges+1),
	}
}

func (ets *edgesTriState) isArmInBounds(
	start model.NodeCoord,
	arm model.Arm,
) bool {
	if start.Row < 0 ||
		start.Col < 0 ||
		start.Row > model.RowIndex(ets.numEdges) ||
		start.Col > model.ColIndex(ets.numEdges) {
		return false
	}
	switch arm.Heading {
	case model.HeadUp:
		return start.Row-model.RowIndex(arm.Len) >= 0
	case model.HeadDown:
		return start.Row+model.RowIndex(arm.Len) <= model.RowIndex(ets.numEdges)
	case model.HeadLeft:
		return start.Col-model.ColIndex(arm.Len) >= 0
	case model.HeadRight:
		return start.Col+model.ColIndex(arm.Len) <= model.ColIndex(ets.numEdges)
	}

	return false
}

func (ets *edgesTriState) isInBounds(
	ep model.EdgePair,
) bool {
	if ep.Row < 0 || ep.Col < 0 {
		// negative coords are bad
		return false
	}
	switch ep.Cardinal {
	case model.HeadRight:
		return ep.Row <= model.RowIndex(ets.numEdges) && ep.Col < model.ColIndex(ets.numEdges)
	case model.HeadDown:
		return ep.Row < model.RowIndex(ets.numEdges) && ep.Col <= model.ColIndex(ets.numEdges)
	default:
		// unexpected input
		return false
	}
}

func (ets *edgesTriState) AllExist(
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

func (ets *edgesTriState) AnyAvoided(
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

func (ets *edgesTriState) Any(
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

func (ets *edgesTriState) GetEdge(
	ep model.EdgePair,
) model.EdgeState {

	if !ets.isInBounds(ep) {
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

func (ets *edgesTriState) SetEdge(
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

func (ets *edgesTriState) AvoidEdge(
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

func (ets *edgesTriState) Copy() *edgesTriState {
	cpy := &edgesTriState{
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
