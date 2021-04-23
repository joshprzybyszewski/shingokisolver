package model

import (
	"fmt"
)

var (
	InvalidEdgePair EdgePair = EdgePair{
		NodeCoord: InvalidNodeCoord,
		Cardinal:  HeadNowhere,
	}
)

type EdgePair struct {
	NodeCoord
	Cardinal
}

func NewEdgePair(
	nc NodeCoord,
	dir Cardinal,
) EdgePair {
	return EdgePair{
		NodeCoord: nc,
		Cardinal:  dir,
	}.Standardize()
}

func (ep EdgePair) Standardize() EdgePair {
	switch ep.Cardinal {
	case HeadUp, HeadLeft:
		return EdgePair{
			NodeCoord: ep.NodeCoord.Translate(ep.Cardinal),
			Cardinal:  ep.Cardinal.Opposite(),
		}

	default:
		// this is good.
		return ep
	}
}

func (ep EdgePair) Next(dir Cardinal) EdgePair {
	return EdgePair{
		NodeCoord: ep.NodeCoord.Translate(dir),
		Cardinal:  ep.Cardinal,
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
