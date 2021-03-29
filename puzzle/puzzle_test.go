package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildGrid(
	t *testing.T,
	g *puzzle,
	startCoord nodeCoord,
	steps ...cardinal,
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
	g := newPuzzle(5, nil)
	assert.Equal(t, 6, g.numNodes())

	g = newPuzzle(3, nil)
	assert.Equal(t, 4, g.numNodes())

	g = newPuzzle(25, nil)
	assert.Equal(t, 26, g.numNodes())
}

func TestIsRangeInvalidWithBoundsCheck(t *testing.T) {
	g := newPuzzle(5, nil)

	assert.False(t, g.isRangeInvalidWithBoundsCheck(-1, MAX_EDGES+1, -55, MAX_EDGES+1))
	assert.True(t, g.isRangeInvalid(-1, MAX_EDGES+1, -55, MAX_EDGES+1))
}

func TestIsEdge(t *testing.T) {
	g := newPuzzle(5, nil)

	g = buildGrid(t, g,
		nodeCoord{},
		headRight,
		headRight,
		headRight,
		headDown,
	)

	assert.True(t, g.IsEdge(headRight, nodeCoord{
		row: 0,
		col: 0,
	}))
	assert.True(t, g.IsEdge(headRight, nodeCoord{
		row: 0,
		col: 1,
	}))
	assert.True(t, g.IsEdge(headRight, nodeCoord{
		row: 0,
		col: 2,
	}))
	assert.True(t, g.IsEdge(headDown, nodeCoord{
		row: 0,
		col: 3,
	}))
}

// defined here as a helper to make unit tests easier
func (g *puzzle) isInvalid() bool {
	return g.isRangeInvalid(0, rowIndex(g.numNodes()), 0, colIndex(g.numNodes()))
}

func TestIsInvalid(t *testing.T) {
	g := newPuzzle(2, nil)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 0,
		},
		headRight,
	)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 1,
		},
		headRight,
	)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 1,
		},
		headDown,
	)
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	g := newPuzzle(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		nodeCoord{
			row: 1,
			col: 0,
		},
		headRight,
		headRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	g := newPuzzle(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g,
		nodeCoord{
			row: 1,
			col: 0,
		},
		headRight,
		headDown,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	g := newPuzzle(4, []NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 0,
		},
		headRight,
		headDown,
	)
	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 1,
		},
		headRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	g := newPuzzle(4, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 0,
		},
		headRight,
		headRight,
		headRight,
	)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	g := newPuzzle(5, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 0,
		},
		headRight,
		headRight,
		headRight,
	)

	assert.True(t, g.isInvalid())
}

func TestGetEdgesFromNode(t *testing.T) {
	g := newPuzzle(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g,
		nodeCoord{
			row: 0,
			col: 0,
		},
		headRight,
		headRight,
		headDown,
		headDown,
	)

	assert.True(t, g.IsEdge(headUp, nodeCoord{
		row: 1,
		col: 2,
	}))
	assert.False(t, g.IsEdge(headUp, nodeCoord{
		row: 0,
		col: 2,
	}))

	assert.True(t, g.IsEdge(headLeft, nodeCoord{
		row: 0,
		col: 1,
	}))
	assert.False(t, g.IsEdge(headLeft, nodeCoord{
		row: 0,
		col: 0,
	}))

	assert.False(t, g.IsEdge(headRight, nodeCoord{
		row: 1,
		col: 1,
	}))
	assert.False(t, g.IsEdge(headDown, nodeCoord{
		row: 1,
		col: 1,
	}))

	efn, ok := g.getEdgesFromNode(nodeCoord{
		row: 1,
		col: 1,
	})
	assert.True(t, ok)

	expEFN := edgesFromNode{
		above: 1,
		left:  1,
	}
	assert.Equal(t, expEFN, efn)
}
