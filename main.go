package main

import "github.com/joshprzybyszewski/shingokisolver/puzzlegrid"

func main() {
	size := 6
	nodes := []puzzlegrid.NodeLocation{{
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
	}}

	s := puzzlegrid.NewSolver(size, nodes)
	s.Solve()
}
