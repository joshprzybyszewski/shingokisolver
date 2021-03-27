package puzzle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve3x3(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	numEdges := 2
	nodes := []NodeLocation{{
		Row:     1,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}}

	for _, st := range AllSolvers {
		s := NewSolver(numEdges, nodes, st)
		err := s.Solve()
		require.NoError(t, err)
	}
}

func TestSolve5x5(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	numEdges := 5
	nodes := []NodeLocation{{
		Row:     3,
		Col:     2,
		IsWhite: false,
		Value:   4,
	}, {
		Row:     3,
		Col:     5,
		IsWhite: true,
		Value:   5,
	}, {
		Row:     4,
		Col:     0,
		IsWhite: true,
		Value:   5,
	}, {
		Row:     5,
		Col:     1,
		IsWhite: false,
		Value:   2,
	}, {
		Row:     5,
		Col:     3,
		IsWhite: false,
		Value:   5,
	}}

	for _, st := range AllSolvers {
		s := NewSolver(numEdges, nodes, st)
		err := s.Solve()
		require.NoError(t, err)
	}
}
