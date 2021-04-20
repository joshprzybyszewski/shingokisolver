package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

var _ evaluator = afterArm{}

type afterArm struct {
	nc model.NodeCoord
	ta model.TwoArms
}

func newAfterArmEvaluator(
	nc model.NodeCoord,
	ta model.TwoArms,
) evaluator {
	return afterArm{
		nc: nc,
		ta: ta,
	}
}

func (ac afterArm) evaluate(ge model.GetEdger) model.EdgeState {
	if !ge.AllExist(ac.nc, ac.ta.One) {
		// arm1 doesn't completely exist.
		// I claim no knowledge
		return model.EdgeUnknown
	}

	if ge.AllExist(ac.nc, ac.ta.Two) {
		// arm1 and arm2 both exist. I know I'm avoided.
		return model.EdgeAvoided
	}

	return model.EdgeUnknown
}
