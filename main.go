package main

import (
	"github.com/joshprzybyszewski/shingokisolver/puzzlegrid"
	"github.com/joshprzybyszewski/shingokisolver/reader"
)

func main() {
	for _, pd := range reader.DefaultPuzzles() {
		err := puzzlegrid.NewDFSSolver(pd.NumEdges, pd.Nodes).Solve()
		// puzzlegrid.NewBFSSolver(numEdges, nodes).Solve()
		if err != nil {
			panic(err)
		}
	}
}
