package puzzle

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type EdgePair struct {
	model.NodeCoord
	model.Cardinal
}

func NewEdgePair(
	nc model.NodeCoord,
	dir model.Cardinal,
) EdgePair {
	switch dir {
	case model.HeadUp, model.HeadLeft:
		return EdgePair{
			NodeCoord: nc.Translate(dir),
			Cardinal:  dir.Opposite(),
		}
	case model.HeadRight, model.HeadDown:
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
	return fmt.Sprintf("coord: %+v, dir: %s", ep.NodeCoord, ep.Cardinal)
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
