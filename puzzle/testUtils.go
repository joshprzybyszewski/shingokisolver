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
	var state model.State
	c := startCoord
	outPuzz := p.DeepCopy()
	for _, s := range steps {
		c, state = outPuzz.AddEdge(s, c)
		switch state {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "unexpected state after adding edge: %+v, %+v\n%s\n", s, c, p)
		}
	}
	t.Logf("BuildTestPuzzle produced: \n%s\n", p)
	return p
}
