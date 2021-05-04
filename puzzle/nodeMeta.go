package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type nodeMeta struct {
	node          model.Node
	nearby        model.NearbyNodes
	twoArmOptions []model.TwoArms

	isComplete bool
}
