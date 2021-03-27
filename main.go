package main

import "github.com/joshprzybyszewski/shingokisolver/puzzlegrid"

func main() {
	for _, pd := range getPuzzles() {
		err := puzzlegrid.NewDFSSolver(pd.numEdges, pd.nodes).Solve()
		// puzzlegrid.NewBFSSolver(numEdges, nodes).Solve()
		if err != nil {
			panic(err)
		}
	}
}

type puzzleDef struct {
	numEdges int
	nodes    []puzzlegrid.NodeLocation
}

func getPuzzles() []puzzleDef {
	return []puzzleDef{{
		numEdges: 2,
		nodes: []puzzlegrid.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}},
	}, {
		numEdges: 5,
		nodes: []puzzlegrid.NodeLocation{{
			Row:     3,
			Col:     2,
			IsWhite: false,
			Value:   4,
		}, {
			Row:     3,
			Col:     5,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     4,
			Col:     0,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     5,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     5,
			Col:     3,
			IsWhite: false,
			Value:   5,
		}},
	}, {
		numEdges: 7,
		nodes: []puzzlegrid.NodeLocation{{
			Row:     0,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     0,
			Col:     5,
			IsWhite: false,
			Value:   3,
		}, {
			Row:     1,
			Col:     0,
			IsWhite: true,
			Value:   4,
		}, {
			Row:     1,
			Col:     4,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     1,
			Col:     5,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     2,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     2,
			Col:     6,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     2,
			Col:     7,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     3,
			Col:     7,
			IsWhite: false,
			Value:   5,
		}, {
			Row:     4,
			Col:     1,
			IsWhite: true,
			Value:   4,
		}, {
			Row:     4,
			Col:     6,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     5,
			Col:     2,
			IsWhite: true,
			Value:   3,
		}, {
			Row:     6,
			Col:     0,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     6,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     6,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     6,
			Col:     5,
			IsWhite: false,
			Value:   2,
		}},
	}}
}
