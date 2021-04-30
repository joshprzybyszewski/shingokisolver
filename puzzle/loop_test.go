package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLooperNumNodesInLoop(t *testing.T) {
	puzz := BuildTestPuzzle(
		t,
		NewPuzzle(5, []model.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: true,
			Value:   2,
		}}),
		model.NewCoord(1, 1),
		model.HeadDown,
	)
	SetNodesComplete(&puzz)

	assert.Zero(t, puzz.loop.NumNodesInLoop())

	puzz, state := AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)
	assert.Equal(t, 1, puzz.loop.NumNodesInLoop())
}

func TestLooperIsLoop(t *testing.T) {
	puzz := BuildTestPuzzle(
		t,
		NewPuzzle(5, []model.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: true,
			Value:   2,
		}}),
		model.NewCoord(1, 1),
		model.HeadDown,
	)
	SetNodesComplete(&puzz)

	assert.False(t, puzz.loop.IsLoop())

	puzz, state := AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)
	assert.True(t, puzz.loop.IsLoop())
}

func TestLooperGetUnknownEdge(t *testing.T) {
	puzz := BuildTestPuzzle(
		t,
		NewPuzzle(5, []model.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: true,
			Value:   2,
		}}),
		model.NewCoord(1, 1),
		model.HeadDown,
	)
	SetNodesComplete(&puzz)

	assert.False(t, puzz.loop.IsLoop())

	potentialUnknowns := []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadRight),
		model.NewEdgePair(model.NewCoord(0, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(0, 1), model.HeadRight),
	}

	ep, state := puzz.loop.GetUnknownEdge(&puzz.edges)
	require.Equal(t, model.Incomplete, state)
	assert.Contains(t, potentialUnknowns, ep)

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)

	ep, state = puzz.loop.GetUnknownEdge(&puzz.edges)
	assert.Equal(t, model.Violation, state)
	assert.Equal(t, model.InvalidEdgePair, ep)
}

func TestLooperWithUpdatedEdges(t *testing.T) {
	// TODO
}
