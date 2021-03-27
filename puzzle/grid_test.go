package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gridBuildStep struct {
	move cardinal
	row  int
	col  int
}

func buildGrid(
	t *testing.T,
	g *grid,
	steps ...gridBuildStep,
) *grid {
	var err error
	for _, s := range steps {
		g, err = g.AddEdge(s.move, s.row, s.col)
		require.NoError(t, err)
	}
	t.Logf("buildGrid produced: \n%s\n", g)
	return g
}

func TestNumEdgesAndNodes(t *testing.T) {
	g := newGrid(5, nil)
	assert.Equal(t, 5, g.numEdges())
	assert.Equal(t, 6, g.numNodes())

	g = newGrid(3, nil)
	assert.Equal(t, 3, g.numEdges())
	assert.Equal(t, 4, g.numNodes())

	g = newGrid(25, nil)
	assert.Equal(t, 25, g.numEdges())
	assert.Equal(t, 26, g.numNodes())
}

func TestIsRangeInvalidWithBoundsCheck(t *testing.T) {
	g := newGrid(5, nil)

	assert.False(t, g.isRangeInvalidWithBoundsCheck(-1, 10000, -55, 10000))
	assert.True(t, g.isRangeInvalid(-1, 10000, -55, 10000))
}

func TestIsEdge(t *testing.T) {
	g := newGrid(5, nil)

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  0,
		col:  0,
	}, {
		move: headRight,
		row:  0,
		col:  1,
	}, {
		move: headRight,
		row:  0,
		col:  2,
	}, {
		move: headDown,
		row:  0,
		col:  2,
	}}...)

	assert.True(t, g.IsEdge(headRight, 0, 0))
	assert.True(t, g.IsEdge(headRight, 0, 1))
	assert.True(t, g.IsEdge(headRight, 0, 2))
	assert.True(t, g.IsEdge(headDown, 0, 2))
}

// defined here as a helper to make unit tests easier
func (g *grid) isInvalid() bool {
	return g.isRangeInvalid(0, len(g.rows), 0, len(g.cols))
}

func TestIsInvalid(t *testing.T) {
	g := newGrid(3, nil)
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g, gridBuildStep{
		move: headRight,
		row:  0,
		col:  0,
	})
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g, gridBuildStep{
		move: headRight,
		row:  0,
		col:  1,
	})
	assert.False(t, g.isInvalid())

	g = buildGrid(t, g, gridBuildStep{
		move: headDown,
		row:  0,
		col:  1,
	})
	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  1,
		col:  0,
	}, {
		move: headRight,
		row:  1,
		col:  1,
	}}...)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  1,
		col:  0,
	}, {
		move: headDown,
		row:  0,
		col:  1,
	}}...)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadBlackNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     0,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  0,
		col:  0,
	}, {
		move: headDown,
		row:  0,
		col:  0,
	}, {
		move: headRight,
		row:  0,
		col:  1,
	}}...)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidBadWhiteNodeTooManyLines(t *testing.T) {
	g := newGrid(4, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  0,
		col:  0,
	}, {
		move: headRight,
		row:  0,
		col:  1,
	}, {
		move: headRight,
		row:  0,
		col:  2,
	}}...)

	assert.True(t, g.isInvalid())
}

func TestIsInvalidGoodWhiteNodeAllowsTheRowToHaveManyEdges(t *testing.T) {
	g := newGrid(5, []NodeLocation{{
		Row:     0,
		Col:     1,
		IsWhite: true,
		Value:   2,
	}})
	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  0,
		col:  0,
	}, {
		move: headRight,
		row:  0,
		col:  1,
	}, {
		move: headRight,
		row:  0,
		col:  3,
	}}...)

	assert.False(t, g.isInvalid())
}

func TestGetEdgesFromNode(t *testing.T) {
	g := newGrid(3, []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}})

	g = buildGrid(t, g, []gridBuildStep{{
		move: headRight,
		row:  0,
		col:  0,
	}, {
		move: headRight,
		row:  1,
		col:  0,
	}, {
		move: headDown,
		row:  0,
		col:  0,
	}, {
		move: headDown,
		row:  0,
		col:  1,
	}}...)

	assert.True(t, g.IsEdge(headUp, 1, 1))
	assert.False(t, g.IsEdge(headUp, 0, 1))

	assert.True(t, g.IsEdge(headLeft, 1, 1))
	assert.False(t, g.IsEdge(headLeft, 1, 0))

	assert.False(t, g.IsEdge(headRight, 1, 1))
	assert.False(t, g.IsEdge(headDown, 1, 1))

	efn := g.getEdgesFromNode(1, 1)
	require.NotNil(t, efn)
	expEFN := edgesFromNode{
		above:       1,
		left:        1,
		totalEdges:  2,
		isabove:     true,
		isleft:      true,
		isPopulated: true,
	}
	assert.Equal(t, expEFN, efn)
}
