package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTwoArmOptions(t *testing.T) {

	for val := int8(2); val <= 25; val++ {
		blackNode := NewNode(NewCoord(1, 1), false, val)
		expTwoArms := longBuildTwoArms(blackNode, 25)
		actTwoArms := BuildTwoArmOptions(blackNode, 25)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, blackNode)

		// black nodes can "point" in the four cardinal directions.
		// the first two segments are require for this pointing, so
		// the rest of the edges can contribute to this
		expNumOptions := 4 * (1 + int(val-2))
		assert.LessOrEqual(t, len(expTwoArms), expNumOptions, `unexpected num options for a black node: %+v`, blackNode)
		assert.LessOrEqual(t, len(actTwoArms), expNumOptions, `unexpected num options for a black node: %+v`, blackNode)

		whiteNode := NewNode(NewCoord(1, 1), true, val)
		expTwoArms = longBuildTwoArms(whiteNode, 25)
		actTwoArms = BuildTwoArmOptions(whiteNode, 25)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, whiteNode)

		expNumOptions = 2 * (1 + int(val-2))
		assert.LessOrEqual(t, len(expTwoArms), expNumOptions, `unexpected num options for a white node: %+v`, whiteNode)
		assert.LessOrEqual(t, len(actTwoArms), expNumOptions, `unexpected num options for a white node: %+v`, whiteNode)
	}
}

func TestGetAllEdges(t *testing.T) {
	ta := TwoArms{
		One: Arm{
			Heading: HeadRight,
			Len:     2,
		},
		Two: Arm{
			Heading: HeadDown,
			Len:     2,
		},
	}

	existing, avoided := ta.GetAllEdges(NewCoord(0, 0))

	expExisting := []EdgePair{
		NewEdgePair(NewCoord(0, 0), HeadRight),
		NewEdgePair(NewCoord(0, 1), HeadRight),
		NewEdgePair(NewCoord(0, 0), HeadDown),
		NewEdgePair(NewCoord(1, 0), HeadDown),
	}
	expAvoided := []EdgePair{
		NewEdgePair(NewCoord(0, 2), HeadRight),
		NewEdgePair(NewCoord(2, 0), HeadDown),
	}

	assert.Equal(t, expExisting, existing)
	assert.Equal(t, expAvoided, avoided)
}
