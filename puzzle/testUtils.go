// +build !prod

package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/require"
)

func BuildTestPuzzle(
	t *testing.T,
	p *Puzzle,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) *Puzzle {
	c := startCoord
	outPuzz := p.DeepCopy()
	for _, s := range steps {
		switch outPuzz.AddEdge(c, s) {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "unexpected state after adding edge: %+v, %+v\n%s\n", s, c, p)
		}
		c = c.Translate(s)
	}
	t.Logf("BuildTestPuzzle produced: \n%s\n", p)
	return p
}
