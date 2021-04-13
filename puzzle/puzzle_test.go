package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
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

// defined here as a helper to make unit tests easier
func (p *Puzzle) isInvalid() bool {
	switch p.GetState() {
	case model.Incomplete, model.Complete:
		return false
	default:
		return true
	}
}

func TestIsInvalid(t *testing.T) {
	p := NewPuzzle(2, nil)
	assert.False(t, p.isInvalid())

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
	)
	assert.False(t, p.isInvalid())

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
	)
	assert.False(t, p.isInvalid())

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadDown,
	)
	assert.True(t, p.isInvalid())
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
		model.HeadRight,
	)

	assert.True(t, p.isInvalid())
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
		model.HeadDown,
	)

	assert.True(t, p.isInvalid())
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
	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
	)

	assert.True(t, p.isInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	p := NewPuzzle(4, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
	)

	assert.True(t, p.isInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	p := NewPuzzle(5, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	p = BuildTestPuzzle(t, p,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
	)

	assert.True(t, p.isInvalid())
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
	assert.False(t, p.IsEdge(model.HeadDown, model.NewCoordFromInts(1, 1)))

	assert.True(t, p.HasTwoOutgoingEdges(model.NewCoordFromInts(1, 1)))
	assert.Equal(t, int8(2), p.GetSumOutgoingStraightLines(model.NewCoordFromInts(1, 1)))
}
