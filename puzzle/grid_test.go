package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEdge(t *testing.T) {
	g := newGrid(5, nil)
	t.Logf("g: \n%s\n", g.String())

	g, err := g.AddEdge(headRight, 0, 0)
	t.Logf("g: \n%s\n", g.String())
	require.NoError(t, err)
	assert.True(t, g.IsEdge(headRight, 0, 0))

	g, err = g.AddEdge(headRight, 0, 1)
	require.NoError(t, err)
	g, err = g.AddEdge(headRight, 0, 2)
	require.NoError(t, err)
	g, err = g.AddEdge(headDown, 0, 2)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsEdge(headRight, 0, 0))
	assert.True(t, g.IsEdge(headRight, 0, 1))
	assert.True(t, g.IsEdge(headRight, 0, 2))
	assert.True(t, g.IsEdge(headDown, 0, 2))
}

func TestIsInvalid(t *testing.T) {
	g := newGrid(3, nil)
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headDown, 0, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 1, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 1, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 1, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headDown, 0, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headDown, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 2)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.isInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	g := newGrid(5, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.isInvalid())

	g, err := g.AddEdge(headRight, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.isInvalid())

	g, err = g.AddEdge(headRight, 0, 3)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.False(t, g.isInvalid())
}

func TestGetEdgesFromNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g, err := g.AddEdge(headRight, 0, 0)
	require.NoError(t, err)

	g, err = g.AddEdge(headRight, 1, 0)
	require.NoError(t, err)
	g, err = g.AddEdge(headDown, 0, 0)
	require.NoError(t, err)
	g, err = g.AddEdge(headDown, 0, 1)
	require.NoError(t, err)

	assert.True(t, g.IsEdge(headUp, 1, 1))
	assert.False(t, g.IsEdge(headUp, 0, 1))

	assert.True(t, g.IsEdge(headLeft, 1, 1))
	assert.False(t, g.IsEdge(headLeft, 1, 0))

	assert.False(t, g.IsEdge(headRight, 1, 1))
	assert.False(t, g.IsEdge(headDown, 1, 1))

	efn := g.getEdgesFromNode(1, 1)
	require.NotNil(t, efn)
	expEFN := &edgesFromNode{
		above:      1,
		left:       1,
		totalEdges: 2,
		isabove:    true,
		isleft:     true,
	}
	assert.Equal(t, expEFN, efn)
}
