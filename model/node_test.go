package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeTypeIsInvalidMotions(t *testing.T) {
	testCases := []struct {
		c1       Cardinal
		c2       Cardinal
		expWhite bool
		expBlack bool
	}{{
		c1:       HeadUp,
		c2:       HeadRight,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadUp,
		c2:       HeadDown,
		expWhite: false,
		expBlack: true,
	}, {
		c1:       HeadUp,
		c2:       HeadLeft,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadRight,
		c2:       HeadDown,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadRight,
		c2:       HeadLeft,
		expWhite: false,
		expBlack: true,
	}, {
		c1:       HeadDown,
		c2:       HeadLeft,
		expWhite: true,
		expBlack: false,
	}}

	for _, tc := range testCases {
		assert.True(t, noNode.isInvalidMotions(tc.c1, tc.c2))
		assert.True(t, noNode.isInvalidMotions(tc.c2, tc.c1))
		assert.Equal(t, tc.expWhite, WhiteNode.isInvalidMotions(tc.c1, tc.c2))
		assert.Equal(t, tc.expWhite, WhiteNode.isInvalidMotions(tc.c2, tc.c1))
		assert.Equal(t, tc.expBlack, BlackNode.isInvalidMotions(tc.c1, tc.c2))
		assert.Equal(t, tc.expBlack, BlackNode.isInvalidMotions(tc.c2, tc.c1))
	}
}

func TestIsInvalidMotions(t *testing.T) {
	w5 := NewNode(NodeCoord{}, true, 5)

	assert.True(t, w5.IsInvalidMotions(HeadUp, HeadRight))
	assert.True(t, w5.IsInvalidMotions(HeadUp, HeadLeft))
	assert.True(t, w5.IsInvalidMotions(HeadDown, HeadRight))
	assert.True(t, w5.IsInvalidMotions(HeadDown, HeadLeft))
	assert.True(t, w5.IsInvalidMotions(HeadRight, HeadUp))
	assert.True(t, w5.IsInvalidMotions(HeadRight, HeadDown))
	assert.True(t, w5.IsInvalidMotions(HeadLeft, HeadUp))
	assert.True(t, w5.IsInvalidMotions(HeadLeft, HeadDown))

	assert.False(t, w5.IsInvalidMotions(HeadUp, HeadDown))
	assert.False(t, w5.IsInvalidMotions(HeadDown, HeadUp))
	assert.False(t, w5.IsInvalidMotions(HeadLeft, HeadRight))
	assert.False(t, w5.IsInvalidMotions(HeadRight, HeadLeft))

	b7 := NewNode(NodeCoord{}, false, 7)
	assert.False(t, b7.IsInvalidMotions(HeadUp, HeadRight))
	assert.False(t, b7.IsInvalidMotions(HeadUp, HeadLeft))
	assert.False(t, b7.IsInvalidMotions(HeadDown, HeadRight))
	assert.False(t, b7.IsInvalidMotions(HeadDown, HeadLeft))
	assert.False(t, b7.IsInvalidMotions(HeadRight, HeadUp))
	assert.False(t, b7.IsInvalidMotions(HeadRight, HeadDown))
	assert.False(t, b7.IsInvalidMotions(HeadLeft, HeadUp))
	assert.False(t, b7.IsInvalidMotions(HeadLeft, HeadDown))

	assert.True(t, b7.IsInvalidMotions(HeadUp, HeadDown))
	assert.True(t, b7.IsInvalidMotions(HeadDown, HeadUp))
	assert.True(t, b7.IsInvalidMotions(HeadLeft, HeadRight))
	assert.True(t, b7.IsInvalidMotions(HeadRight, HeadLeft))
}

func TestGetFilteredOptionsSimpleAvoids(t *testing.T) {
	nEdges := 15
	b3 := NewNode(NewCoord(7, 7), false, 3)
	allTAs := BuildTwoArmOptions(b3, nEdges)

	tge := testGetEdger{
		numEdges: nEdges,
	}
	filtered := b3.GetFilteredOptions(allTAs, tge, nil)
	assert.Equal(t, filtered, allTAs)

	tge = testGetEdger{
		numEdges: nEdges,
		avoided: []EdgePair{
			NewEdgePair(NewCoord(7, 7), HeadUp),
		},
	}
	filtered = b3.GetFilteredOptions(allTAs, tge, nil)
	assert.NotEqual(t, filtered, allTAs)
	for _, ta := range filtered {
		assert.NotEqual(t, HeadUp, ta.One.Heading)
		assert.NotEqual(t, HeadUp, ta.Two.Heading)
	}

	tge = testGetEdger{
		numEdges: nEdges,
		avoided: []EdgePair{
			NewEdgePair(NewCoord(7, 7), HeadUp),
			NewEdgePair(NewCoord(7, 7), HeadLeft),
		},
	}
	filtered = b3.GetFilteredOptions(allTAs, tge, nil)
	assert.NotEqual(t, filtered, allTAs)
	for i, ta := range filtered {
		assert.NotEqual(t, HeadUp, ta.One.Heading, `issue with TwoArms at index %d: %v`, i, ta)
		assert.NotEqual(t, HeadLeft, ta.One.Heading, `issue with TwoArms at index %d: %v`, i, ta)
		assert.NotEqual(t, HeadUp, ta.Two.Heading, `issue with TwoArms at index %d: %v`, i, ta)
		assert.NotEqual(t, HeadLeft, ta.Two.Heading, `issue with TwoArms at index %d: %v`, i, ta)

		if ta.One.Heading != HeadRight {
			// one of the two needs to HeadRight
			assert.Equal(t, HeadRight, ta.Two.Heading, `issue with TwoArms at index %d: %v`, i, ta)
		}
		if ta.One.Heading != HeadDown {
			// one of the two needs to HeadDown
			assert.Equal(t, HeadDown, ta.Two.Heading, `issue with TwoArms at index %d: %v`, i, ta)
		}
	}
}
