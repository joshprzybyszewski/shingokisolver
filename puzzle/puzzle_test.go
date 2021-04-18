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

	s := p.AddEdge(
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadDown,
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

	s := p.AddEdge(
		model.NewCoordFromInts(
			1,
			1,
		),
		model.HeadRight,
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

	s := p.AddEdge(
		model.NewCoordFromInts(
			1,
			1,
		),
		model.HeadDown,
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
	s := p.AddEdge(
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
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
	s := p.ClaimGimmes()
	require.Equal(t, model.Incomplete, s)

	s = p.AddEdge(
		model.NewCoordFromInts(
			0,
			2,
		),
		model.HeadRight,
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

	nOut, isMax := p.GetSumOutgoingStraightLines(model.NewCoordFromInts(1, 1))
	assert.Equal(t, int8(2), nOut)
	assert.False(t, isMax)
}
