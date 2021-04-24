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
		map[model.Cardinal][]*model.Node{
			model.HeadRight: []*model.Node{
				nil,
			},
			model.HeadLeft: []*model.Node{
				nil,
			},
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

	twoArmOptions := buildTwoArmsCache(
		outPuzz.nodes,
		outPuzz.numEdges(),
		&outPuzz.edges,
	)

	expOptionsByNode := map[model.Node][]model.TwoArms{
		model.NewNode(model.NewCoord(0, 1), false, 4): []model.TwoArms{{
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
	for _, tao := range twoArmOptions {
		assert.NotEmpty(t, tao.Options)
		if expOptions, ok := expOptionsByNode[tao.Node]; ok {
			assert.Equal(t, expOptions, tao.Options)
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
