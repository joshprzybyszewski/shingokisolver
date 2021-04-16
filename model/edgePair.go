package model

import (
	"fmt"
)

type GetEdger interface {
	GetEdge(EdgePair) EdgeState

	AllExist(NodeCoord, Arm) bool
	Any(NodeCoord, Arm) (bool, bool)
}

type EdgePair struct {
	NodeCoord
	Cardinal
}

func NewEdgePair(
	nc NodeCoord,
	dir Cardinal,
) EdgePair {
	switch dir {
	case HeadUp, HeadLeft:
		return EdgePair{
			NodeCoord: nc.Translate(dir),
			Cardinal:  dir.Opposite(),
		}

	case HeadRight, HeadDown:
		// this is good.
		return EdgePair{
			NodeCoord: nc,
			Cardinal:  dir,
		}

	default:
		panic(`unexpected cardinal: ` + dir.String())
	}
}

func (ep EdgePair) String() string {
	return fmt.Sprintf("EdgePair{coord: %+v, dir: %s}", ep.NodeCoord, ep.Cardinal)
}

func (ep EdgePair) IsIn(
	others ...EdgePair,
) bool {
	return ep.IndexOf(others...) >= 0
}

func (ep EdgePair) IndexOf(
	others ...EdgePair,
) int {
	for i, o := range others {
		if o == ep {
			return i
		}
	}
	return -1
}
