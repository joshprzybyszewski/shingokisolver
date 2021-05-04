package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumEdgesAndNodes(t *testing.T) {
	p := NewPuzzle(5, nil)
	assert.Equal(t, 6, p.numNodes())

	p = NewPuzzle(3, nil)
	assert.Equal(t, 4, p.numNodes())

	p = NewPuzzle(25, nil)
	assert.Equal(t, 26, p.numNodes())
}

func TestIsEdge(t *testing.T) {
	p := NewPuzzle(5, nil)

	p = BuildTestPuzzle(t, p,
		model.NodeCoord{},
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
		model.HeadDown,
	)

	assert.True(t, p.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		0,
	)))
	assert.True(t, p.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		1,
	)))
	assert.True(t, p.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		2,
	)))
	assert.True(t, p.IsEdge(model.HeadDown, model.NewCoordFromInts(
		0,
		3,
	)))
}

func TestIsInvalid(t *testing.T) {
	p := NewPuzzle(2, nil)

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
	)
	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
	)

	_, s := AddEdge(p,
		model.NewEdgePair(
			model.NewCoordFromInts(
				0,
				1,
			),
			model.HeadDown,
		),
	)
	assert.Equal(t, model.Violation, s)
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	p := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			1,
			0,
		),
		model.HeadRight,
	)

	_, s := AddEdge(p,
		model.NewEdgePair(
			model.NewCoordFromInts(
				1,
				1,
			),
			model.HeadRight,
		),
	)
	assert.Equal(t, model.Violation, s)
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	p := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			1,
			0,
		),
		model.HeadRight,
	)

	_, s := AddEdge(p,
		model.NewEdgePair(
			model.NewCoordFromInts(
				1,
				1,
			),
			model.HeadDown,
		),
	)
	assert.Equal(t, model.Violation, s)
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	p := NewPuzzle(4, []model.NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadDown,
	)
	_, s := AddEdge(p,
		model.NewEdgePair(
			model.NewCoordFromInts(
				0,
				1,
			),
			model.HeadRight,
		),
	)
	assert.Equal(t, model.Violation, s)
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	p := NewPuzzle(4, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	p, s := ClaimGimmes(p)
	require.Equal(t, model.Incomplete, s)

	_, s = AddEdge(p,
		model.NewEdgePair(
			model.NewCoordFromInts(
				0,
				2,
			),
			model.HeadRight,
		),
	)
	assert.Equal(t, model.Violation, s)
}

func TestGetEdgesFromNode(t *testing.T) {
	p := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(0, 0),
		model.HeadRight,
		model.HeadRight,
		model.HeadDown,
		model.HeadDown,
	)

	assert.True(t, p.IsEdge(model.HeadUp, model.NewCoordFromInts(1, 2)))
	assert.False(t, p.IsEdge(model.HeadUp, model.NewCoordFromInts(0, 2)))

	assert.True(t, p.IsEdge(model.HeadLeft, model.NewCoordFromInts(0, 1)))
	assert.False(t, p.IsEdge(model.HeadLeft, model.NewCoordFromInts(0, 0)))

	assert.False(t, p.IsEdge(model.HeadRight, model.NewCoordFromInts(1, 1)))
	assert.True(t, p.IsEdge(model.HeadDown, model.NewCoordFromInts(1, 1)))

	nOut, isMax := getSumOutgoingStraightLines(model.NewCoordFromInts(1, 1), &p.edges)
	assert.Equal(t, int8(2), nOut)
	assert.False(t, isMax)
}

func TestGetFilteredOptionsDoubleAvoid(t *testing.T) {
	puzz := NewPuzzle(10, []model.NodeLocation{{
		Row:     10,
		Col:     2,
		Value:   6,
		IsWhite: false,
	}})
	puzz, ms := ClaimGimmes(puzz)
	require.Equal(t, model.Incomplete, ms)

	puzz, ms = AvoidEdge(puzz, model.NewEdgePair(model.NewCoord(10, 0), model.HeadRight))
	require.Equal(t, model.Incomplete, ms)
	t.Logf("puzz: \n%s\n", puzz)

	ft, ms := puzz.GetFirstTarget()
	require.Equal(t, model.Incomplete, ms)
	assert.Equal(t, model.NewNode(model.NewCoord(10, 2), false, 6), ft.Node)

	require.NotEmpty(t, ft.Options)
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 1}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 2}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 3}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 4}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 5}))
	optionsBefore := ft.Options

	// Now avoid edges on either side of "HU2"
	puzz2, ms := AvoidEdge(puzz, model.NewEdgePair(model.NewCoord(8, 1), model.HeadRight))
	require.Equal(t, model.Incomplete, ms)
	puzz2, ms = AvoidEdge(puzz2, model.NewEdgePair(model.NewCoord(8, 2), model.HeadRight))
	require.Equal(t, model.Incomplete, ms)
	t.Logf("puzz2: \n%s\n", puzz2)

	ft, ms = puzz2.GetFirstTarget()
	require.Equal(t, model.Incomplete, ms)
	assert.Equal(t, model.NewNode(model.NewCoord(10, 2), false, 6), ft.Node)

	require.NotEmpty(t, ft.Options)
	assert.NotEqual(t, optionsBefore, ft.Options)
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 1}))
	assert.False(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 2}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 3}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 4}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 5}))

	// Now avoid an edge on one side of "HU2" and add the other
	puzz3, ms := AddEdge(puzz, model.NewEdgePair(model.NewCoord(8, 1), model.HeadRight))
	require.Equal(t, model.Incomplete, ms)
	puzz3, ms = AvoidEdge(puzz3, model.NewEdgePair(model.NewCoord(8, 2), model.HeadRight))
	require.Equal(t, model.Incomplete, ms)
	t.Logf("puzz3: \n%s\n", puzz3)

	ft, ms = puzz3.GetFirstTarget()
	require.Equal(t, model.Incomplete, ms)
	assert.Equal(t, model.NewNode(model.NewCoord(10, 2), false, 6), ft.Node)

	require.NotEmpty(t, ft.Options)
	assert.NotEqual(t, optionsBefore, ft.Options)
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 1}))
	assert.True(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 2}))
	assert.False(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 3}))
	assert.False(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 4}))
	assert.False(t, containsArm(ft.Options, model.Arm{Heading: model.HeadUp, Len: 5}))
}

func containsArm(options []model.TwoArms, arm model.Arm) bool {
	for _, ta := range options {
		if ta.One == arm || ta.Two == arm {
			return true
		}
	}
	return false
}

func Test5817105Example(t *testing.T) {
	puzz := NewPuzzle(25, []model.NodeLocation{{
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
	t.Logf("puzz:\n%s", puzz)
	puzz, ms := ClaimGimmes(puzz)
	require.Equal(t, model.Incomplete, ms)
	t.Logf("puzz:\n%s", puzz)

	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 0)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 1)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 2)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 3)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 4)))
	assert.True(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 5)))
	assert.False(t, puzz.IsEdge(model.HeadRight, model.NewCoord(25, 6)))
}
