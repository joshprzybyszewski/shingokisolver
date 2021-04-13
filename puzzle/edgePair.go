package puzzle

import (
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

type edgePair struct {
	coord model.NodeCoord
	dir   model.Cardinal
}

func (ep edgePair) String() string {
	return fmt.Sprintf("coord: %+v, dir: %s", ep.coord, ep.dir)
}

func standardizeInput(
	nc model.NodeCoord,
	dir model.Cardinal,
) (edgePair, error) {
	switch dir {
	case model.HeadUp, model.HeadLeft:
		return edgePair{
			coord: nc.Translate(dir),
			dir:   dir.Opposite(),
		}, nil
	case model.HeadRight, model.HeadDown:
		// this is good.
		return edgePair{
			coord: nc,
			dir:   dir,
		}, nil
	default:
		return edgePair{}, errors.New(`unexpected cardinal`)
	}
}
