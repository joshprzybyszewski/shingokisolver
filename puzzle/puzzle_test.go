package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildGrid(
	t *testing.T,
	g *puzzle,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) *puzzle {
	var err error
	c := startCoord
	for _, s := range steps {
		c, g, err = g.AddEdge(s, c)
		require.NoError(t, err, `bad grid! %+v`, g)
	}
	t.Logf("buildGrid produced: \n%s\n", g)
	return g
}

func TestNumEdgesAndNodes(t *testing.T) {
	g := NewPuzzle(5, nil)
	assert.Equal(t, 6, g.numNodes())

	g = NewPuzzle(3, nil)
	assert.Equal(t, 4, g.numNodes())

	g = NewPuzzle(25, nil)
	assert.Equal(t, 26, g.numNodes())
}

func TestIsRangeInvalidWithBoundsCheck(t *testing.T) {
	g := NewPuzzle(5, nil)

	assert.False(t, g.IsRangeInvalid(-1, model.MAX_EDGES+1, -55, model.MAX_EDGES+1))
	assert.True(t, g.isRangeInvalid(-1, model.MAX_EDGES+1, -55, model.MAX_EDGES+1))
}

func TestIsEdge(t *testing.T) {
	g := NewPuzzle(5, nil)

	g = buildGrid(t, g,
		model.NodeCoord{},
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
		model.HeadDown,
	)

	assert.True(t, g.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		0,
	)))
	assert.True(t, g.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		1,
	)))
	assert.True(t, g.IsEdge(model.HeadRight, model.NewCoordFromInts(
		0,
		2,
	)))
	assert.True(t, g.IsEdge(model.HeadDown, model.NewCoordFromInts(
		0,
		3,
	)))
}

// defined here as a helper to make unit tests easier
func (g *puzzle) isInvalid() bool {
	return g.isRangeInvalid(0, model.RowIndex(g.numNodes()), 0, model.ColIndex(g.numNodes()))
}

func TestIsInvalid(t *testing.T) {
	g := NewPuzzle(2, nil)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
	)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
	)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadDown,
	)
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	g := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			1,
			0,
		),
		model.HeadRight,
		model.HeadRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	g := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			1,
			0,
		),
		model.HeadRight,
		model.HeadDown,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	g := NewPuzzle(4, []model.NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadDown,
	)
	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			1,
		),
		model.HeadRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	g := NewPuzzle(4, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	g := NewPuzzle(5, []model.NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	g = buildGrid(t, g,
		model.NewCoordFromInts(
			0,
			0,
		),
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
	)

	assert.True(t, g.isInvalid())
}

func TestGetEdgesFromNode(t *testing.T) {
	g := NewPuzzle(3, []model.NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		model.NewCoordFromInts(0, 0),
		model.HeadRight,
		model.HeadRight,
		model.HeadDown,
		model.HeadDown,
	)

	assert.True(t, g.IsEdge(model.HeadUp, model.NewCoordFromInts(1, 2)))
	assert.False(t, g.IsEdge(model.HeadUp, model.NewCoordFromInts(0, 2)))

	assert.True(t, g.IsEdge(model.HeadLeft, model.NewCoordFromInts(0, 1)))
	assert.False(t, g.IsEdge(model.HeadLeft, model.NewCoordFromInts(0, 0)))

	assert.False(t, g.IsEdge(model.HeadRight, model.NewCoordFromInts(1, 1)))
	assert.False(t, g.IsEdge(model.HeadDown, model.NewCoordFromInts(1, 1)))

	efn, ok := g.getOutgoingEdgesFrom(model.NewCoordFromInts(1, 1))
	assert.True(t, ok)

	expEFN := model.OutgoingEdges{
		// TODO
		// above: 1,
		// left:  1,
	}
	assert.Equal(t, expEFN, efn)
}
