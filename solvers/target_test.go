package solvers

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
	targets := buildTargets(puzz)
	expTargets := []*target{{
		coord: model.NewCoordFromInts(2, 1),
		val:   6,
	}, {
		coord: model.NewCoordFromInts(7, 7),
		val:   4,
	}, {
		coord: model.NewCoordFromInts(3, 0),
		val:   4,
	}, {
		coord: model.NewCoordFromInts(7, 2),
		val:   4,
	}, {
		coord: model.NewCoordFromInts(3, 5),
		val:   4,
	}, {
		coord: model.NewCoordFromInts(1, 3),
		val:   3,
	}, {
		coord: model.NewCoordFromInts(5, 2),
		val:   3,
	}, {
		coord: model.NewCoordFromInts(3, 7),
		val:   2,
	}, {
		coord: model.NewCoordFromInts(1, 1),
		val:   2,
	}, {
		coord: model.NewCoordFromInts(2, 6),
		val:   2,
	}, {
		coord: model.NewCoordFromInts(3, 1),
		val:   2,
	}, {
		coord: model.NewCoordFromInts(4, 2),
		val:   2,
	}}

	assert.Len(t, targets, len(expTargets))
	for i, actTarget := range targets {
		require.NotNil(t, actTarget)
		if i == len(targets)-1 {
			assert.Nil(t, actTarget.next, `final target should have a nil next target`)
		} else {
			assert.Equal(t, actTarget.next, targets[i+1], `target at index %d did not point to the next target`, i)
		}
		assert.Equal(t, expTargets[i].coord, actTarget.coord, `target at index %d had unexpected coords`, i)
		assert.Equal(t, expTargets[i].val, actTarget.val, `target at index %d had unexpected coords`, i)
	}
}
