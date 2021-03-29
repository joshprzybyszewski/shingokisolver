// +build !prod

package puzzle

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/require"
)

func BuildTestPuzzle(
	t *testing.T,
	g *Puzzle,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) *Puzzle {
	var err error
	c := startCoord
	for _, s := range steps {
		c, g, err = g.AddEdge(s, c)
		require.NoError(t, err, `bad grid! %+v`, g)
	}
	t.Logf("BuildTestPuzzle produced: \n%s\n", g)
	return g
}
