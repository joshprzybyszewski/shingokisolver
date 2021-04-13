package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTwoArmOptions(t *testing.T) {

	for val := int8(2); val <= 50; val++ {
		blackNode := NewNode(NewCoord(1, 1), false, val)
		expTwoArms := longBuildTwoArms(blackNode)
		actTwoArms := BuildTwoArmOptions(blackNode)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, blackNode)

		// black nodes can "point" in the four cardinal directions.
		// the first two segments are require for this pointing, so
		// the rest of the edges can contribute to this
		expNumOptions := 4 * (1 + int(val-2))
		assert.Equal(t, expNumOptions, len(expTwoArms), `unexpected num options for a black node: %+v`, blackNode)
		assert.Equal(t, expNumOptions, len(actTwoArms), `unexpected num options for a black node: %+v`, blackNode)

		whiteNode := NewNode(NewCoord(1, 1), true, val)
		expTwoArms = longBuildTwoArms(whiteNode)
		actTwoArms = BuildTwoArmOptions(whiteNode)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, whiteNode)

		expNumOptions = 2 * (1 + int(val-2))
		assert.Equal(t, expNumOptions, len(expTwoArms), `unexpected num options for a white node: %+v`, whiteNode)
		assert.Equal(t, expNumOptions, len(actTwoArms), `unexpected num options for a white node: %+v`, whiteNode)
	}
}
