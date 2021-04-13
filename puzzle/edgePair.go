package puzzle

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type edgePair struct {
	model.NodeCoord
	model.Cardinal
}

func (ep edgePair) String() string {
	return fmt.Sprintf("coord: %+v, dir: %s", ep.NodeCoord, ep.Cardinal)
}

func newEdgePair(
	nc model.NodeCoord,
	dir model.Cardinal,
) edgePair {
	switch dir {
	case model.HeadUp, model.HeadLeft:
		return edgePair{
			NodeCoord: nc.Translate(dir),
			Cardinal:  dir.Opposite(),
		}
	case model.HeadRight, model.HeadDown:
		// this is good.
		return edgePair{
			NodeCoord: nc,
			Cardinal:  dir,
		}
	default:
		panic(`unexpected cardinal: ` + dir.String())
	}
}
