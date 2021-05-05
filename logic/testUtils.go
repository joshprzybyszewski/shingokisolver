package logic

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func AssertHasAdvancedNode(
	t *testing.T,
	r *Rules,
	expNearbyNodes model.NearbyNodes,
	expOptions []model.TwoArms,
	expNode model.Node,
	expDir model.Cardinal,
	expIndex int8,
) {
	// TODO remove usages...
}
