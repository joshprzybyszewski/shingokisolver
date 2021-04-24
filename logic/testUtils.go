package logic

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
)

func AssertHasAdvancedNode(
	t *testing.T,
	r *Rules,
	expNearbyNodes map[model.Cardinal][]*model.Node,
	expOptions []model.TwoArms,
	expNode model.Node,
	expDir model.Cardinal,
	expIndex int8,
) {
	for _, other := range r.otherEvals {
		an, ok := other.(advancedNode)
		if !ok {
			continue
		}
		if an.node != expNode {
			continue
		}
		assert.Equal(t, expIndex, an.index)
		assert.Equal(t, expDir, an.dir)
		assert.Equal(t, expNearbyNodes, an.nearbyNodes)
		assert.Equal(t, expOptions, an.options)
		return
	}
	assert.Fail(t, `did not find advnacedNode evaluator`)
}
