package puzzlegrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolve3x3(t *testing.T) {
	numEdges := 2
	nodes := []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}}

	s := NewSolver(numEdges, nodes)
	err := s.Solve()
	require.NoError(t, err)
	assert.Error(t, err)
}
