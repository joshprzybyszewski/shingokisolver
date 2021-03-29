package solvers

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPartialSolutionsForNode(t *testing.T) {
	pd, err := reader.FromString(`........
.b2.w3....
.b6....w2.
w4w2...w4.b2
..w2.....
..b3.....
........
..b4....b4`)
	require.NoError(t, err)

	s := newTargetSolver(pd.NumEdges, pd.Nodes)
	ts, ok := s.(*targetSolver)
	require.True(t, ok)
	puzz := ts.puzzle.DeepCopy()

	actPartials := ts.getPartialSolutionsForNode(
		&partialSolutionItem{
			puzzle: puzz.DeepCopy(),
			targeting: &target{
				coord: model.NewCoordFromInts(4, 2),
				val:   2,
				next:  nil,
			},
			looseEnds: nil,
		},
	)
	assert.NotEmpty(t, actPartials)
	assert.Empty(t, actPartials)
	assert.Len(t, actPartials, 4)
	var actPuzzles []string
	for i, p := range actPartials {
		actStr := p.puzzle.String()
		actPuzzles = append(actPuzzles, actStr)
		t.Logf("p[%d].puzzle =\n%s\n", i, actStr)
	}

	var expPuzzles []string
	for _, expPuzz := range []*puzzle.Puzzle{
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(4, 2),
			model.HeadRight,
			model.HeadRight,
		),
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(4, 1),
			model.HeadRight,
			model.HeadRight,
		),
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(4, 0),
			model.HeadRight,
			model.HeadRight,
		),
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(4, 2),
			model.HeadUp,
			model.HeadUp,
		),
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(5, 2),
			model.HeadUp,
			model.HeadUp,
		),
	} {
		expPuzzles = append(expPuzzles, expPuzz.String())
	}
	assert.Equal(t, expPuzzles, actPuzzles)
}

func TestGetSolutionsForNodeInDirections(t *testing.T) {
	pd, err := reader.FromString(`........
.b2.w3....
.b6....w2.
w4w2...w4.b2
..w2.....
..b3.....
........
..b4....b4`)
	require.NoError(t, err)

	s := newTargetSolver(pd.NumEdges, pd.Nodes)
	ts, ok := s.(*targetSolver)
	require.True(t, ok)

	actPartials := ts.getSolutionsForNodeInDirections(
		nil,
		model.NewCoordFromInts(4, 2),
		model.Node{},
		model.HeadRight, model.HeadLeft,
	)
	assert.Empty(t, actPartials)
}
