package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/shingokisolver/model"
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
	require.NotNil(t, puzz.loop)

	assert.Zero(t, puzz.loop.NumNodesInLoop())

	puzz, state := AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)
	require.NotNil(t, puzz.loop)
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
	require.NotNil(t, puzz.loop)

	assert.False(t, puzz.loop.IsLoop())

	puzz, state := AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)
	require.NotNil(t, puzz.loop)
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
	require.NotNil(t, puzz.loop)

	assert.False(t, puzz.loop.IsLoop())

	potentialUnknowns := []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadRight),
		model.NewEdgePair(model.NewCoord(0, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(0, 1), model.HeadRight),
	}

	require.NotNil(t, puzz.loop)
	ep, state := puzz.loop.GetUnknownEdge(&puzz.edges)
	require.Equal(t, model.Incomplete, state)
	assert.Contains(t, potentialUnknowns, ep)

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadUp),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)

	require.NotNil(t, puzz.loop)
	ep, state = puzz.loop.GetUnknownEdge(&puzz.edges)
	assert.Equal(t, model.Violation, state)
	assert.Equal(t, model.InvalidEdgePair, ep)
}

func TestLooperWithThreeNodes(t *testing.T) {
	puzz := BuildTestPuzzle(
		t,
		NewPuzzle(5, []model.NodeLocation{{
			Row:     0,
			Col:     1,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     5,
			Col:     4,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     2,
			Col:     0,
			IsWhite: false,
			Value:   3,
		}}),
		model.NewCoord(1, 5),
		model.HeadDown,
		model.HeadDown,
	)
	puzz, state := ClaimGimmes(puzz)
	require.Equal(t, model.Incomplete, state)

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 1), model.HeadRight),
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadDown),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Now puzzle is: \n%s\n", puzz)

	SetNodesComplete(&puzz)
	require.NotNil(t, puzz.loop)

	assert.False(t, puzz.loop.IsLoop())

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(1, 2), model.HeadDown),
		model.NewEdgePair(model.NewCoord(3, 1), model.HeadDown),
		model.NewEdgePair(model.NewCoord(3, 5), model.HeadDown),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Complete puzzle is: \n%s\n", puzz)
	require.NotNil(t, puzz.loop)
	assert.True(t, puzz.loop.IsLoop())
	assert.Equal(t, 3, puzz.loop.NumNodesInLoop())
}

func TestLooperWithFiveNodes(t *testing.T) {
	puzz := BuildTestPuzzle(
		t,
		NewPuzzle(5, []model.NodeLocation{{
			Row:     0,
			Col:     1,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     5,
			Col:     4,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     1,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     1,
			Col:     4,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     3,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     3,
			Col:     4,
			IsWhite: false,
			Value:   2,
		}}),
		model.NewCoord(1, 5),
		model.HeadDown,
	)
	puzz, state := ClaimGimmes(puzz)
	require.Equal(t, model.Incomplete, state)
	t.Logf("After ClaimGimmes: \n%s\n", puzz)

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(1, 1), model.HeadRight),
		model.NewEdgePair(model.NewCoord(3, 0), model.HeadRight),
		model.NewEdgePair(model.NewCoord(3, 1), model.HeadDown),
		model.NewEdgePair(model.NewCoord(3, 4), model.HeadLeft),
		model.NewEdgePair(model.NewCoord(3, 4), model.HeadDown),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("After adding edges: \n%s\n", puzz)

	SetNodesComplete(&puzz)
	require.NotNil(t, puzz.loop)

	assert.False(t, puzz.loop.IsLoop())

	puzz, state = AddEdges(puzz, []model.EdgePair{
		model.NewEdgePair(model.NewCoord(2, 0), model.HeadRight),
		model.NewEdgePair(model.NewCoord(2, 3), model.HeadDown),
		model.NewEdgePair(model.NewCoord(2, 4), model.HeadRight),
		model.NewEdgePair(model.NewCoord(4, 4), model.HeadRight),
	})
	require.Equal(t, model.Incomplete, state)
	t.Logf("Complete puzzle is: \n%s\n", puzz)
	require.NotNil(t, puzz.loop)
	assert.True(t, puzz.loop.IsLoop())
	assert.Equal(t, 6, puzz.loop.NumNodesInLoop())
}
