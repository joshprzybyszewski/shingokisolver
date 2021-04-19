package solvers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

func BenchmarkSolve(b *testing.B) {
	allPuzzles := reader.GetAllPuzzles()

	for _, size := range []int{
		5,
		7,
		10,
		15,
		20,
		25,
	} {
		puzzles := reader.GetPuzzleWithSize(allPuzzles, size)
		b.Run(fmt.Sprintf("Solve %dx%d", size, size), func(b *testing.B) {
			for _, pd := range puzzles {
				_, err := NewSolver(
					pd.NumEdges,
					pd.Nodes,
				).Solve()
				require.NoError(b, err)
			}
		})
	}
}
