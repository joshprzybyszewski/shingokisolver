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
	outPuzz := p.DeepCopy()
	for _, dir := range steps {
		switch s := outPuzz.AddEdge(c, dir); s {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "failure building puzzle", "unexpected state (%s) after adding edge: %+v, %+v\n%s\n", s, dir, c, p)
		}
		c = c.Translate(dir)
	}
	t.Logf("BuildTestPuzzle produced: \n%s\n", outPuzz)
	return outPuzz
}
