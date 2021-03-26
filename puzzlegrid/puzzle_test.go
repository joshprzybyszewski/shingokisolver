package puzzlegrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdges(t *testing.T) {
	e := newGridEdges()
	for i := 0; i < MAX_EDGES; i++ {
		assert.False(t, e.isEdge(i))
	}
	t.Logf("e: %b\n", e)

	v2 := 5
	e2 := e.addEdge(v2)
	t.Logf("e2: %b\n", e2)
	assert.True(t, e2.isEdge(v2))
	assert.False(t, e.isEdge(v2))

	v3 := 13
	e3 := e2.addEdge(v3)
	t.Logf("e3: %b\n", e3)
	assert.True(t, e3.isEdge(v3))
	assert.False(t, e2.isEdge(v3))
	assert.False(t, e.isEdge(v3))

	v4 := 0
	e4 := e3.addEdge(v4)
	t.Logf("e4: %b\n", e4)
	assert.True(t, e4.isEdge(v4))
	assert.False(t, e3.isEdge(v4))
	assert.False(t, e2.isEdge(v4))
	assert.False(t, e.isEdge(v4))

	v5 := -3
	e5 := e4.addEdge(v5)
	t.Logf("e5: %b\n", e5)
	assert.False(t, e5.isEdge(v5))
	assert.False(t, e4.isEdge(v5))
	assert.False(t, e3.isEdge(v5))
	assert.False(t, e2.isEdge(v5))
	assert.False(t, e.isEdge(v5))

	v6 := MAX_EDGES + 1
	e6 := e5.addEdge(v6)
	t.Logf("e6: %b\n", e6)
	assert.False(t, e6.isEdge(v6))
	assert.False(t, e5.isEdge(v6))
	assert.False(t, e4.isEdge(v6))
	assert.False(t, e3.isEdge(v6))
	assert.False(t, e2.isEdge(v6))
	assert.False(t, e.isEdge(v6))
}

func TestIsEdge(t *testing.T) {
	g := newGrid(5, nil)
	t.Logf("g: \n%s\n", g.String())

	g, err := g.AddEdge(rowDir, 0, 0)
	t.Logf("g: \n%s\n", g.String())
	require.NoError(t, err)
	assert.True(t, g.IsEdge(rowDir, 0, 0))

	g, err = g.AddEdge(rowDir, 0, 1)
	require.NoError(t, err)
	g, err = g.AddEdge(rowDir, 0, 2)
	require.NoError(t, err)
	g, err = g.AddEdge(colDir, 2, 0)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsEdge(rowDir, 0, 0))
	assert.True(t, g.IsEdge(rowDir, 0, 1))
	assert.True(t, g.IsEdge(rowDir, 0, 2))
	assert.True(t, g.IsEdge(colDir, 2, 0))
}

func TestIsInvalid(t *testing.T) {
	g := newGrid(3, nil)
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(colDir, 1, 0)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsInvalid())
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 1, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 1, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsInvalid())
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 1, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(colDir, 1, 0)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsInvalid())
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(colDir, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 1)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 2)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.True(t, g.IsInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	g := newGrid(5, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	assert.False(t, g.IsInvalid())

	g, err := g.AddEdge(rowDir, 0, 0)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 1)
	require.NoError(t, err)
	assert.False(t, g.IsInvalid())

	g, err = g.AddEdge(rowDir, 0, 3)
	require.NoError(t, err)
	t.Logf("g: \n%s\n", g.String())
	assert.False(t, g.IsInvalid())
}
