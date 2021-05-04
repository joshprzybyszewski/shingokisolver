// +build !prod

package puzzle

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

func BuildTestPuzzleWithNoRules(
	t *testing.T,
	numEdges int,
	nls []model.NodeLocation,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) Puzzle {
	if numEdges > state.MaxEdges {
		t.Error(`bad input numEdges`)
	}

	nodeMetas := make([]*nodeMeta, 0, len(nls))
	for _, nl := range nls {
		nc := model.NewCoordFromInts(nl.Row, nl.Col)
		n := model.NewNode(nc, nl.IsWhite, nl.Value)
		nodeMetas = append(nodeMetas, &nodeMeta{
			node:          n,
			twoArmOptions: model.BuildTwoArmOptions(n, numEdges),
		})
	}

	puzz := Puzzle{
		nodes: nodeMetas,
		edges: state.New(numEdges),
		rules: &logic.RuleSet{},
	}

	return BuildTestPuzzle(t, puzz, startCoord, steps...)
}

func BuildTestPuzzle(
	t *testing.T,
	p Puzzle,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) Puzzle {
	c := startCoord
	outPuzz := p
	var s model.State
	for _, dir := range steps {
		outPuzz, s = AddEdge(outPuzz, model.NewEdgePair(c, dir))
		switch s {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "failure building puzzle", "unexpected state (%s) after adding edge: %+v, %+v\n%s\n", s, dir, c, p)
		}
		c = c.Translate(dir)
	}
	t.Logf("BuildTestPuzzle produced: \n%s\n", outPuzz)
	return outPuzz
}
