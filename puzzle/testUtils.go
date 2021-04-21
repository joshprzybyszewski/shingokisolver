// +build !prod

package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/require"
)

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
