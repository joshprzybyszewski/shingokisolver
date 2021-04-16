package model

import (
	"fmt"
)

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
	for _, o := range others {
		if o == ep {
			return true
		}
	}
	return false
}
