package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
)

func TestBuildTargets(t *testing.T) {
	pd, err := reader.FromString(`........
.b2.w3....
.b6....w2.
w4w2...w4.b2
..w2.....
..b3.....
........
..b4....b4`)
	require.NoError(t, err)

	puzz := puzzle.NewPuzzle(pd.NumEdges, pd.Nodes)
	expTargets := []*model.Target{{
		Coord: model.NewCoordFromInts(3, 7),
		Node:  model.NewNode(false, 2),
	}, {
		Coord: model.NewCoordFromInts(1, 1),
		Node:  model.NewNode(false, 2),
	}, {
		Coord: model.NewCoordFromInts(2, 6),
		Node:  model.NewNode(false, 2),
	}, {
		Coord: model.NewCoordFromInts(3, 1),
		Node:  model.NewNode(false, 2),
	}, {
		Coord: model.NewCoordFromInts(4, 2),
		Node:  model.NewNode(false, 2),
	}, {
		Coord: model.NewCoordFromInts(1, 3),
		Node:  model.NewNode(false, 3),
	}, {
		Coord: model.NewCoordFromInts(5, 2),
		Node:  model.NewNode(false, 3),
	}, {
		Coord: model.NewCoordFromInts(7, 7),
		Node:  model.NewNode(false, 4),
	}, {
		Coord: model.NewCoordFromInts(3, 0),
		Node:  model.NewNode(false, 4),
	}, {
		Coord: model.NewCoordFromInts(7, 2),
		Node:  model.NewNode(false, 4),
	}, {
		Coord: model.NewCoordFromInts(3, 5),
		Node:  model.NewNode(false, 4),
	}, {
		Coord: model.NewCoordFromInts(2, 1),
		Node:  model.NewNode(false, 6),
	}}

	actTargets := puzz.Targets()
	assert.Len(t, actTargets, len(expTargets))
	for i, actTarget := range actTargets {
		require.NotNil(t, actTarget)
		if i == len(actTargets)-1 {
			assert.Nil(t, actTarget.Next, `final target should have a nil next target`)
		} else {
			assert.Equal(t, actTarget.Next, actTargets[i+1], `target at index %d did not point to the next target`, i)
		}
		assert.Equal(t, expTargets[i].Coord, actTarget.Coord, `target at index %d had unexpected coords`, i)
		assert.Equal(t, expTargets[i].Node, actTarget.Node, `target at index %d had unexpected coords`, i)
	}
}
