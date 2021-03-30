package solvers

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	normal7x7 = `........
.b2.w3....
.b6....w2.
w4w2...w4.b2
..w2.....
..b3.....
........
..b4....b4`
)

func TestGetPartialSolutionsForNode(t *testing.T) {
	pd, err := reader.FromString(normal7x7)
	require.NoError(t, err)

	s := newTargetSolver(pd.NumEdges, pd.Nodes)
	ts, ok := s.(*targetSolver)
	require.True(t, ok)
	puzz := ts.puzzle.DeepCopy()

	actPartials := ts.getAllPartialSolutionsForItem(
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
	require.NotEmpty(t, actPartials)
	assert.Len(t, actPartials, 2)
	var actPuzzles []string
	for i, p := range actPartials {
		actStr := p.puzzle.String()
		actPuzzles = append(actPuzzles, actStr)
		t.Logf("p[%d].puzzle =\n%s\n", i, actStr)
	}

	var expPuzzles []string
	for _, expPuzz := range []*puzzle.Puzzle{
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(4, 1),
			model.HeadRight,
			model.HeadRight,
		),
		puzzle.BuildTestPuzzle(t, puzz,
			model.NewCoordFromInts(5, 2),
			model.HeadUp,
			model.HeadUp,
		),
	} {
		expPuzzles = append(expPuzzles, expPuzz.String())
	}
	assert.ElementsMatch(t, expPuzzles, actPuzzles)
}

func TestTargetSolverBuildAllPartials(t *testing.T) {
	pd, err := reader.FromString(normal7x7)
	require.NoError(t, err)

	s := newTargetSolver(pd.NumEdges, pd.Nodes)
	ts, ok := s.(*targetSolver)
	require.True(t, ok)

	puzz := puzzle.NewPuzzle(pd.NumEdges, pd.Nodes)
	targets := buildTargets(puzz)

	ts.buildAllPartials(targets)
	assert.NotEmpty(t, ts.looseEndConnector.partials)

	for i, p := range ts.looseEndConnector.partials {
		t.Logf("partials[%d] =\n%s\n", i, p.puzzle)
	}
}
