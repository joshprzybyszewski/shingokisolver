package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func TestGimmesPuzzle90104(t *testing.T) {
	fresh := NewPuzzle(25, []model.NodeLocation{{
		Row:     0,
		Col:     13,
		IsWhite: true,
		Value:   9,
	}, {
		Row:     0,
		Col:     16,
		IsWhite: false,
		Value:   5,
	}})

	puzz, s := ClaimGimmes(fresh)
	require.Equal(t, model.Incomplete, s)
	t.Logf("puzz: \n%s\n", puzz)

	// inspect an edge's rules to verify we built it correctly.
	r := puzz.rules.Get(model.NewEdgePair(model.NewCoord(0, 11), model.HeadRight))
	require.NotNil(t, r)
	logic.AssertHasAdvancedNode(
		t,
		r,
		model.NearbyNodes{
			nil, // HeadNowhere
			{
				nil,
			}, // HeadRight
			nil, // HeadUp
			{
				nil,
			}, // HeadLeft
			nil, // HeadDown
		},
		[]model.TwoArms{{
			Two: model.Arm{
				Heading: model.HeadLeft,
				Len:     8,
			},
			One: model.Arm{
				Heading: model.HeadRight,
				Len:     1,
			},
		}, {
			Two: model.Arm{
				Heading: model.HeadLeft,
				Len:     7,
			},
			One: model.Arm{
				Heading: model.HeadRight,
				Len:     2,
			},
		}},
		model.NewNode(model.NewCoord(0, 13), true, 9),
		model.HeadLeft,
		1,
	)

	// The white node should have the two going straight through it.
	// This should be very easy.
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(0, 13)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 13)))

	// The black node should have an edge going down.
	// This should be very easy.
	assert.True(t, puzz.IsEdge(model.HeadDown, model.NewCoord(0, 16)))

	// The white node should also be extending to the left, otherwise
	// it will run into the black node and won't work!
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 12)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 11)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 10)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 9)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 8)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(0, 7)))

	// The black node should not extend right. It could go one edge
	// to the left, and then 4 down.
	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(0, 16)))
}

func TestGimmesPuzzle5817105(t *testing.T) {
	fresh := NewPuzzle(25, []model.NodeLocation{{
		Row:     24,
		Col:     0,
		Value:   7,
		IsWhite: true,
	}, {
		Row:     25,
		Col:     3,
		Value:   6,
		IsWhite: true,
	}, {
		Row:     25,
		Col:     8,
		Value:   2,
		IsWhite: false,
	}})

	puzz, s := ClaimGimmes(fresh)
	require.Equal(t, model.Incomplete, s)
	t.Logf("puzz: \n%s\n", puzz)

	// inspect an edge's rules to verify we built it correctly.
	r := puzz.rules.Get(model.NewEdgePair(model.NewCoord(25, 1), model.HeadRight))
	require.NotNil(t, r)
	logic.AssertHasAdvancedNode(
		t,
		r,
		model.NearbyNodes{
			nil, // HeadNowhere
			{
				nil,
			}, // HeadRight
			nil, // HeadUp
			{
				nil,
			}, // HeadLeft
			nil, // HeadDown
		},
		[]model.TwoArms{{
			Two: model.Arm{
				Heading: model.HeadLeft,
				Len:     3,
			},
			One: model.Arm{
				Heading: model.HeadRight,
				Len:     3,
			},
		}},
		model.NewNode(model.NewCoord(25, 3), true, 6),
		model.HeadLeft,
		1,
	)

	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 0)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 1)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 2)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 3)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 4)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 5)))

	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 6)))
	assert.True(t, puzz.edges.IsAvoided(model.NewEdgePair(model.NewCoord(25, 6), model.HeadRight)))

	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 7)))
	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 8)))
}

func TestBuildTwoArmsCache(t *testing.T) {
	// PuzzleID: 530,864
	fresh := NewPuzzle(5, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: false,
		Value:   4,
	}, {
		Row:     0,
		Col:     3,
		IsWhite: true,
		Value:   3,
	}, {
		Row:     1,
		Col:     4,
		IsWhite: false,
		Value:   5,
	}, {
		Row:     2,
		Col:     0,
		IsWhite: true,
		Value:   4,
	}, {
		Row:     2,
		Col:     5,
		IsWhite: true,
		Value:   5,
	}, {
		Row:     4,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}, {
		Row:     5,
		Col:     2,
		IsWhite: false,
		Value:   2,
	}})

	outPuzz, ms := claimGimmes(fresh)
	require.Equal(t, model.Incomplete, ms)

	updateCache(&outPuzz)

	expOptionsByNode := map[model.Node][]model.TwoArms{
		model.NewNode(model.NewCoord(0, 1), false, 4): {{
			One: model.Arm{
				Heading: model.HeadLeft,
				Len:     1,
			},
			Two: model.Arm{
				Heading: model.HeadDown,
				Len:     3,
			},
		}},
	}
	for i, tao := range outPuzz.twoArmOptions {
		assert.NotEmpty(t, tao)
		if expOptions, ok := expOptionsByNode[outPuzz.nodes[i]]; ok {
			assert.Equal(t, expOptions, tao)
		}
	}

	/*
		(   )---(b 4) X (   )---(w 3)---(   )---(   )
		  |       |       |       X       X       |
		(   ) X (   )   (   )   (   )---(b 5) X (   )
		  |                               |       |
		(w 4) X (   )   (   )   (   )   (   ) X (w 5)
		  |                                       |
		(   ) X (   )   (   )   (   )   (   ) X (   )
		  |                                       |
		(   )---(b 2) X (   )   (   )   (   ) X (   )
		  X               |                       |
		(   ) X (   )   (b 2)   (   )   (   )---(   )

	*/
}

func TestGimmesPuzzle5817105(t *testing.T) {
	// like puzzle 5,817,105
	fresh := NewPuzzle(25, []model.NodeLocation{{
		Row:     24,
		Col:     0,
		IsWhite: true,
		Value:   7,
	}, {
		Row:     24,
		Col:     1,
		IsWhite: false,
		Value:   4,
	}, {
		Row:     25,
		Col:     3,
		IsWhite: true,
		Value:   6,
	}, {
		Row:     24,
		Col:     5,
		IsWhite: false,
		Value:   3,
	}, {
		Row:     23,
		Col:     5,
		IsWhite: false,
		Value:   3,
	}, {
		Row:     25,
		Col:     8,
		IsWhite: false,
		Value:   2,
	}})

	puzz, s := ClaimGimmes(fresh)
	require.Equal(t, model.Incomplete, s)
	t.Logf("puzz: \n%s\n", puzz)

	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 0)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 1)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 2)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 3)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 4)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 5)))
	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 6)))

	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(25, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(24, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(23, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(22, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(21, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(20, 0)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(19, 0)))
	assert.False(t, puzz.IsEdge(model.HeadUp, model.NewCoord(18, 0)))

	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(25, 6)))

	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(24, 5)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(24, 5)))
	assert.True(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(24, 4)))
	assert.False(t, puzz.IsEdge(model.HeadLeft, model.NewCoord(24, 3)))

	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(24, 1)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(24, 1)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(23, 1)))
	assert.True(t, puzz.IsEdge(model.HeadUp, model.NewCoord(22, 1)))
}
